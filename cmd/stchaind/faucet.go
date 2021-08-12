package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/ReneKroon/ttlcache/v2"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	tmtypes "github.com/tendermint/tendermint/types"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	flagFrom = "from" // optional
	flagAmt  = "amt"  // denom fixed as ustos
	flagPort = "port"

	flagAddrCap = "addr-cap"
	flagIpCap   = "ip-cap"

	defaultNodeURI        = "tcp://127.0.0.1:26657"
	defaultKeyringBackend = "test"
	defaultHome           = "build/node/stchaincli"
	defaultDenom          = "ustos"
	defaultChainId        = "test-chain"
	defaultAddrCap        = 1
	defaultIpCap          = 3

	maxAmtFaucet = 100000000000
)

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
var faucetArgs = FaucetArgs{}

// struct to hold the command-line args
type FaucetArgs struct {
	from  sdk.AccAddress
	coins sdk.Coins
	port  string
}

// AddFaucetCmd returns faucet cobra Command
func AddFaucetCmd(
	ctx *server.Context, cdc *codec.Codec, defaultNodeHome, defaultClientHome string,
) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "faucet",
		Short: "Run a faucet cmd",
		Args:  cobra.RangeArgs(0, 7),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			addrCap := viper.GetInt(flagAddrCap)
			ipCap := viper.GetInt(flagIpCap)

			ctx.Logger.Info("Set addrCap = " + strconv.Itoa(addrCap) + ", ipCap = " + strconv.Itoa(ipCap))

			faucetToCache := ttlcache.NewCache()
			faucetToCache.SetTTL(60 * time.Minute)
			faucetToCache.SkipTTLExtensionOnHit(true)
			faucetToCache.SetCacheSizeLimit(65535)
			ftm := FaucetToMiddleware{AddrCache: faucetToCache, Cap: addrCap}

			fromIpCache := ttlcache.NewCache()
			fromIpCache.SetTTL(60 * time.Minute)
			fromIpCache.SkipTTLExtensionOnHit(true)
			fromIpCache.SetCacheSizeLimit(65535)
			fim := FromIpMiddleware{IpCache: fromIpCache, Cap: ipCap}

			if viper.IsSet(flagFrom) {
				fromAddr := viper.GetString(flagFrom)
				fromAddrBytes, err := sdk.AccAddressFromBech32(fromAddr)
				if err != nil {
					return fmt.Errorf("failed to parse bech32 address fro FROM Address: %w", err)
				}
				faucetArgs.from = fromAddrBytes
			}

			if faucetArgs.from == nil {
				genesis := ctx.Config.GenesisFile()
				faucetArgs.from, err = getFirstAccAddressFromGenesis(cdc, genesis)
				if err != nil {
					return fmt.Errorf("failed to parse genesis: %w", err)
				}
				fmt.Printf("No sender account specified, using account 0 for faucet\n")
			}
			faucetArgs.port = viper.GetString(flagPort)

			ctx.Logger.Info("funding address: ", "addr", faucetArgs.from.String())
			var toTransferAmt int
			if toTransferAmt = viper.GetInt(flagAmt); toTransferAmt <= 0 || toTransferAmt > maxAmtFaucet {
				return fmt.Errorf("Invalid amount in faucet")
			}
			coin := sdk.Coin{Amount: sdk.NewInt(int64(toTransferAmt)), Denom: defaultDenom}
			faucetArgs.coins = sdk.Coins{coin}

			ctx.Logger.Info("Starting faucet...")

			// start threads
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			viper.Set(flags.FlagBroadcastMode, "sync")
			if !viper.IsSet(flags.FlagChainID) {
				viper.Set(flags.FlagChainID, defaultChainId)
			}
			viper.Set(flags.FlagSkipConfirmation, true)
			if !viper.IsSet(flags.FlagKeyringBackend) {
				viper.Set(flags.FlagKeyringBackend, defaultKeyringBackend)
			}
			if !viper.IsSet(flags.FlagNode) {
				viper.Set(flags.FlagNode, defaultNodeURI)
			}
			if !viper.IsSet(flags.FlagHome) {
				viper.Set(flags.FlagHome, defaultHome)
			}
			viper.Set(flags.FlagTrustNode, true)
			viper.Set(cli.OutputFlag, defaultOutputFlag)
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, faucetArgs.from.String()).WithCodec(cdc)

			// listen to localhost:26600
			listener, err := net.Listen("tcp", ":"+faucetArgs.port)
			ctx.Logger.Info("listen to [" + ":" + faucetArgs.port + "]")
			// router
			r := mux.NewRouter()
			// health check
			r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("ok\n"))
			})
			//faucet
			r.HandleFunc("/faucet/{address}", func(writer http.ResponseWriter, request *http.Request) {
				vars := mux.Vars(request)
				addr := vars["address"]
				realIp := getRealAddr(request)
				toAddr, err := sdk.AccAddressFromBech32(addr)
				if err != nil {
					writer.WriteHeader(http.StatusBadRequest)
					writer.Write([]byte(err.Error()))
				}
				ctx.Logger.Info("received faucet request: ", "toAddr", addr, "fromIp", realIp)
				err = doTransfer(cliCtx, txBldr.WithChainID(viper.GetString(flags.FlagChainID)), toAddr, faucetArgs.from, coin) // send coin to temp account
				if err != nil {
					writer.WriteHeader(http.StatusBadRequest)
					writer.Write([]byte(err.Error()))
				}
				ctx.Logger.Info("send", "addr", addr, "amount", coin.String())

				msg := fmt.Sprintf("Faucet successfully triggered, sending %s to %s\n", coin.String(), addr)
				writer.WriteHeader(http.StatusOK)
				writer.Write([]byte(msg))
				return
			}).Methods("POST")
			// ipCap check has higher priority than toAddrCap
			r.Use(fim.Middleware)
			r.Use(ftm.Middleware)
			//start the server
			err = http.Serve(listener, r)
			if err != nil {
				fmt.Println(err.Error())
			}
			// print stats
			fmt.Println("####################################################################")
			fmt.Println("################        Terminating faucet        ##################")
			fmt.Println("####################################################################")
			return nil
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().String(flagClientHome, defaultClientHome, "client's home directory")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	cmd.Flags().String(flagAmt, "", "amt to transfer in faucet")
	cmd.Flags().String(flagFrom, "", "from address")
	cmd.Flags().String(flagPort, "26600", "port of faucet server")
	cmd.Flags().String(flags.FlagChainID, "", "chain id")
	cmd.Flags().Int(flagAddrCap, defaultAddrCap, "hourly cap of faucet to a particular account address")
	cmd.Flags().Int(flagIpCap, defaultIpCap, "hourly cap of faucet from a particular IP")

	return cmd
}

