package keeper

import (
	"sort"

	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// Iteration for sorting map to slice
func sortDetailMapToSlice(rewardDetailMap map[string]types.Reward) (rewardDetailList []types.Reward) {
	keys := make([]string, 0, len(rewardDetailMap))
	for key := range rewardDetailMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		reward := rewardDetailMap[key]
		rewardDetailList = append(rewardDetailList, reward)
	}
	return rewardDetailList
}
