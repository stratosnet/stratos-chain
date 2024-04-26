package app

import (
	"os"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibcfee "github.com/cosmos/ibc-go/v7/modules/apps/29-fee"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	ibctransfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcclientclient "github.com/cosmos/ibc-go/v7/modules/core/02-client/client"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	solomachine "github.com/cosmos/ibc-go/v7/modules/light-clients/06-solomachine"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"

	authmodulev1 "cosmossdk.io/api/cosmos/auth/module/v1"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	evmclient "github.com/stratosnet/stratos-chain/x/evm/client"

	"github.com/stratosnet/stratos-chain/x/evm"
	evmtypes "github.com/stratosnet/stratos-chain/x/evm/types"
	"github.com/stratosnet/stratos-chain/x/pot"
	pottypes "github.com/stratosnet/stratos-chain/x/pot/types"
	"github.com/stratosnet/stratos-chain/x/register"
	registertypes "github.com/stratosnet/stratos-chain/x/register/types"
	"github.com/stratosnet/stratos-chain/x/sds"
	sdstypes "github.com/stratosnet/stratos-chain/x/sds/types"
)

const (
	appName = "stchain"
)

var (
	// DefaultNodeHome sets the folder where the application data and configuration will be stored
	DefaultNodeHome = os.ExpandEnv("$HOME/.stchaind")
	powerReduction  = sdkmath.NewInt(1e18)

	// ModuleBasics is in charge of setting up basic module elements
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(
			[]govclient.ProposalHandler{
				paramsclient.ProposalHandler,
				upgradeclient.LegacyProposalHandler,
				upgradeclient.LegacyCancelProposalHandler,
				ibcclientclient.UpdateClientProposalHandler,
				ibcclientclient.UpgradeProposalHandler,
				evmclient.EVMChangeProxyImplementationHandler,
			},
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		vesting.AppModuleBasic{},
		consensus.AppModuleBasic{},

		// IBC modules
		ibc.AppModuleBasic{},
		solomachine.AppModuleBasic{},
		ibctm.AppModuleBasic{},
		ibctransfer.AppModuleBasic{},
		ibcfee.AppModuleBasic{},
		ica.AppModuleBasic{},

		// Stratos modules
		register.AppModuleBasic{},
		pot.AppModuleBasic{},
		sds.AppModuleBasic{},
		evm.AppModuleBasic{},
	)

	maccPerms = []*authmodulev1.ModuleAccountPermission{
		{Account: authtypes.FeeCollectorName},
		{Account: distrtypes.ModuleName, Permissions: []string{authtypes.Burner}},
		{Account: minttypes.ModuleName, Permissions: []string{authtypes.Minter}},
		{Account: stakingtypes.BondedPoolName, Permissions: []string{authtypes.Burner, authtypes.Staking}},
		{Account: stakingtypes.NotBondedPoolName, Permissions: []string{authtypes.Burner, authtypes.Staking}},
		{Account: govtypes.ModuleName, Permissions: []string{authtypes.Burner}},
		{Account: ibctransfertypes.ModuleName, Permissions: []string{authtypes.Minter, authtypes.Burner}},
		{Account: ibcfeetypes.ModuleName},
		{Account: icatypes.ModuleName},

		{Account: registertypes.ModuleName, Permissions: []string{authtypes.Minter, authtypes.Burner}},
		{Account: registertypes.ResourceNodeBondedPool, Permissions: []string{authtypes.Minter}},
		{Account: registertypes.ResourceNodeNotBondedPool, Permissions: []string{authtypes.Minter}},
		{Account: registertypes.MetaNodeBondedPool, Permissions: []string{authtypes.Minter}},
		{Account: registertypes.MetaNodeNotBondedPool, Permissions: []string{authtypes.Minter}},
		{Account: registertypes.TotalUnissuedPrepay, Permissions: []string{authtypes.Minter}},

		{Account: pottypes.ModuleName, Permissions: []string{authtypes.Minter}},
		{Account: pottypes.FoundationAccount, Permissions: []string{authtypes.Minter, authtypes.Burner}},
		{Account: pottypes.TotalRewardPool},

		{Account: sdstypes.ModuleName},
		{Account: evmtypes.ModuleName, Permissions: []string{authtypes.Minter, authtypes.Burner}},
	}

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: upgrade module must go first to handle software upgrades.
	// NOTE: staking module is required if HistoricalEntries param > 0.
	// NOTE: capability module's beginblocker must come before any modules using capabilities (e.g. IBC)
	beginBlockerOrder = []string{
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		evmtypes.ModuleName,
		minttypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		consensustypes.ModuleName,

		// IBC modules
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		ibcfeetypes.ModuleName,
		icatypes.ModuleName,

		// Stratos modules
		registertypes.ModuleName,
		pottypes.ModuleName,
		sdstypes.ModuleName,
	}

	endBlockerOrder = []string{
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		registertypes.ModuleName,
		sdstypes.ModuleName,
		evmtypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		pottypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		consensustypes.ModuleName,

		// IBC modules
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		ibcfeetypes.ModuleName,
		icatypes.ModuleName,
	}

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: The genutils module must also occur after auth so that it can access the params from auth.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	genesisModuleOrder = []string{
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		consensustypes.ModuleName,

		// IBC modules
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		ibcfeetypes.ModuleName,
		icatypes.ModuleName,

		// Stratos modules
		registertypes.ModuleName,
		pottypes.ModuleName,
		sdstypes.ModuleName,
		evmtypes.ModuleName,

		// NOTE: crisis module must go at the end to check for invariants on each module
		crisistypes.ModuleName,
	}

	// NOTE: The auth module must occur before everyone else. All other modules can be sorted
	// alphabetically (default order)
	// NOTE: The relationships module must occur before the profiles module, or all relationships will be deleted
	migrationModuleOrder = []string{
		authtypes.ModuleName,
		authz.ModuleName,
		banktypes.ModuleName,
		capabilitytypes.ModuleName,
		distrtypes.ModuleName,
		evidencetypes.ModuleName,
		feegrant.ModuleName,
		genutiltypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		slashingtypes.ModuleName,
		stakingtypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		consensustypes.ModuleName,

		// IBC modules
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		ibcfeetypes.ModuleName,
		icatypes.ModuleName,

		// Stratos modules
		registertypes.ModuleName,
		pottypes.ModuleName,
		sdstypes.ModuleName,
		evmtypes.ModuleName,

		crisistypes.ModuleName,
	}
)

func init() {
	version.AppName = appName + "d"

	//reset DefaultPowerReduction to prevent voting power overflow.
	sdk.DefaultPowerReduction = powerReduction
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for _, perms := range maccPerms {
		dupMaccPerms[perms.Account] = perms.Permissions
	}
	return dupMaccPerms
}
