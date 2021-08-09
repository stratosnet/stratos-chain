package register

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/stratosnet/stratos-chain/x/register/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

func TestRegister(t *testing.T) {

	/********************* initialize mock app *********************/
	SetConfig()
	//mApp, k, accountKeeper, bankKeeper, stakingKeeper, registerKeeper := getMockApp(t)
	mApp, k, _, _ := getMockApp(t)
	accounts := setupAccounts(mApp)
	mock.SetGenesis(mApp, accounts)

	header := abci.Header{}
	ctx := mApp.BaseApp.NewContext(true, header)

	/********************* sign twice and send register resource node msg *********************/

	header = abci.Header{Height: mApp.LastBlockHeight() + 1}
	ctx = mApp.BaseApp.NewContext(true, header)

	registerResNodeMsg := types.NewMsgCreateResourceNode(
		"sds://resourceNode3",
		resNodePubKey3,
		sdk.NewCoin(k.BondDenom(ctx), resNodeInitStake),
		resOwnerAddr3,
		NewDescription("sds://resourceNode3", "", "", "", ""),
		"4",
	)
	t.Log("registerResNodeMsg: ", registerResNodeMsg)

	resOwnerAcc3 := mApp.AccountKeeper.GetAccount(ctx, resOwnerAddr3)
	accNumOwner := resOwnerAcc3.GetAccountNumber()
	t.Log("accNumOwner: ", accNumOwner)
	accSeqOwner := resOwnerAcc3.GetSequence()
	t.Log("accSeqOwner: ", accSeqOwner)
	t.Log("resOwnerPrivKey3: ", resOwnerPrivKey3)
	t.Log("resOwnerPubKey3: ", resOwnerPrivKey3.PubKey())

	accNumNode := resOwnerAcc3.GetAccountNumber()
	t.Log("accNumNode: ", accNumNode)
	accSeqNode := resOwnerAcc3.GetSequence()
	t.Log("accSeqNode: ", accSeqNode)

	gasInfo, result, e := mock.SignCheckDeliver(
		t,
		mApp.Cdc,
		mApp.BaseApp,
		header,
		[]sdk.Msg{registerResNodeMsg},
		[]uint64{accNumOwner, accNumNode},
		[]uint64{accSeqOwner, accSeqNode},
		true,
		true,
		resOwnerPrivKey3,
	)

	if e != nil {
		return
	}

	t.Log("gasInfo: ", gasInfo)
	t.Log("Result: ", result)
}
