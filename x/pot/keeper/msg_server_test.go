package keeper_test

import (
	"fmt"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stratosnet/stratos-chain/crypto/merkle"
	stratostestutil "github.com/stratosnet/stratos-chain/testutil"
	sttypes "github.com/stratosnet/stratos-chain/types"
	pottypes "github.com/stratosnet/stratos-chain/x/pot/types"
	regtypes "github.com/stratosnet/stratos-chain/x/register/types"
)

func (s *KeeperTestSuite) fakeAndMockMsgVolumeReportForMerkleTest(epochN int64, mpd merkle.MerkleProofData) *pottypes.MsgVolumeReport {
	epoch := sdkmath.NewInt(epochN)
	ownerAcc := sdk.AccAddress([]byte("ownerAcc1_______________"))
	walletAcc := sdk.AccAddress([]byte("walletAcc1_______________"))
	reporterAcc := sttypes.SdsAddress([]byte("reporterAcc1_______________"))
	reporterOwnerAcc := sdk.AccAddress([]byte("reporterOAcc1_______________"))
	walletVolume := pottypes.NewSingleWalletVolume(walletAcc, sdkmath.NewInt(10))

	metaNodeP2PPrivKey := ed25519.GenPrivKey()

	msg := pottypes.NewMsgVolumeReport(
		[]pottypes.SingleWalletVolume{walletVolume},
		reporterAcc,
		epoch,
		"test",
		reporterOwnerAcc,
	)

	metaNodeData := make(map[string][]byte)
	metaNode, _ := regtypes.NewMetaNode(
		sttypes.SdsAddress(ownerAcc.String()),
		metaNodeP2PPrivKey.PubKey(),
		ownerAcc,
		ownerAcc,
		regtypes.NewDescription("test", "test", "test", "test", "test"),
		time.Now(),
	)
	metaNodeData["1"] = regtypes.MustMarshalMetaNode(s.encCfg.Codec, metaNode)
	mockMetaNodeIterator := stratostestutil.NewMockIterator(metaNodeData)

	msg.MerkleProofData = pottypes.MerkleProofData{
		Root:        mpd.GetRoot(),
		Commitments: mpd.GetCommitments(),
	}

	msg, _ = stratostestutil.SignVolumeReport(
		msg,
		metaNodeP2PPrivKey.Bytes(),
	)

	// mocks start
	s.mockRegOwnMetaNode(reporterOwnerAcc, reporterAcc).Return(true).AnyTimes()
	s.mockRegGetBondedMetaNodeCnt().Return(sdkmath.NewInt(1)).AnyTimes()
	s.mockRegGetRemainingOzoneLimit().Return(sdkmath.NewInt(10)).AnyTimes()
	s.mockRegGetDepositNozRate().Return(sdkmath.LegacyNewDec(1)).AnyTimes()
	s.mockRegGetEffectiveTotalDeposit().Return(sdkmath.NewInt(20)).AnyTimes()
	s.mockRegGetTotalUnissuedPrepay().Return(sdk.NewCoin(pottypes.DefaultBondDenom, sdk.NewInt(10))).AnyTimes()
	s.mockRegGetMetaNodeIterator().Return(mockMetaNodeIterator).AnyTimes()
	s.mockRegSetRemainingOzoneLimit(sdkmath.NewInt(20)).AnyTimes()

	m2mAmt1, _ := sdk.NewIntFromString("16000000000000000000")
	s.mockBankSendCoinsFromModuleToModule(
		pottypes.FoundationAccount, authtypes.FeeCollectorName,
		sdk.NewCoins(sdk.NewCoin(pottypes.DefaultBondDenom, m2mAmt1)),
	).Return(nil).AnyTimes()

	m2mAmt2, _ := sdk.NewIntFromString("3")
	s.mockBankSendCoinsFromModuleToModule(
		regtypes.TotalUnissuedPrepay, authtypes.FeeCollectorName,
		sdk.NewCoins(sdk.NewCoin(pottypes.DefaultBondDenom, m2mAmt2)),
	).Return(nil).AnyTimes()

	m2mAmt3, _ := sdk.NewIntFromString("48000000000000000000")
	s.mockBankSendCoinsFromModuleToModule(
		pottypes.FoundationAccount, pottypes.TotalRewardPool,
		sdk.NewCoins(sdk.NewCoin(pottypes.DefaultBondDenom, m2mAmt3)),
	).Return(nil).AnyTimes()

	m2mAmt4, _ := sdk.NewIntFromString("8")
	s.mockBankSendCoinsFromModuleToModule(
		regtypes.TotalUnissuedPrepay, pottypes.TotalRewardPool,
		sdk.NewCoins(sdk.NewCoin(pottypes.DefaultBondDenom, m2mAmt4)),
	).Return(nil).AnyTimes()

	s.mockAccGetModuleAddress(regtypes.TotalUnissuedPrepay).Return(prepAcct).AnyTimes()
	s.mockDistrFundCommunityPool(sdk.NewCoins(sdk.NewCoin(pottypes.DefaultBondDenom, sdkmath.NewInt(0))), prepAcct).Return(nil).AnyTimes()
	// mocks end

	return msg
}

