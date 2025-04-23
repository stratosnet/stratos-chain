package app

import (
	v012 "github.com/stratosnet/stratos-chain/app/upgrades/v012"
	v013 "github.com/stratosnet/stratos-chain/app/upgrades/v013"
)

// registerUpgradeHandlers registers all the upgrade handlers that are supported by the app
func (app *StratosApp) registerUpgradeHandlers() {
	app.registerUpgrade(v012.NewUpgrade(app.ModuleManager, app.Configurator(), app.paramsKeeper, app.consensusKeeper))
	app.registerUpgrade(v013.NewUpgrade(app.ModuleManager, app.Configurator()))
}
