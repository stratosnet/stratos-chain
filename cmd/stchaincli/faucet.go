package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
)

const (
	flagFundFrom = "from" // optional
	flagAmt      = "amt"  // denom fixed as ustos
	flagPort     = "port"
	flagChainId  = "chain-id"
	flagAddrCap  = "addr-cap"
	flagIpCap    = "ip-cap"

	defaultOutputFlag     = "text"
	defaultKeyringBackend = "test"
	defaultDenom          = "ustos"
	defaultChainId        = "test-chain"
	defaultAddrCap        = 1
	defaultIpCap          = 3
	capDuration           = 60 // in minutes

	maxAmtFaucet    = 100000000000
	requestInterval = 100 * time.Millisecond
)

// used in request channel
type FaucetReq struct {
	ToAddr  sdk.AccAddress
	resChan chan FaucetRsp
}

// used in response channel
type FaucetRsp struct {
	ErrorMsg   string
	TxResponse sdk.TxResponse
	Seq        uint64
}

// used for restful response
type RestFaucetRsp struct {
	ErrorMsg   string
	TxResponse sdk.TxResponse
}

type FaucetToMiddleware struct {
	Cap       int // maximum faucet cap to an individual addr during an hour
	AddrCache ttlcache.SimpleCache
}

type FromIpMiddleware struct {
	Cap     int // maximum accessing times during an hour
	IpCache ttlcache.SimpleCache
}

func (ftm *FaucetToMiddleware) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		addr := vars["address"]
		if ftm.checkCap(addr) {
			h.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Faucet request to address [" + addr + "] exceeds hourly cap (" + strconv.Itoa(ftm.Cap) + " request(s) per hour)"))
		}
	})
}

func (ftm *FaucetToMiddleware) checkCap(toAddr string) bool {
	val, _ := ftm.AddrCache.Get(toAddr)
	if val == nil {
		ftm.AddrCache.Set(toAddr, 1)
		return true
	}

	if val.(int) >= ftm.Cap {
		return false
	}
	ftm.AddrCache.Set(toAddr, val.(int)+1)
	return true
}

func (fim *FromIpMiddleware) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		realIp := getRealAddr(r)
		if fim.checkCap(realIp) {
			h.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Faucet request from Ip " + realIp + " exceeds hourly cap (" + strconv.Itoa(fim.Cap) + " request(s) per hour)!"))
		}
	})
}

func (fim *FromIpMiddleware) checkCap(fromIp string) bool {
	val, _ := fim.IpCache.Get(fromIp)
	if val == nil {
		fim.IpCache.Set(fromIp, 1)
		return true
	}

	if val.(int) >= fim.Cap {
		return false
	}
	fim.IpCache.Set(fromIp, val.(int)+1)
	return true
}

// global to load command line args
var (
	faucetArgs = FaucetArgs{}
	seqInfo    = SeqInfo{startSeq: 0, lastSuccSeq: 0}
)

// struct to hold the command-line args
type FaucetArgs struct {
	from  sdk.AccAddress
	coins sdk.Coins
	port  string
}

type SeqInfo struct {
	lastSuccSeq int
	startSeq    int
	mu          sync.Mutex
}

func (si *SeqInfo) incrLastSuccSeq(succSeq uint64) {
	si.mu.Lock()
	defer si.mu.Unlock()
	if si.lastSuccSeq < int(succSeq) {
		si.lastSuccSeq = int(succSeq)
	}
}

func (si *SeqInfo) getNewSeq(newStartSeq int) int {
	si.mu.Lock()
	defer si.mu.Unlock()

	if si.lastSuccSeq < newStartSeq-1 {
		si.lastSuccSeq = newStartSeq - 1
		return newStartSeq
	} else {
		return si.lastSuccSeq + 1
	}
}

func FaucetJobFromCh(faucetReq *chan FaucetReq, cliCtx context.CLIContext, txBldr authtypes.TxBuilder, from sdk.AccAddress, coin sdk.Coin, quit chan os.Signal) {
	for {
		select {
		case <-quit:

			return
		case fReq := <-*faucetReq:
			resChan := fReq.resChan
			// get latest seq
			_, latestSeq, err := authtypes.NewAccountRetriever(cliCtx).GetAccountNumberSequence(from)
			if err != nil {
				faucetRsp := FaucetRsp{ErrorMsg: "Node is under maintenance, please try again later!"}
				resChan <- faucetRsp
				continue
			}
			newSeq := seqInfo.getNewSeq(int(latestSeq))
			//fmt.Print(fmt.Sprintf("sequence in this tx: %d\n", newSeq))
			err = doTransfer(cliCtx,
				txBldr.
					WithSequence(uint64(newSeq)).
					WithChainID(viper.GetString(flags.FlagChainID)).
					WithGas(uint64(400000)).
					WithMemo(strconv.Itoa(newSeq)),
				fReq.ToAddr, faucetArgs.from, coin, &resChan)
			if err != nil {
				faucetRsp := FaucetRsp{ErrorMsg: err.Error()}
				resChan <- faucetRsp
			}
		}
	}
}

