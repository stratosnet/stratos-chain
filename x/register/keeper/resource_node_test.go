package keeper

//
//import (
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	"github.com/tendermint/tendermint/crypto/ed25519"
//	"log"
//	"testing"
//)
//
//var (
//	ppNodeOwner1   = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
//	ppNodeOwner2   = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
//	ppNodeOwner3   = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
//	ppNodeOwner4   = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
//	ppNodeOwnerNew = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
//
//	ppNodePubKey1   = ed25519.GenPrivKey().PubKey()
//	ppNodeAddr1     = sdk.AccAddress(ppNodePubKey1.Address())
//	ppInitialStake1 = sdk.NewInt(100000000)
//
//	ppNodePubKey2   = ed25519.GenPrivKey().PubKey()
//	ppNodeAddr2     = sdk.AccAddress(ppNodePubKey2.Address())
//	ppInitialStake2 = sdk.NewInt(100000000)
//
//	ppNodePubKey3   = ed25519.GenPrivKey().PubKey()
//	ppNodeAddr3     = sdk.AccAddress(ppNodePubKey3.Address())
//	ppInitialStake3 = sdk.NewInt(100000000)
//
//	ppNodePubKey4   = ed25519.GenPrivKey().PubKey()
//	ppNodeAddr4     = sdk.AccAddress(ppNodePubKey4.Address())
//	ppInitialStake4 = sdk.NewInt(100000000)
//
//	ppNodePubKeyNew = ed25519.GenPrivKey().PubKey()
//	ppNodeAddrNew   = sdk.AccAddress(ppNodePubKeyNew.Address())
//	ppNodeStakeNew  = sdk.NewInt(100000000)
//)
//
//func TestOzoneLimitChange(t *testing.T) {
//	ctx, _, _, k, _ := CreateTestInput(t, false)
//
//	initialStakeTotal := sdk.NewInt(43000)
//	k.SetInitialGenesisStakeTotal(ctx, initialStakeTotal)
//	k.SetRemainingOzoneLimit(ctx, initialStakeTotal)
//
//	resouceNodeTokens := make([]sdk.Int, 0)
//	numSeq := 100
//	resourceNodeStake := sdk.NewInt(19000)
//	for i := 0; i < numSeq; i++ {
//		resouceNodeTokens = append(resouceNodeTokens, resourceNodeStake)
//	}
//	log.Printf("Before: remaining ozone limit is %v", k.GetRemainingOzoneLimit(ctx))
//	for i, val := range resouceNodeTokens {
//		ozoneLimitChange := k.increaseOzoneLimitByAddStake(ctx, val)
//		log.Printf("Add resourceNode #%v(stake=%v), ozone limit increases by %v, remaining ozone limit is %v", i, resourceNodeStake, ozoneLimitChange, k.GetRemainingOzoneLimit(ctx))
//	}
//}
