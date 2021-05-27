package cli

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/viper"
	"github.com/stratosnet/stratos-chain/x/register/types"

	//"github.com/Workiva/go-datastructures/threadsafe/err"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/go-bip39"
	"github.com/stratosnet/stratos-chain/app"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/bech32"
	"strings"
	"testing"
)

func TestDoubleSignForCreateResourceNode(t *testing.T) {
	cdc := app.MakeCodec()
	params := ""
	inBuf := bufio.NewReader(strings.NewReader(params))
	//fmt.Println("inBuf: ", inBuf)
	//txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
	txBldr := createFakeTxBuilder()
	//fmt.Println("txBldr: ", txBldr)
	cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

	mnemonic1 := "floor spirit faith hour six reward shoot general judge concert bus drip potato thunder emerge permit salon globe celery reunion mail raccoon output love"
	mnemonic2 := "gossip broccoli vehicle light anchor notable tissue message husband deputy swift sister glimpse nominee basic company child view hand cement holiday age prize fame"

	privKey1, e := createFakePrivKey(mnemonic1)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	fmt.Println("private key1: " + hex.EncodeToString(privKey1.Bytes()))

	pub1 := privKey1.PubKey()
	disc := types.Description{Moniker: "r1", Identity: "", Website: "", SecurityContact: "", Details: ""}
	msg := types.NewMsgCreateResourceNode("sds://resourcenode1", pub1, sdk.Coin{Denom: "stos", Amount: sdk.NewInt(10000)},
		sdk.AccAddress("st1qr9set2jaayzjjpm9tw4f3n6f5zfu3hef8wtaw"), disc, "7")
	//fmt.Println("msg: ", msg)
	//fmt.Println("txBldr: ", txBldr)

	stdsignmsg, sig1, e := GetSignInfo(privKey1, txBldr, msg)
	if e != nil {
		fmt.Println(e.Error())
		return
	}

	privKey2, e := createFakePrivKey(mnemonic2)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	fmt.Println("private key2: " + hex.EncodeToString(privKey2.Bytes()))
	stdsignmsg, sig2, e := GetSignInfo(privKey2, txBldr, msg)
	if e != nil {
		fmt.Println(e.Error())
		return
	}

	tx := authtypes.NewStdTx(stdsignmsg.Msgs, stdsignmsg.Fee, []authtypes.StdSignature{sig1, sig2}, stdsignmsg.Memo)
	fmt.Printf("tx: %#v", tx)
	payload := cliCtx.Codec.MustMarshalBinaryBare(tx)
	fmt.Println()
	fmt.Printf("payload: %#v", payload)
	cliCtx.BroadcastMode = "block"
	broadcastTx, e := cliCtx.BroadcastTx(payload)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	fmt.Printf("broadcastTx: %#v", broadcastTx)
}

func GetSignInfo(privKey crypto.PrivKey, txBldr authtypes.TxBuilder, msg sdk.Msg) (authtypes.StdSignMsg, authtypes.StdSignature, error) {
	pub := privKey.PubKey()
	pubstr, _ := bech32.ConvertAndEncode("stpub", pub.Bytes())
	fmt.Println("pubkey: " + pubstr)
	addrst, e := bech32.ConvertAndEncode("st", pub.Address().Bytes())
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println("address : " + addrst)
	stdsignmsg, e := txBldr.BuildSignMsg([]sdk.Msg{msg})
	if e != nil {
		return authtypes.StdSignMsg{}, authtypes.StdSignature{}, e
	}
	fmt.Println("stdsignmsg : ", stdsignmsg)
	sigBytes, e := privKey.Sign(stdsignmsg.Bytes())
	if e != nil {
		return authtypes.StdSignMsg{}, authtypes.StdSignature{}, e
	}
	sig := authtypes.StdSignature{
		PubKey:    pub,
		Signature: sigBytes,
	}
	return stdsignmsg, sig, nil
}