// GetFaucetCmd returns faucet cobra Command
func GetFaucetCmd(cdc *codec.Codec) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "faucet",
		Short: "Run a faucet server",
		Args:  cobra.RangeArgs(0, 7),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if !viper.IsSet(flagFundFrom) {
				return fmt.Errorf("fund-from not specified")
			}
			if !viper.IsSet(flags.FlagChainID) {
				return fmt.Errorf("chain-id not specified")
			}
			if !viper.IsSet(flags.FlagKeyringBackend) {
				viper.Set(flags.FlagKeyringBackend, defaultKeyringBackend)
			}

			addrCap := viper.GetInt(flagAddrCap)
			ipCap := viper.GetInt(flagIpCap)

			fmt.Print("Set hourly addrCap = " + strconv.Itoa(addrCap) + ", hourly ipCap = " + strconv.Itoa(ipCap))

			faucetToCache := ttlcache.NewCache()
			faucetToCache.SetTTL(capDuration * time.Minute)
			faucetToCache.SkipTTLExtensionOnHit(true)
			faucetToCache.SetCacheSizeLimit(65535)
			ftm := FaucetToMiddleware{AddrCache: faucetToCache, Cap: addrCap}

			fromIpCache := ttlcache.NewCache()
			fromIpCache.SetTTL(capDuration * time.Minute)
			fromIpCache.SkipTTLExtensionOnHit(true)
			fromIpCache.SetCacheSizeLimit(65535)
			fim := FromIpMiddleware{IpCache: fromIpCache, Cap: ipCap}

			fromAddr := viper.GetString(flagFundFrom)
			fromAddrBytes, err := sdk.AccAddressFromBech32(fromAddr)
			if err != nil {
				return fmt.Errorf("failed to parse bech32 address fro FROM Address: %w", err)
			}
			faucetArgs.from = fromAddrBytes

			// start threads
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			viper.Set(flags.FlagSkipConfirmation, true)
			viper.Set(cli.OutputFlag, defaultOutputFlag)

			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, faucetArgs.from.String()).WithCodec(cdc)

			faucetArgs.port = viper.GetString(flagPort)

			fmt.Print("\nfunding address: ", "addr", faucetArgs.from.String())
			var toTransferAmt int
			if toTransferAmt = viper.GetInt(flagAmt); toTransferAmt <= 0 || toTransferAmt > maxAmtFaucet {
				return fmt.Errorf("invalid amount in faucet")
			}
			coin := sdk.Coin{Amount: sdk.NewInt(int64(toTransferAmt)), Denom: defaultDenom}
			faucetArgs.coins = sdk.Coins{coin}

			fmt.Print("\nStarting faucet...")

			// listen to localhost:26600
			listener, err := net.Listen("tcp", ":"+faucetArgs.port)
			fmt.Print("\nlisten to [" + ":" + faucetArgs.port + "]")
			// router
			r := mux.NewRouter()
			// health check
			r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("ok\n"))
			})
			faucetReqCh := make(chan FaucetReq, 10000)
			//faucet
			r.HandleFunc("/faucet/{address}", func(writer http.ResponseWriter, request *http.Request) {
				vars := mux.Vars(request)
				addr := vars["address"]
				fmt.Println("get request from addr: ", addr)
				toAddr, err := sdk.AccAddressFromBech32(addr)
				if err != nil {
					writer.WriteHeader(http.StatusBadRequest)
					writer.Write([]byte(err.Error()))
				}
				resChan := make(chan FaucetRsp)
				faucetReq := FaucetReq{ToAddr: toAddr, resChan: resChan}
				faucetReqCh <- faucetReq

				faucetRsp := <-resChan
				if int(faucetRsp.TxResponse.Code) < 1 && len(faucetRsp.ErrorMsg) == 0 {
					// sigverify pass
					seqInfo.incrLastSuccSeq(faucetRsp.Seq)
				}
				fmt.Println("tx send=", faucetRsp.TxResponse.TxHash, ", height=", faucetRsp.TxResponse.Height, ", errorMsg=", faucetRsp.ErrorMsg)
				restRsp := &RestFaucetRsp{ErrorMsg: faucetRsp.ErrorMsg, TxResponse: faucetRsp.TxResponse}
				rest.PostProcessResponseBare(writer, cliCtx, restRsp)
				return
			}).Methods("POST")
			// ipCap check has higher priority than toAddrCap
			r.Use(fim.Middleware)
			r.Use(ftm.Middleware)

			quit := make(chan os.Signal, 1)
			signal.Notify(quit,
				syscall.SIGTERM,
				syscall.SIGINT,
				syscall.SIGQUIT,
				syscall.SIGKILL,
				syscall.SIGHUP,
			)

			go FaucetJobFromCh(&faucetReqCh, cliCtx, txBldr, faucetArgs.from, coin, quit)
			//start the server
			err = http.Serve(listener, r)
			if err != nil {
				fmt.Println(err.Error())
			}
			close(quit)
			// print stats
			fmt.Println("####################################################################")
			fmt.Println("################        Terminating faucet        ##################")
			fmt.Println("####################################################################")
			return nil
		},
	}

	cmd.Flags().String(flags.FlagKeyringBackend, defaultKeyringBackend, "Select keyring's backend (os|file|test)")
	cmd.Flags().String(flagAmt, "", "amt to transfer in faucet")
	cmd.Flags().String(flagFundFrom, "", "fund from address")
	cmd.Flags().String(flagPort, "26600", "port of faucet server")
	cmd.Flags().Int(flagAddrCap, defaultAddrCap, "hourly cap of faucet to a particular account address")
	cmd.Flags().Int(flagIpCap, defaultIpCap, "hourly cap of faucet from a particular IP")

	return cmd
}

