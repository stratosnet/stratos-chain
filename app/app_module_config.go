package app

import (
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"

	runtimev1alpha1 "cosmossdk.io/api/cosmos/app/runtime/v1alpha1"
	appv1alpha1 "cosmossdk.io/api/cosmos/app/v1alpha1"
	authmodulev1 "cosmossdk.io/api/cosmos/auth/module/v1"
	authzmodulev1 "cosmossdk.io/api/cosmos/authz/module/v1"
	bankmodulev1 "cosmossdk.io/api/cosmos/bank/module/v1"
	capabilitymodulev1 "cosmossdk.io/api/cosmos/capability/module/v1"
	consensusmodulev1 "cosmossdk.io/api/cosmos/consensus/module/v1"
	crisismodulev1 "cosmossdk.io/api/cosmos/crisis/module/v1"
	distrmodulev1 "cosmossdk.io/api/cosmos/distribution/module/v1"
	evidencemodulev1 "cosmossdk.io/api/cosmos/evidence/module/v1"
	feegrantmodulev1 "cosmossdk.io/api/cosmos/feegrant/module/v1"
	genutilmodulev1 "cosmossdk.io/api/cosmos/genutil/module/v1"
	govmodulev1 "cosmossdk.io/api/cosmos/gov/module/v1"
	mintmodulev1 "cosmossdk.io/api/cosmos/mint/module/v1"
	paramsmodulev1 "cosmossdk.io/api/cosmos/params/module/v1"
	slashingmodulev1 "cosmossdk.io/api/cosmos/slashing/module/v1"
	stakingmodulev1 "cosmossdk.io/api/cosmos/staking/module/v1"
	txconfigv1 "cosmossdk.io/api/cosmos/tx/config/v1"
	upgrademodulev1 "cosmossdk.io/api/cosmos/upgrade/module/v1"
	vestingmodulev1 "cosmossdk.io/api/cosmos/vesting/module/v1"
	sdkappconfig "cosmossdk.io/core/appconfig"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	evmmodulev1 "github.com/stratosnet/stratos-chain/api/stratos/evm/module/v1"
	potmodulev1 "github.com/stratosnet/stratos-chain/api/stratos/pot/module/v1"
	registermodulev1 "github.com/stratosnet/stratos-chain/api/stratos/register/module/v1"
	sdsmodulev1 "github.com/stratosnet/stratos-chain/api/stratos/sds/module/v1"
	stratos "github.com/stratosnet/stratos-chain/types"
	evmtypes "github.com/stratosnet/stratos-chain/x/evm/types"
	pottypes "github.com/stratosnet/stratos-chain/x/pot/types"
	registertypes "github.com/stratosnet/stratos-chain/x/register/types"
	sdstypes "github.com/stratosnet/stratos-chain/x/sds/types"
)

