package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stratosnet/stratos-chain/x/register/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

var (
	spNodeOwner1   = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	spNodeOwner2   = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	spNodeOwner3   = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	spNodeOwner4   = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	spNodeOwnerNew = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	spNodePubKey1 = ed25519.GenPrivKey().PubKey()
	spNodeAddr1   = sdk.AccAddress(spNodePubKey1.Address())
	initialStake1 = sdk.NewInt(100000000)

	spNodePubKey2 = ed25519.GenPrivKey().PubKey()
	spNodeAddr2   = sdk.AccAddress(spNodePubKey2.Address())
	initialStake2 = sdk.NewInt(100000000)

	spNodePubKey3 = ed25519.GenPrivKey().PubKey()
	spNodeAddr3   = sdk.AccAddress(spNodePubKey3.Address())
	initialStake3 = sdk.NewInt(100000000)

	spNodePubKey4 = ed25519.GenPrivKey().PubKey()
	spNodeAddr4   = sdk.AccAddress(spNodePubKey4.Address())
	initialStake4 = sdk.NewInt(100000000)

	spNodePubKeyNew = ed25519.GenPrivKey().PubKey()
	spNodeAddrNew   = sdk.AccAddress(spNodePubKeyNew.Address())
	spNodeStakeNew  = sdk.NewInt(100000000)
)

func TestExpiredVote(t *testing.T) {

	ctx, accountKeeper, bankKeeper, k, _ := CreateTestInput(t, false)

	//genesis init Sp nodes.
	genesisSpNode1 := types.NewIndexingNode("sds://indexingNode1", spNodePubKey1, spNodeOwner1, types.NewDescription("sds://indexingNode1", "", "", "", ""))
	genesisSpNode2 := types.NewIndexingNode("sds://indexingNode2", spNodePubKey2, spNodeOwner2, types.NewDescription("sds://indexingNode2", "", "", "", ""))
	genesisSpNode3 := types.NewIndexingNode("sds://indexingNode3", spNodePubKey3, spNodeOwner3, types.NewDescription("sds://indexingNode3", "", "", "", ""))
	genesisSpNode4 := types.NewIndexingNode("sds://indexingNode4", spNodePubKey4, spNodeOwner4, types.NewDescription("sds://indexingNode4", "", "", "", ""))
	genesisSpNode1.Tokens = genesisSpNode1.Tokens.Add(initialStake1)
	genesisSpNode2.Tokens = genesisSpNode2.Tokens.Add(initialStake2)
	genesisSpNode3.Tokens = genesisSpNode3.Tokens.Add(initialStake3)
	genesisSpNode4.Tokens = genesisSpNode3.Tokens.Add(initialStake4)
	genesisSpNode1.Status = sdk.Bonded
	genesisSpNode2.Status = sdk.Bonded
	genesisSpNode3.Status = sdk.Bonded
	genesisSpNode4.Status = sdk.Bonded

	k.SetIndexingNode(ctx, genesisSpNode1)
	k.SetIndexingNode(ctx, genesisSpNode2)
	k.SetIndexingNode(ctx, genesisSpNode3)
	k.SetIndexingNode(ctx, genesisSpNode4)
	k.SetLastIndexingNodeStake(ctx, spNodeAddr1, initialStake1)
	k.SetLastIndexingNodeStake(ctx, spNodeAddr2, initialStake2)
	k.SetLastIndexingNodeStake(ctx, spNodeAddr3, initialStake3)
	k.SetLastIndexingNodeStake(ctx, spNodeAddr4, initialStake4)
	k.SetIndexingNodeBondedToken(ctx, sdk.NewCoin(k.BondDenom(ctx), initialStake1.Add(initialStake2).Add(initialStake3).Add(initialStake4)))
	k.SetInitialGenesisStakeTotal(ctx, initialStake1.Add(initialStake2).Add(initialStake3).Add(initialStake4))

	//Register new SP node after genesis initialized
	createAccount(t, ctx, accountKeeper, bankKeeper, spNodeOwnerNew, sdk.NewCoins(sdk.NewCoin("ustos", spNodeStakeNew)))
	err := k.RegisterIndexingNode(ctx, "sds://newIndexingNode", spNodePubKeyNew, spNodeOwnerNew,
		types.NewDescription("sds://newIndexingNode", "", "", "", ""), sdk.NewCoin("ustos", spNodeStakeNew))
	require.NoError(t, err)

	//set expireTime of voting to 7 days before
	votePool, found := k.GetIndexingNodeRegistrationVotePool(ctx, spNodeAddrNew)
	require.True(t, found)
	require.NotNil(t, votePool)
	votePool.ExpireTime = votePool.ExpireTime.AddDate(0, 0, -7)
	k.SetIndexingNodeRegistrationVotePool(ctx, votePool)

	//After registration, the status of new SP node is UNBONDED
	_, found = k.GetIndexingNode(ctx, spNodeAddrNew)
	require.True(t, found)
	require.Error(t, types.ErrVoteExpired)
}