func doTransfer(cliCtx context.CLIContext, txBldr authtypes.TxBuilder, to sdk.AccAddress, from sdk.AccAddress, coin sdk.Coin, resChan *chan FaucetRsp) error {
	//// build and sign the transaction, then broadcast to Tendermint
	msg := bank.NewMsgSend(from, to, sdk.Coins{coin})
	msgs := []sdk.Msg{msg}
	txBldr, err := utils.PrepareTxBuilder(txBldr, cliCtx)
	if err != nil {
		return err
	}

	fromName := cliCtx.GetFromName()

	if txBldr.SimulateAndExecute() || cliCtx.Simulate {
		txBldr, err = utils.EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			return err
		}

		gasEst := utils.GasEstimateResponse{GasEstimate: txBldr.Gas()}
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", gasEst.String())
	}

	if !cliCtx.SkipConfirm {
		stdSignMsg, err := txBldr.BuildSignMsg(msgs)
		if err != nil {
			return err
		}

		var json []byte
		if viper.GetBool(flags.FlagIndentResponse) {
			json, err = cliCtx.Codec.MarshalJSONIndent(stdSignMsg, "", "  ")
			if err != nil {
				panic(err)
			}
		} else {
			json = cliCtx.Codec.MustMarshalJSON(stdSignMsg)
		}

		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", json)

		buf := bufio.NewReader(os.Stdin)
		ok, err := input.GetConfirmation("confirm transaction before signing and broadcasting", buf)
		if err != nil || !ok {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", "cancelled transaction")
			return err
		}
	}
	// build and sign the transaction
	txBytes, err := txBldr.BuildAndSign(fromName, keys.DefaultKeyPass, msgs)
	if err != nil {
		return err
	}

	// broadcast to a Tendermint node
	res, err := cliCtx.BroadcastTxCommit(txBytes)
	if err != nil {
		return err
	}
	faucetRsp := FaucetRsp{TxResponse: res, Seq: txBldr.Sequence()}
	*resChan <- faucetRsp
	return nil
}

func getRealAddr(r *http.Request) string {
	fmt.Printf("parsing IP: remote_addr=[%s], X-Forwarded-For=[%s], X-Real-Ip=[%s]", r.RemoteAddr, r.Header.Get("X-Forwarded-For"), r.Header.Get("X-Real-Ip"))
	remoteIP := ""
	// the default is the originating ip. but we try to find better options because this is almost
	// never the right IP
	if parts := strings.Split(r.RemoteAddr, ":"); len(parts) == 2 {
		remoteIP = parts[0]
		fmt.Printf("==== remoteIp[%s] parsed from RemoteAddr", remoteIP)
	}
	// If we have a forwarded-for header, take the address from there
	if xff := strings.Trim(r.Header.Get("X-Forwarded-For"), ","); len(xff) > 0 {
		addrs := strings.Split(xff, ",")
		lastFwd := addrs[len(addrs)-1]
		if ip := net.ParseIP(lastFwd); ip != nil {
			remoteIP = ip.String()
		}
		fmt.Printf("==== remoteIp[%s] parsed from X-Forwarded-For", remoteIP)
		// parse X-Real-Ip header
	} else if xri := r.Header.Get("X-Real-Ip"); len(xri) > 0 {
		if ip := net.ParseIP(xri); ip != nil {
			remoteIP = ip.String()
		}
		fmt.Printf("==== remoteIp[%s] parsed from X-Real-Ip", remoteIP)
	}
	return remoteIP
}