var (
	// blocked account addresses
	blockAccAddrs = []string{
		govtypes.ModuleName,
		minttypes.ModuleName,
		stakingtypes.BondedPoolName,
		stakingtypes.NotBondedPoolName,
		ibctransfertypes.ModuleName,
		ibcfeetypes.ModuleName,
		icatypes.ModuleName,
		registertypes.ModuleName,
		registertypes.ResourceNodeBondedPool,
		registertypes.ResourceNodeNotBondedPool,
		registertypes.MetaNodeBondedPool,
		registertypes.MetaNodeNotBondedPool,
		registertypes.TotalUnissuedPrepay,
		pottypes.ModuleName,
		pottypes.FoundationAccount,
		pottypes.TotalRewardPool,
		sdstypes.ModuleName,
		evmtypes.ModuleName,

		// We allow the following module accounts to receive funds:
		// distrtypes.ModuleName
		// authtypes.FeeCollectorName
	}

	// application configuration (used by depinject)
	StratosAppConfig = sdkappconfig.Compose(&appv1alpha1.Config{
		Modules: []*appv1alpha1.ModuleConfig{
			// SDK modules
			{
				Name: "runtime",
				Config: sdkappconfig.WrapAny(&runtimev1alpha1.Module{
					AppName: appName,
					// During begin block slashing happens after distr.BeginBlocker so that
					// there is nothing left over in the validator fee pool, so as to keep the
					// CanWithdrawInvariant invariant.
					// NOTE: staking module is required if HistoricalEntries param > 0
					// NOTE: capability module's beginblocker must come before any modules using capabilities (e.g. IBC)
					BeginBlockers: beginBlockerOrder,
					EndBlockers:   endBlockerOrder,
					OverrideStoreKeys: []*runtimev1alpha1.StoreKeyConfig{
						{
							ModuleName: authtypes.ModuleName,
							KvStoreKey: "acc",
						},
					},
					InitGenesis: genesisModuleOrder,
					// When ExportGenesis is not specified, the export genesis module order
					// is equal to the init genesis order
					// ExportGenesis: genesisModuleOrder,
					// Uncomment if you want to set a custom migration order here.
					// OrderMigrations: nil,
					OrderMigrations: migrationModuleOrder,
				}),
			},
			{
				Name: authtypes.ModuleName,
				Config: sdkappconfig.WrapAny(&authmodulev1.Module{
					Bech32Prefix:             stratos.AccountAddressPrefix,
					ModuleAccountPermissions: maccPerms,
					// By default, modules authority is the governance module. This is configurable with the following:
					// Authority: "group", // A custom module authority can be set using a module name
					// Authority: "cosmos1cwwv22j5ca08ggdv9c2uky355k908694z577tv", // or a specific address
				}),
			},
			{
				Name:   vestingtypes.ModuleName,
				Config: sdkappconfig.WrapAny(&vestingmodulev1.Module{}),
			},
			{
				Name: banktypes.ModuleName,
				Config: sdkappconfig.WrapAny(&bankmodulev1.Module{
					BlockedModuleAccountsOverride: blockAccAddrs,
				}),
			},
			{
				Name:   stakingtypes.ModuleName,
				Config: sdkappconfig.WrapAny(&stakingmodulev1.Module{}),
			},
			{
				Name:   slashingtypes.ModuleName,
				Config: sdkappconfig.WrapAny(&slashingmodulev1.Module{}),
			},
			{
				Name:   paramstypes.ModuleName,
				Config: sdkappconfig.WrapAny(&paramsmodulev1.Module{}),
			},
			{
				Name:   "tx",
				Config: sdkappconfig.WrapAny(&txconfigv1.Config{}),
			},
			{
				Name:   genutiltypes.ModuleName,
				Config: sdkappconfig.WrapAny(&genutilmodulev1.Module{}),
			},
			{
				Name:   authz.ModuleName,
				Config: sdkappconfig.WrapAny(&authzmodulev1.Module{}),
			},
			{
				Name:   upgradetypes.ModuleName,
				Config: sdkappconfig.WrapAny(&upgrademodulev1.Module{}),
			},
			{
				Name:   distrtypes.ModuleName,
				Config: sdkappconfig.WrapAny(&distrmodulev1.Module{}),
			},
			{
				Name: capabilitytypes.ModuleName,
				Config: sdkappconfig.WrapAny(&capabilitymodulev1.Module{
					SealKeeper: true,
				}),
			},
			{
				Name:   evidencetypes.ModuleName,
				Config: sdkappconfig.WrapAny(&evidencemodulev1.Module{}),
			},
			{
				Name:   minttypes.ModuleName,
				Config: sdkappconfig.WrapAny(&mintmodulev1.Module{}),
			},
			{
				Name:   feegrant.ModuleName,
				Config: sdkappconfig.WrapAny(&feegrantmodulev1.Module{}),
			},
			{
				Name:   govtypes.ModuleName,
				Config: sdkappconfig.WrapAny(&govmodulev1.Module{}),
			},
			{
				Name:   crisistypes.ModuleName,
				Config: sdkappconfig.WrapAny(&crisismodulev1.Module{}),
			},
			{
				Name:   consensustypes.ModuleName,
				Config: sdkappconfig.WrapAny(&consensusmodulev1.Module{}),
			},

			// Stratos modules
			{
				Name:   registertypes.ModuleName,
				Config: sdkappconfig.WrapAny(&registermodulev1.Module{}),
			},
			{
				Name:   pottypes.ModuleName,
				Config: sdkappconfig.WrapAny(&potmodulev1.Module{}),
			},
			{
				Name:   sdstypes.ModuleName,
				Config: sdkappconfig.WrapAny(&sdsmodulev1.Module{}),
			},
			{
				Name:   evmtypes.ModuleName,
				Config: sdkappconfig.WrapAny(&evmmodulev1.Module{}),
			},
		},
	})
)