func TestDuplicateVote(t *testing.T) {

	ctx, accountKeeper, bankKeeper, k, _ := CreateTestInput(t, false)

	//genesis init Sp nodes.
	genesisSpNode1 := types.NewIndexingNode("sds://indexingNode1", spNodePubKey1, spNodeOwner1, types.NewDescription("sds://indexingNode1", "", "", "", ""))
	genesisSpNode2 := types.NewIndexingNode("sds://indexingNode2", spNodePubKey2, spNodeOwner2, types.NewDescription("sds://indexingNode2", "", "", "", ""))
	genesisSpNode3 := types.NewIndexingNode("sds://indexingNode3", spNodePubKey3, spNodeOwner3, types.NewDescription("sds://indexingNode3", "", "", "", ""))
	genesisSpNode4 := types.NewIndexingNode("sds://indexingNode4", spNodePubKey4, spNodeOwner4, types.NewDescription("sds://indexingNode4", "", "", "", ""))
	genesisSpNode1.Tokens = genesisSpNode1.Tokens.Add(initialStake1)
	genesisSpNode2.Tokens = genesisSpNode2.Tokens.Add(initialStake2)
	genesisSpNode3.Tokens = genesisSpNode3.Tokens.Add(initialStake3)
	genesisSpNode4.Tokens = genesisSpNode3.Tokens.Add(initialStake4)
	genesisSpNode1.Status = sdk.Bonded
	genesisSpNode2.Status = sdk.Bonded
	genesisSpNode3.Status = sdk.Bonded
	genesisSpNode4.Status = sdk.Bonded

	k.SetIndexingNode(ctx, genesisSpNode1)
	k.SetIndexingNode(ctx, genesisSpNode2)
	k.SetIndexingNode(ctx, genesisSpNode3)
	k.SetIndexingNode(ctx, genesisSpNode4)
	k.SetLastIndexingNodeStake(ctx, spNodeAddr1, initialStake1)
	k.SetLastIndexingNodeStake(ctx, spNodeAddr2, initialStake2)
	k.SetLastIndexingNodeStake(ctx, spNodeAddr3, initialStake3)
	k.SetLastIndexingNodeStake(ctx, spNodeAddr4, initialStake4)
	k.SetIndexingNodeBondedToken(ctx, sdk.NewCoin(k.BondDenom(ctx), initialStake1.Add(initialStake2).Add(initialStake3).Add(initialStake4)))
	k.SetInitialGenesisStakeTotal(ctx, initialStake1.Add(initialStake2).Add(initialStake3).Add(initialStake4))

	//Register new SP node after genesis initialized
	createAccount(t, ctx, accountKeeper, bankKeeper, spNodeOwnerNew, sdk.NewCoins(sdk.NewCoin("ustos", spNodeStakeNew)))
	err := k.RegisterIndexingNode(ctx, "sds://newIndexingNode", spNodePubKeyNew, spNodeOwnerNew,
		types.NewDescription("sds://newIndexingNode", "", "", "", ""), sdk.NewCoin("ustos", spNodeStakeNew))
	require.NoError(t, err)

	//After registration, the status of new SP node is UNBONDED
	newNode, found := k.GetIndexingNode(ctx, spNodeAddrNew)
	require.True(t, found)
	require.Equal(t, newNode.Status, sdk.Unbonded)

	//Exist SP Node1 vote to approve, the status of new SP node is UNBONDED
	err = handlerSimulate(ctx, k, types.NewMsgIndexingNodeRegistrationVote(spNodeAddrNew, spNodeOwnerNew, types.Approve, spNodeAddr1, spNodeOwner1))
	require.NoError(t, err)
	newNode, found = k.GetIndexingNode(ctx, spNodeAddrNew)
	require.True(t, found)
	require.Equal(t, newNode.Status, sdk.Unbonded)

	//Exist SP Node1 vote to approve, the status of new SP node is UNBONDED
	err = handlerSimulate(ctx, k, types.NewMsgIndexingNodeRegistrationVote(spNodeAddrNew, spNodeOwnerNew, types.Approve, spNodeAddr1, spNodeOwner1))
	require.Error(t, types.ErrDuplicateVoting)
}

