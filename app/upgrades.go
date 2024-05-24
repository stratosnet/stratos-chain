package app

import (
	v012 "github.com/stratosnet/stratos-chain/app/upgrades/v012"
)

// registerUpgradeHandlers registers all the upgrade handlers that are supported by the app
func (app *StratosApp) registerUpgradeHandlers() {
	app.registerUpgrade(v012.NewUpgrade(app.ModuleManager, app.Configurator(), app.paramsKeeper, app.consensusKeeper))
}