func (s *KeeperTestSuite) TestMsgVolumeReport() {
	require := s.Require()
	epochCounter := int64(0)

	s.T().Run("check merkle proof data and ok", func(t *testing.T) {
		s.Reset(t)
		epochCounter++

		ctx, msgServer := s.ctx, s.msgServer
		regKeeper := s.GetNoMockRegKeeper()
		proover := merkle.NewRelayerMerkleProver()
		regKeeper.SetProover(proover)

		// proof creation
		commitment1 := []byte("c1")

		root := regKeeper.GetMerkleRoot(ctx)

		mpd, _ := proover.CreateProofs(
			[][]byte{
				root,
			},
			[][]byte{
				commitment1,
			})

		// create commitment to process
		err := regKeeper.CreateMerkleCommitment(ctx, commitment1, root)
		require.NoError(err)

		msg := s.fakeAndMockMsgVolumeReportForMerkleTest(epochCounter, mpd)

		mp := msg.GetMerkleProofData()
		// real execute
		s.mockRegProcessMerkleProofs(&mp).DoAndReturn(regKeeper.ProcessMerkleProofs).AnyTimes()

		_, err = msgServer.HandleMsgVolumeReport(ctx, msg)
		require.NoError(err)
		require.Equal(2, len(ctx.EventManager().ABCIEvents()))

		evt1 := ctx.EventManager().ABCIEvents()[0]
		msg1, _ := sdk.ParseTypedEvent(evt1)
		revt1 := msg1.(*regtypes.EventMerkleDataUpdated)
		require.Equal(root, revt1.Root)
		require.Equal(commitment1, revt1.Commitment)
		require.Equal(regtypes.EventMerkleDataUpdated_NULLIFY, revt1.ActionType)
	})

	s.T().Run("same/not exist proceed commitment and fail", func(t *testing.T) {
		s.Reset(t)
		epochCounter++

		ctx, msgServer := s.ctx, s.msgServer
		regKeeper := s.GetNoMockRegKeeper()
		proover := merkle.NewRelayerMerkleProver()
		regKeeper.SetProover(proover)

		// proof creation
		commitment1 := []byte("c1")

		root := regKeeper.GetMerkleRoot(ctx)

		mpd, _ := proover.CreateProofs(
			[][]byte{
				root,
			},
			[][]byte{
				commitment1,
			})

		msg := s.fakeAndMockMsgVolumeReportForMerkleTest(epochCounter, mpd)

		mp := msg.GetMerkleProofData()
		// real execute
		s.mockRegProcessMerkleProofs(&mp).DoAndReturn(regKeeper.ProcessMerkleProofs).AnyTimes()

		_, err := msgServer.HandleMsgVolumeReport(ctx, msg)
		require.Error(err)
		require.Contains(err.Error(), fmt.Sprintf("commitment '%s' does not exist", commitment1))
		require.Equal(0, len(ctx.EventManager().ABCIEvents()))
	})
}
