package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stratosnet/stratos-chain/crypto/merkle"
	regtypes "github.com/stratosnet/stratos-chain/x/register/types"
	sdstypes "github.com/stratosnet/stratos-chain/x/sds/types"
)

func (s *KeeperTestSuite) fakeAndMockMsgPrepayForMerkleTest(sender sdk.AccAddress) (*sdstypes.MsgPrepay, []byte) {
	ozonePrice := sdkmath.NewInt(2)
	remainingOzone := sdkmath.NewInt(100)

	amount := sdk.NewCoin(sdstypes.DefaultBondDenom, sdk.NewInt(1))

	msg := &sdstypes.MsgPrepay{
		Sender:      sender.String(),
		Beneficiary: sender.String(),
		Amount:      sdk.NewCoins(amount),
	}

	tev := &sdstypes.EventPrePay{
		Sender:       msg.GetSender(),
		Beneficiary:  msg.GetBeneficiary(),
		Amount:       msg.GetAmount().String(),
		PurchasedNoz: ozonePrice.String(),
	}

	data, _ := tev.Marshal()

	s.mockBankHasBalance(sender, amount).Return(true).AnyTimes()
	s.mockRegCalculatePurchaseAmount(amount.Amount).Return(ozonePrice, remainingOzone, nil).AnyTimes()
	s.mockBankSendCoinsFromAccountToModule(sender, regtypes.TotalUnissuedPrepay, sdk.NewCoins(amount)).Return(nil).AnyTimes()
	s.mockRegSetRemainingOzoneLimit(remainingOzone).AnyTimes()

	acc := authtypes.NewBaseAccount(sender, nil, 1, 1)
	s.mockAccGetAccount(sender).Return(acc).AnyTimes()

	return msg, data
}

func (s *KeeperTestSuite) TestMsgPrepay() {
	require := s.Require()

	sender := sdk.AccAddress([]byte("addr1_______________"))

	s.T().Run("create merkle proof with data and ok", func(t *testing.T) {
		s.Reset(t)

		ctx, msgServer := s.ctx, s.msgServer

		regKeeper := s.GetNoMockRegKeeper()
		proover := merkle.NewRelayerMerkleProver()
		regKeeper.SetProover(proover)

		msg, data := s.fakeAndMockMsgPrepayForMerkleTest(sender)
		s.mockRegGenerateMerkleProofs(sender, data).DoAndReturn(regKeeper.GenerateMerkleProofs).AnyTimes()

		root := regKeeper.GetMerkleRoot(ctx)
		require.Equal(merkle.NullCommitment[:], root)

		commitment1 := merkle.CreateSdkCommitment(sender, 1, data)

		newRoot := proover.GetRoot(root, [][]byte{commitment1})

		_, err := msgServer.HandleMsgPrepay(ctx, msg)
		require.NoError(err)
		require.Equal(2, len(ctx.EventManager().ABCIEvents()))

		evt1 := ctx.EventManager().ABCIEvents()[0]
		msg1, _ := sdk.ParseTypedEvent(evt1)
		revt1 := msg1.(*regtypes.EventMerkleDataUpdated)
		require.Equal(newRoot, revt1.Root)
		require.Equal(commitment1, revt1.Commitment)
	})
}