func checkCap(s string, s2 string) bool {
	return true
}

func getFirstAccAddressFromGenesis(cdc *codec.Codec, genesisFilePath string) (accAddr sdk.AccAddress, err error) {
	var genDoc *tmtypes.GenesisDoc
	if genDoc, err = tmtypes.GenesisDocFromFile(strings.ReplaceAll(genesisFilePath, "cli", "d")); err != nil {
		return nil, fmt.Errorf("error loading genesis doc from %s: %s", genesisFilePath, err.Error())
	}
	var genState map[string]json.RawMessage
	if err = cdc.UnmarshalJSON(genDoc.AppState, &genState); err != nil {
		return nil, fmt.Errorf("error unmarshalling genesis doc %s: %s", genesisFilePath, err.Error())
	}
	var addresses []sdk.AccAddress
	auth.GenesisAccountIterator{}.IterateGenesisAccounts(
		cdc, genState, func(acc exported.Account) (stop bool) {
			addresses = append(addresses, acc.GetAddress())
			return false
		},
	)
	if len(addresses) > 0 {
		return addresses[0], nil
	}
	return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, "No account initiated in genesis")
}

func doTransfer(cliCtx context.CLIContext, txBldr authtypes.TxBuilder, to sdk.AccAddress, from sdk.AccAddress, coin sdk.Coin) error {
	//// build and sign the transaction, then broadcast to Tendermint
	msg := bank.NewMsgSend(from, to, sdk.Coins{coin})
	cliCtx.BroadcastMode = "sync"
	err := utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
	return err
}

func getRealAddr(r *http.Request) string {
	remoteIP := ""
	// the default is the originating ip. but we try to find better options because this is almost
	// never the right IP
	if parts := strings.Split(r.RemoteAddr, ":"); len(parts) == 2 {
		remoteIP = parts[0]
	}
	// If we have a forwarded-for header, take the address from there
	if xff := strings.Trim(r.Header.Get("X-Forwarded-For"), ","); len(xff) > 0 {
		addrs := strings.Split(xff, ",")
		lastFwd := addrs[len(addrs)-1]
		if ip := net.ParseIP(lastFwd); ip != nil {
			remoteIP = ip.String()
		}
		// parse X-Real-Ip header
	} else if xri := r.Header.Get("X-Real-Ip"); len(xri) > 0 {
		if ip := net.ParseIP(xri); ip != nil {
			remoteIP = ip.String()
		}
	}
	return remoteIP
}
