package v0_11_0

import (
	"bytes"
	"context"
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

func MigrateStore(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.Codec, aminoCodec *codec.LegacyAmino) error {
	if err := migrateTotalRewardStore(ctx, storeKey, cdc, aminoCodec); err != nil {
		return err
	}
	return nil
}

func migrateTotalRewardStore(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.Codec, aminoCodec *codec.LegacyAmino) error {
	store := ctx.KVStore(storeKey)
	lastDistributedEpochBz := store.Get(types.LastDistributedEpochKeyPrefix)
	intValue := stratos.Int{}
	cdc.MustUnmarshalLengthPrefixed(lastDistributedEpochBz, &intValue)
	lastDistributedEpoch := *intValue.Value

	// init clientCtx
	clientCtx := client.Context{}.WithViper("")
	clientCtx, err := config.ReadFromClientConfig(clientCtx)
	if err != nil {
		return err
	}

	node, err := clientCtx.GetNode()
	if err != nil {
		return err
	}

	var matureEpochParam int64
	paramBz := store.Get(types.KeyMatureEpoch)
	if err = aminoCodec.UnmarshalJSON(paramBz, &matureEpochParam); err != nil {
		panic(err)
	}

	latestHeight := ctx.BlockHeight()
	currentEpoch := sdk.ZeroInt()
	for {
		miningRewardTotal := sdk.Coins{}
		trafficRewardTotal := sdk.Coins{}

		currentEpoch.Add(sdk.OneInt())
		if currentEpoch.GT(lastDistributedEpoch) {
			break
		}

		// query volumeReportRecord to get tx Hash
		var volumeReport types.VolumeReportRecord
		ctx = ctx.WithBlockHeight(latestHeight)
		store = ctx.KVStore(storeKey)
		volumeReportEpochbz := store.Get(types.VolumeReportStoreKey(currentEpoch))
		if volumeReportEpochbz == nil {
			continue
		}
		cdc.MustUnmarshalLengthPrefixed(volumeReportEpochbz, &volumeReport)
		volumeReportTxHash, err := hex.DecodeString(volumeReport.TxHash)
		if err != nil {
			continue
		}

		// query Tx to get the height of volumeReport tx execution
		resTx, err := node.Tx(context.Background(), volumeReportTxHash, true)
		if err != nil {
			return err
		}
		txHeight := resTx.Height

		// iterator individual at the height of volumeReport tx is executed
		ctx = ctx.WithBlockHeight(txHeight)
		store = ctx.KVStore(storeKey)
		iterator := store.Iterator(nil, nil)
		defer iterator.Close()
		for ; iterator.Valid(); iterator.Next() {
			key := iterator.Key()
			keyPrefixWithEpoch := bytes.Split(key, []byte("_"))[0]

			bEpoch := []byte(currentEpoch.AddRaw(matureEpochParam).String())
			matureEpochkey := append(types.IndividualRewardKeyPrefix, bEpoch...)

			if bytes.Equal(keyPrefixWithEpoch, matureEpochkey) {
				var individualReward types.Reward
				cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &individualReward)
				miningRewardTotal = miningRewardTotal.Add(individualReward.RewardFromMiningPool...)
				trafficRewardTotal = trafficRewardTotal.Add(individualReward.RewardFromTrafficPool...)
			}
		}

		//save totalReward by current epoch
		totalReward := types.TotalReward{
			MiningReward:  miningRewardTotal,
			TrafficReward: trafficRewardTotal,
		}
		totalRewardbz := cdc.MustMarshalLengthPrefixed(&totalReward)
		store.Set(types.GetTotalRewardKey(currentEpoch), totalRewardbz)
	}

	return nil
}