func createFakePrivKey(mnemonic string) (crypto.PrivKey, error) {
	pass := ""
	seed, e := bip39.NewSeedWithErrorChecking(mnemonic, pass)
	if e != nil {
		fmt.Println(e)
		return nil, e
	}
	fmt.Println("seed: " + hex.EncodeToString(seed))
	masterPriv, ch := hd.ComputeMastersFromSeed(seed)
	fmt.Println("masterkey: " + hex.EncodeToString(masterPriv[:]))
	fmt.Println("chain code: " + hex.EncodeToString(ch[:]))
	hdPath := "m/44'/606'/0'/0/0"
	fmt.Println("path: " + hdPath)
	derivedKey, e := hd.DerivePrivateKeyForPath(masterPriv, ch, hdPath)
	privkey := keys.SecpPrivKeyGen(derivedKey[:])
	fmt.Println("pri key: ", privkey)
	fmt.Println("private key: " + hex.EncodeToString(privkey.Bytes()))
	pub := privkey.PubKey()
	pubstr, _ := bech32.ConvertAndEncode("sdspub", pub.Bytes())
	fmt.Println("pubkey: " + pubstr)
	addrst, e := bech32.ConvertAndEncode("st", pub.Address().Bytes())
	fmt.Println("addrst: ", addrst)
	if e != nil {
		return nil, e
	}
	return privkey, nil
}

////// makes a new MsgCreateIndexingNode.
//////todo: this is temporary for local testing, remove before prod.
func createIndexingNodeMsg(cliCtx context.CLIContext, txBldr auth.TxBuilder) (auth.TxBuilder, sdk.Msg, error) {
	fmt.Println("inside func")
	//amountStr := viper.GetString(FlagAmount)
	amount := sdk.Coin{Denom: "stos", Amount: sdk.NewInt(1000)}
	fmt.Println("amount: ", amount)
	//if err != nil {
	//	return txBldr, nil, err
	//}
	networkAddr := "sds://resourcenode1"
	ownerAddr := "st1qr9set2jaayzjjpm9tw4f3n6f5zfu3hef8wtaw"

	////pkStr := viper.GetString()
	//_, pubk, _ := bech32.DecodeAndConvert("stpub1addwnpepqf8gwzx32nt7fstqy7k6dlj6a923a80cccz3pkuhmhz7lv7qg00wzvrqzru")
	//fmt.Println("pubk: ", pubk)
	//newstr, _ := bech32.ConvertAndEncode("pub", pubk)
	//fmt.Println("newstr: ", newstr)
	pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, "stpub1addwnpepqf8gwzx32nt7fstqy7k6dlj6a923a80cccz3pkuhmhz7lv7qg00wzvrqzru")
	fmt.Println("pk: ", pk)
	if err != nil {
		return txBldr, nil, err
	}
	desc := types.NewDescription(
		viper.GetString("r2"),
		viper.GetString(""),
		viper.GetString(""),
		viper.GetString(""),
		viper.GetString(""),
	)
	msg := types.NewMsgCreateIndexingNode(networkAddr, pk, amount, sdk.AccAddress(ownerAddr), desc)
	fmt.Println()
	fmt.Printf("msg inside func %#v: ", msg)
	return txBldr, msg, nil
}

func createFakeTxBuilder() auth.TxBuilder {
	cdc := codec.New()
	return auth.NewTxBuilder(
		utils.GetTxEncoder(cdc),
		123,
		987,
		0,
		0,
		false,
		"test-chain-localnet",
		"Unit test",
		sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0))),
		sdk.DecCoins{sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.NewDecWithPrec(0, sdk.Precision))},
	)
}