func TestSpRegistrationApproval(t *testing.T) {

	ctx, accountKeeper, bankKeeper, k, _ := CreateTestInput(t, false)

	//genesis init Sp nodes.
	genesisSpNode1 := types.NewIndexingNode("sds://indexingNode1", spNodePubKey1, spNodeOwner1, types.NewDescription("sds://indexingNode1", "", "", "", ""))
	genesisSpNode2 := types.NewIndexingNode("sds://indexingNode2", spNodePubKey2, spNodeOwner2, types.NewDescription("sds://indexingNode2", "", "", "", ""))
	genesisSpNode3 := types.NewIndexingNode("sds://indexingNode3", spNodePubKey3, spNodeOwner3, types.NewDescription("sds://indexingNode3", "", "", "", ""))
	genesisSpNode4 := types.NewIndexingNode("sds://indexingNode4", spNodePubKey4, spNodeOwner4, types.NewDescription("sds://indexingNode4", "", "", "", ""))
	genesisSpNode1.Tokens = genesisSpNode1.Tokens.Add(initialStake1)
	genesisSpNode2.Tokens = genesisSpNode2.Tokens.Add(initialStake2)
	genesisSpNode3.Tokens = genesisSpNode3.Tokens.Add(initialStake3)
	genesisSpNode4.Tokens = genesisSpNode3.Tokens.Add(initialStake4)
	genesisSpNode1.Status = sdk.Bonded
	genesisSpNode2.Status = sdk.Bonded
	genesisSpNode3.Status = sdk.Bonded
	genesisSpNode4.Status = sdk.Bonded

	k.SetIndexingNode(ctx, genesisSpNode1)
	k.SetIndexingNode(ctx, genesisSpNode2)
	k.SetIndexingNode(ctx, genesisSpNode3)
	k.SetIndexingNode(ctx, genesisSpNode4)
	k.SetLastIndexingNodeStake(ctx, spNodeAddr1, initialStake1)
	k.SetLastIndexingNodeStake(ctx, spNodeAddr2, initialStake2)
	k.SetLastIndexingNodeStake(ctx, spNodeAddr3, initialStake3)
	k.SetLastIndexingNodeStake(ctx, spNodeAddr4, initialStake4)
	k.SetIndexingNodeBondedToken(ctx, sdk.NewCoin(k.BondDenom(ctx), initialStake1.Add(initialStake2).Add(initialStake3).Add(initialStake4)))
	k.SetInitialGenesisStakeTotal(ctx, initialStake1.Add(initialStake2).Add(initialStake3).Add(initialStake4))

	//Register new SP node after genesis initialized
	createAccount(t, ctx, accountKeeper, bankKeeper, spNodeOwnerNew, sdk.NewCoins(sdk.NewCoin("ustos", spNodeStakeNew)))
	//_, err := k.bankKeeper.AddCoins(ctx, spNodeAddr4, sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), sdk.NewInt(10000000000000))))
	//require.NoError(t, err)

	err := k.RegisterIndexingNode(ctx, "sds://newIndexingNode", spNodePubKeyNew, spNodeOwnerNew,
		types.NewDescription("sds://newIndexingNode", "", "", "", ""), sdk.NewCoin("ustos", spNodeStakeNew))
	require.NoError(t, err)

	//After registration, the status of new SP node is UNBONDED
	newNode, found := k.GetIndexingNode(ctx, spNodeAddrNew)
	require.True(t, found)
	require.Equal(t, newNode.Status, sdk.Unbonded)

	//Exist SP Node1 vote to approve, the status of new SP node is UNBONDED
	err = handlerSimulate(ctx, k, types.NewMsgIndexingNodeRegistrationVote(spNodeAddrNew, spNodeOwnerNew, types.Approve, spNodeAddr1, spNodeOwner1))
	require.NoError(t, err)
	newNode, found = k.GetIndexingNode(ctx, spNodeAddrNew)
	require.True(t, found)
	require.Equal(t, newNode.Status, sdk.Unbonded)

	//Exist SP Node2 vote to approve, the status of new SP node is UNBONDED
	err = handlerSimulate(ctx, k, types.NewMsgIndexingNodeRegistrationVote(spNodeAddrNew, spNodeOwnerNew, types.Approve, spNodeAddr2, spNodeOwner2))
	require.NoError(t, err)
	newNode, found = k.GetIndexingNode(ctx, spNodeAddrNew)
	require.True(t, found)
	require.Equal(t, newNode.Status, sdk.Unbonded)

	//Exist SP Node3 vote to approve, the status of new SP node changes to BONDED
	err = handlerSimulate(ctx, k, types.NewMsgIndexingNodeRegistrationVote(spNodeAddrNew, spNodeOwnerNew, types.Approve, spNodeAddr3, spNodeOwner3))
	require.NoError(t, err)
	newNode, found = k.GetIndexingNode(ctx, spNodeAddrNew)
	require.True(t, found)
	require.Equal(t, newNode.Status, sdk.Bonded)

	//Exist SP Node4 vote to approve, the status of new SP node is BONDED
	err = handlerSimulate(ctx, k, types.NewMsgIndexingNodeRegistrationVote(spNodeAddrNew, spNodeOwnerNew, types.Approve, spNodeAddr4, spNodeOwner4))
	require.NoError(t, err)
	newNode, found = k.GetIndexingNode(ctx, spNodeAddrNew)
	require.True(t, found)
	require.Equal(t, newNode.Status, sdk.Bonded)
}

