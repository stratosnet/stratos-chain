package rest

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"
	"github.com/stratosnet/stratos-chain/x/pot/keeper"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/pot/rewards/epoch/{epoch}", getPotRewardsByEpochHandlerFn(cliCtx, keeper.QueryPotRewardsByEpoch)).Methods("GET")
	r.HandleFunc("/pot/rewards/owner/{ownerAddress}", getPotRewardsHandlerFn(cliCtx, keeper.QueryPotRewardsByOwner)).Methods("GET")
	r.HandleFunc("/pot/report/epoch/{epoch}", getVolumeReportHandlerFn(cliCtx, keeper.QueryVolumeReport)).Methods("GET")
}

func getPotRewardsByEpochHandlerFn(cliCtx context.CLIContext, queryPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get and verify params
		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		epochStr := mux.Vars(r)["epoch"]
		epoch, ok := checkEpoch(w, r, epochStr)
		if !ok {
			return
		}

		ownerAddressStr := ""
		if v := r.URL.Query().Get(RestOwnerAddress); len(v) != 0 {
			ownerAddressStr = v
		}
		ownerAddress, err := sdk.AccAddressFromBech32(ownerAddressStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// get volumeReportRecord from the given epoch
		volumeReportRecord, queryHeight := getVolumeReport(w, cliCtx, epoch)
		if volumeReportRecord.TxHash == "" {
			cliCtx = cliCtx.WithHeight(queryHeight)
			rest.PostProcessResponse(w, cliCtx, fmt.Sprintf("no volume report at epoch: %s", epoch.String()))
			return
		}

		// get nodeVolumes from volumeReportRecord.TxHash
		reportMsg := getNodeVolumes(w, cliCtx, volumeReportRecord)
		if len(reportMsg.NodesVolume) == 0 {
			cliCtx = cliCtx.WithHeight(queryHeight)
			rest.PostProcessResponse(w, cliCtx, fmt.Sprintf("no nodesVolumes in volume report at epoch: %s", epoch.String()))
			return
		}

		// create params with reportMsg.NodesVolume
		params := keeper.NewQueryPotRewardsByEpochParams(page, limit, ownerAddress, epoch, reportMsg.NodesVolume)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			cliCtx = cliCtx.WithHeight(queryHeight)
			rest.PostProcessResponse(w, cliCtx, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, queryPath)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			cliCtx = cliCtx.WithHeight(queryHeight)
			rest.PostProcessResponse(w, cliCtx, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getNodeVolumes(w http.ResponseWriter, cliCtx context.CLIContext, volumeReportRecord types.QueryVolumeReportRecord) types.MsgVolumeReport {
	output, err := utils.QueryTx(cliCtx, volumeReportRecord.TxHash)
	if err != nil {
		if strings.Contains(err.Error(), volumeReportRecord.TxHash) {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return types.MsgVolumeReport{}
		}
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return types.MsgVolumeReport{}
	}

	if output.Empty() {
		rest.WriteErrorResponse(w, http.StatusNotFound, fmt.Sprintf("no transaction found with hash %s", volumeReportRecord.TxHash))
	}
	v := output.Tx.GetMsgs()[0]
	reportMsg := v.(types.MsgVolumeReport)
	return reportMsg
}

func getVolumeReport(w http.ResponseWriter, cliCtx context.CLIContext, epoch sdk.Int) (types.QueryVolumeReportRecord, int64) {
	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, keeper.QueryVolumeReport)
	volumeReportRecordBz, height, err := cliCtx.QueryWithData(route, []byte(epoch.String()))
	if err != nil {
		return types.QueryVolumeReportRecord{}, height
	}

	if bytes.Contains(volumeReportRecordBz, []byte("no volume report at epoch")) {
		return types.QueryVolumeReportRecord{}, height
	}
	var volumeReportRecord types.QueryVolumeReportRecord
	err = cliCtx.Codec.UnmarshalJSON(volumeReportRecordBz, &volumeReportRecord)
	if err != nil {
		return types.QueryVolumeReportRecord{}, height
	}
	return volumeReportRecord, height
}

// GET request handler to query Volume report info
func getVolumeReportHandlerFn(cliCtx context.CLIContext, queryPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		v := mux.Vars(r)["epoch"]
		if len(v) == 0 {
			return
		}
		epoch, ok := checkEpoch(w, r, v)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, queryPath)
		res, height, err := cliCtx.QueryWithData(route, []byte(epoch.String()))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func checkEpoch(w http.ResponseWriter, r *http.Request, epochStr string) (sdk.Int, bool) {
	epoch, ok := sdk.NewIntFromString(epochStr)
	if !ok {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid epoch")
		return sdk.NewInt(-1), false
	}
	return epoch, true
}

// GET request handler to query potRewards info by nodeWalletAddr
func getPotRewardsHandlerFn(cliCtx context.CLIContext, queryPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		ownerAddrStr := mux.Vars(r)["ownerAddress"]
		ownerAddr, err := sdk.AccAddressFromBech32(ownerAddrStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		var (
			queryHeight int64
		)

		if v := r.URL.Query().Get(RestHeight); len(v) != 0 {
			queryHeight, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		params := keeper.NewQueryPotRewardsWithOwnerHeightParams(page, limit, ownerAddr, queryHeight)

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(queryHeight)
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, queryPath)
		res, currentHeight, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(currentHeight)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