// CreateResourceNodeCmd will create a file upload tx and sign it with the given key.
//func TestCreateResourceNodeCmd(t *testing.T) {
//	mnemonic := "floor spirit faith hour six reward shoot general judge concert bus drip potato thunder emerge permit salon globe celery reunion mail raccoon output love"
//	pass := ""
//	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, pass)
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println("seed: " + hex.EncodeToString(seed))
//	masterPriv, ch := hd.ComputeMastersFromSeed(seed)
//	fmt.Println("masterkey: " + hex.EncodeToString(masterPriv[:]))
//	fmt.Println("chain code: " + hex.EncodeToString(ch[:]))
//	hdPath := "m/44'/606'/0'/0/0"
//	fmt.Println("path: " + hdPath)
//	derivedKey, err := hd.DerivePrivateKeyForPath(masterPriv, ch, hdPath)
//	privkey := keys.SecpPrivKeyGen(derivedKey[:])
//	fmt.Println("pri key: " , privkey)
//	fmt.Println("private key: " + hex.EncodeToString(privkey.Bytes()))
//	pub := privkey.PubKey()
//	pubstr, _ := bech32.ConvertAndEncode("sdspub", pub.Bytes())
//	fmt.Println("pubkey: " + pubstr)
//	addrst, err := bech32.ConvertAndEncode("st", pub.Address().Bytes())
//	if err != nil {
//		fmt.Println(err)
//	}
//	//fmt.Println("address:"+ pub.Address().String())
//	fmt.Println("address : " + addrst)
//
//	txBldr := createFakeTxBuilder()
//	fmt.Println("Fees: ", txBldr.Fees())
//	fmt.Println("GasPrice: ", txBldr.GasPrices())
//
//	disc := types.Description{Moniker: "r1", Identity: "", Website: "", SecurityContact: "", Details: ""}
//	msg := types.NewMsgCreateResourceNode("sds://resourcenode1", pub, sdk.Coin{Denom: "stos", Amount: sdk.NewInt(10000)},
//		sdk.AccAddress("st1qr9set2jaayzjjpm9tw4f3n6f5zfu3hef8wtaw"), disc, "computation")
//	fmt.Println("MSG: ", msg)
//	//txBldr := auth.NewTxBuilder(msg).WithTxEncoder(utils.GetTxEncoder(auth.ModuleCdc))
//	stdsignmsg, err := txBldr.BuildSignMsg([]sdk.Msg{msg})
//	if err != nil {
//		fmt.Println("SignMsg Error: ", err)
//	}
//	sigBytes, err := privkey.Sign(stdsignmsg.Bytes())
//	fmt.Printf("sigBytes: %#v" , sigBytes)
//	fmt.Println()
//
//	sig1 := authtypes.StdSignature{
//		PubKey:    pub,
//		Signature: sigBytes,
//	}
//	tx := authtypes.NewStdTx(stdsignmsg.Msgs, stdsignmsg.Fee, []authtypes.StdSignature{sig1}, stdsignmsg.Memo)
//	fmt.Printf("Tx: %#v" , tx)
//	fmt.Println()
//	//cliCtx := context.NewCLIContext().WithCodec(auth.ModuleCdc)
//	cliCtx := context.CLIContext{}
//	cliCtx.BroadcastMode = "async"
//	//broadcastTx, err := cliCtx.BroadcastTx(auth.ModuleCdc.MustMarshalJSON(tx))
//	broadcastTx, err := cliCtx.BroadcastTx(msg.GetSignBytes())
//	if err != nil {
//		fmt.Println("Broadcast Error: ", err)
//	}
//	fmt.Printf("broadcastTx: %#v", broadcastTx)

//txBldr := auth.NewTxBuilder(inBuf).WithTxEncoder(utils.GetTxEncoder(auth.ModuleCdc))
//msg := types.NewMsgCreateResourceNode("sds://resourcenode1", pubstr, 10000000, sdk.AccAddress("st1qr9set2jaayzjjpm9tw4f3n6f5zfu3hef8wtaw"), "", 7)
//Bldr := auth.NewTxBuilderFromCLI(msg).WithTxEncoder(utils.GetTxEncoder(auth.ModuleCdc))
//stdsignmsg, err := txBldr.BuildSignMsg([]sdk.Msg{msg})
//sigBytes, err := privkey.Sign(stdsignmsg.Bytes())
//sig1 := authtypes.StdSignature{
//	PubKey:    pub,
//	Signature: sigBytes,
//}
//tx := authtypes.NewStdTx(stdsignmsg.Msgs, stdsignmsg.Fee, []authtypes.StdSignature{sig1}, stdsignmsg.Memo)
//cliCtx := context.NewCLIContext().WithCodec(auth.ModuleCdc)
//broadcastTx, err := cliCtx.BroadcastTx(tx)
//if err != nil {
//	return err
//}
//cmd.Flags().AddFlagSet(FsPk)
//cmd.Flags().AddFlagSet(FsAmount)
//cmd.Flags().AddFlagSet(FsNetworkID)
//
//cmd.MarkFlagRequired(flags.FlagFrom)
//cmd.MarkFlagRequired(FlagAmount)
//cmd.MarkFlagRequired(FlagPubKey)
//cmd.MarkFlagRequired(FlagNetworkID)
//return cmd
//}