func handlerSimulate(ctx sdk.Context, k Keeper, msg types.MsgIndexingNodeRegistrationVote) error {
	nodeToApprove, found := k.GetIndexingNode(ctx, msg.CandidateNetworkAddress)
	if !found {
		return types.ErrNoIndexingNodeFound
	}
	if !nodeToApprove.GetOwnerAddr().Equals(msg.CandidateOwnerAddress) {
		return types.ErrInvalidOwnerAddr
	}

	approver, found := k.GetIndexingNode(ctx, msg.VoterNetworkAddress)
	if !found {
		return types.ErrInvalidVoterAddr
	}
	if !approver.Status.Equal(sdk.Bonded) || approver.IsSuspended() {
		return types.ErrInvalidVoterStatus
	}

	_, err := k.HandleVoteForIndexingNodeRegistration(ctx, msg.CandidateNetworkAddress, msg.CandidateOwnerAddress, msg.Opinion, msg.VoterNetworkAddress)
	if err != nil {
		return err
	}

	return nil
}

func createAccount(t *testing.T, ctx sdk.Context, accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper, acc sdk.AccAddress, coins sdk.Coins) {
	account := accountKeeper.GetAccount(ctx, acc)
	if account == nil {
		account = accountKeeper.NewAccountWithAddress(ctx, acc)
		//fmt.Printf("create account: " + account.String() + "\n")
	}
	coins, err := bankKeeper.AddCoins(ctx, acc, coins)
	require.NoError(t, err)
}
