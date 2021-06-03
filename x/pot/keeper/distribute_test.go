package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	"github.com/stratosnet/stratos-chain/x/register"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	ownerAddr1          = "st1qr9set2jaayzjjpm9tw4f3n6f5zfu3hef8wtaw"
	ownerAddr2          = "st1lc9sg3wq7guvkqv2d8vd2wycvj9fsspxq6qtg3"
	ownerAddr3          = "st15lm2e9h79j4d2zhyyf99j40uuy5vr404vurd0e"
	ownerAddr4          = "st1sysfc0hrjt63zqdywtzu7wu367uc23mcq7cz99"
	ownerAddr5          = ""
	ownerAddr6          = ""
	resourceNodePubKey1 = ""
	resourceNodePubKey2 = ""
	resourceNodePubKey3 = ""
	indexingNodePubKey1 = ""
	indexingNodePubKey2 = ""
	indexingNodePubKey3 = ""

	resourceNodeAddr1   = ""
	resourceNodeAddr2   = ""
	resourceNodeAddr3   = ""
	resourceNodeVolume1 = 10000000
	resourceNodeVolume2 = 20000000
	resourceNodeVolume3 = 30000000
	epoch1              = 1

	indexingNodeAddr1 = ""
	indexingNodeAddr2 = ""
	indexingNodeAddr3 = ""
)

func Test(t *testing.T) {
	ctx, _, _, potKeeper, _, _, _, registerKeeper := CreateTestInput(t, false)

	registerHandler := register.NewHandler(registerKeeper)
	//msg := staking.NewMsgCreateValidator(valOpAddr1, valConsPk1,
	//	sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100)), staking.Description{}, commission, sdk.OneInt())

	owner1, err := sdk.AccAddressFromBech32(ownerAddr1)
	require.NoError(t, err)
	owner2, err := sdk.AccAddressFromBech32(ownerAddr2)
	require.NoError(t, err)
	owner3, err := sdk.AccAddressFromBech32(ownerAddr3)
	require.NoError(t, err)
	owner4, err := sdk.AccAddressFromBech32(ownerAddr4)
	require.NoError(t, err)
	owner5, err := sdk.AccAddressFromBech32(ownerAddr5)
	require.NoError(t, err)
	owner6, err := sdk.AccAddressFromBech32(ownerAddr6)
	require.NoError(t, err)
	resourceNode1, err := sdk.AccAddressFromBech32(resourceNodeAddr1)
	require.NoError(t, err)
	resourceNode2, err := sdk.AccAddressFromBech32(resourceNodeAddr2)
	require.NoError(t, err)
	resourceNode3, err := sdk.AccAddressFromBech32(resourceNodeAddr3)
	require.NoError(t, err)

	pubKey1, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, resourceNodePubKey1)
	pubKey2, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, resourceNodePubKey2)
	pubKey3, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, resourceNodePubKey3)
	pubKey4, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, indexingNodePubKey1)
	pubKey5, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, indexingNodePubKey2)
	pubKey6, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, indexingNodePubKey3)

	//register resource node1
	msgRes := register.NewMsgCreateResourceNode("sds://resourceNode1", pubKey1, sdk.NewCoin(registerKeeper.BondDenom(ctx), sdk.NewInt(10000000)), owner1, register.Description{})
	res, err := registerHandler(ctx, msgRes)
	require.NoError(t, err)
	require.NotNil(t, res)

	//register resource node2
	msgRes = register.NewMsgCreateResourceNode("sds://resourceNode2", pubKey2, sdk.NewCoin(registerKeeper.BondDenom(ctx), sdk.NewInt(10000000)), owner2, register.Description{})
	res, err = registerHandler(ctx, msgRes)
	require.NoError(t, err)
	require.NotNil(t, res)

	//register resource node3
	msgRes = register.NewMsgCreateResourceNode("sds://resourceNode3", pubKey3, sdk.NewCoin(registerKeeper.BondDenom(ctx), sdk.NewInt(10000000)), owner3, register.Description{})
	res, err = registerHandler(ctx, msgRes)
	require.NoError(t, err)
	require.NotNil(t, res)

	//register indexing node1
	msgIdx := register.NewMsgCreateIndexingNode("sds://indexingNode1", pubKey4, sdk.NewCoin(registerKeeper.BondDenom(ctx), sdk.NewInt(10000000)), owner4, register.Description{})
	res, err = registerHandler(ctx, msgIdx)
	require.NoError(t, err)
	require.NotNil(t, res)

	//register indexing node2
	msgIdx = register.NewMsgCreateIndexingNode("sds://indexingNode2", pubKey5, sdk.NewCoin(registerKeeper.BondDenom(ctx), sdk.NewInt(10000000)), owner5, register.Description{})
	res, err = registerHandler(ctx, msgIdx)
	require.NoError(t, err)
	require.NotNil(t, res)

	//register indexing node3
	msgIdx = register.NewMsgCreateIndexingNode("sds://indexingNode3", pubKey6, sdk.NewCoin(registerKeeper.BondDenom(ctx), sdk.NewInt(10000000)), owner6, register.Description{})
	res, err = registerHandler(ctx, msgIdx)
	require.NoError(t, err)
	require.NotNil(t, res)

	var trafficList []types.SingleNodeVolume
	trafficList = append(trafficList, types.NewSingleNodeVolume(resourceNode1, sdk.NewInt(resourceNodeVolume1)))
	trafficList = append(trafficList, types.NewSingleNodeVolume(resourceNode2, sdk.NewInt(resourceNodeVolume2)))
	trafficList = append(trafficList, types.NewSingleNodeVolume(resourceNode3, sdk.NewInt(resourceNodeVolume3)))

	potKeeper.DistributePotReward(ctx, trafficList, sdk.NewInt(epoch1))
}