//func buildCreateResourceNodeMsg_test(cliCtx context.CLIContext, txBldr auth.TxBuilder) (auth.TxBuilder, sdk.Msg, error) {
//	amountStr := viper.GetString(FlagAmount)
//	amount, err := sdk.ParseCoin(amountStr)
//	if err != nil {
//		return txBldr, nil, err
//	}
//
//	networkID := viper.GetString(FlagNetworkID)
//	ownerAddr := cliCtx.GetFromAddress()
//	pkStr := viper.GetString(FlagPubKey)
//	nodeTypeRef := viper.GetInt(FlagNodeType)
//
////pk, er := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, pkStr)
//////pk, er := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypesdsPub, pkStr)
////if er != nil {
////	return txBldr, nil, err
////}
//
////desc := types.NewDescription(
////	viper.GetString(FlagMoniker),
////	viper.GetString(FlagIdentity),
////	viper.GetString(FlagWebsite),
////	viper.GetString(FlagSecurityContact),
////	viper.GetString(FlagDetails),
////)
////if !ValueInSlice(types.NodeType(nodeTypeRef).Type(), types.NodeTypes) {
////	return txBldr, nil, types.ErrNodeType
////}
//
//// validate nodeTypeRef
////	if t := types.NodeType(nodeTypeRef).Type(); t == "UNKNOWN" {
////		return txBldr, nil, types.ErrNodeType
////	}
////	msg := types.NewMsgCreateResourceNode(networkID, pk, amount, ownerAddr, desc, fmt.Sprintf("%d: %s", nodeTypeRef, types.NodeType(nodeTypeRef).Type()))
////	return txBldr, msg, nil
////}
////
//
////// CreateIndexingNodeCmd will create a file upload tx and sign it with the given key.
////func CreateIndexingNodeCmd(cdc *codec.Codec) *cobra.Command {
////	cmd := &cobra.Command{
////		Use:   "create-indexing-node",
////		Short: "create new indexing node",
////		RunE: func(cmd *cobra.Command, args []string) error {
////			inBuf := bufio.NewReader(cmd.InOrStdin())
////			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
////			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
////			txBldr, msg, err := buildCreateIndexingNodeMsg(cliCtx, txBldr)
////			if err != nil {
////				return err
////			}
////			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
////		},
////	}
////	cmd.Flags().AddFlagSet(FsPk)
////	cmd.Flags().AddFlagSet(FsAmount)
////	cmd.Flags().AddFlagSet(FsNetworkID)
////
////	cmd.MarkFlagRequired(flags.FlagFrom)
////	cmd.MarkFlagRequired(FlagAmount)
////	cmd.MarkFlagRequired(FlagPubKey)
////	cmd.MarkFlagRequired(FlagNetworkID)
////	return cmd
////}
////
////// makes a new CreateResourceNodeMsg.
////func buildCreateResourceNodeMsg(cliCtx context.CLIContext) (sdk.Msg, error) {
////	amountStr := viper.GetString(FlagAmount)
////	amount, err := sdk.ParseCoin(amountStr)
////	if err != nil {
////		return nil, err
////	}
////	networkAddr := viper.GetString(FlagNetworkID)
////	ownerAddr := cliCtx.GetFromAddress()
////	pkStr := viper.GetString(FlagPubKey)
////	pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, pkStr)
////	if err != nil {
////		return nil, err
////	}
////	desc := types.NewDescription(
////		viper.GetString(FlagMoniker),
////		viper.GetString(FlagIdentity),
////		viper.GetString(FlagWebsite),
////		viper.GetString(FlagSecurityContact),
////		viper.GetString(FlagDetails),
////	)
////	msg := types.NewMsgCreateResourceNode(networkAddr, pk, amount, ownerAddr, desc)
////	return msg, nil
////}
////func RemoveResourceNodeCmd(cdc *codec.Codec) *cobra.Command {
////	cmd := &cobra.Command{
////		Use:   "remove-resource-node [resource_node_address] [owner_address]",
////		Args:  cobra.ExactArgs(2),
////		Short: "remove resource node",
////		RunE: func(cmd *cobra.Command, args []string) error {
////			inBuf := bufio.NewReader(cmd.InOrStdin())
////			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
////			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[1]).WithCodec(cdc)
////			resourceNodeAddr, err := sdk.AccAddressFromBech32(args[0])
////			if err != nil {
////				return err
////			}
////			ownerAddr := cliCtx.GetFromAddress()
////			msg := types.NewMsgRemoveResourceNode(resourceNodeAddr, ownerAddr)
////			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
////		},
////	}
////	return cmd
////}
////func RemoveIndexingNodeCmd(cdc *codec.Codec) *cobra.Command {
////	cmd := &cobra.Command{
////		Use:   "remove-indexing-node [indexing_node_address] [owner_address]",
////		Args:  cobra.ExactArgs(2),
////		Short: "remove indexing node",
////		RunE: func(cmd *cobra.Command, args []string) error {
////			inBuf := bufio.NewReader(cmd.InOrStdin())
////			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
////			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[1]).WithCodec(cdc)
////			indexingNodeAddr, err := sdk.AccAddressFromBech32(args[0])
////			if err != nil {
////				return err
////			}
////			ownerAddr := cliCtx.GetFromAddress()
////			msg := types.NewMsgRemoveIndexingNode(indexingNodeAddr, ownerAddr)
////			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
////		},
////	}
////	return cmd
////}
////
