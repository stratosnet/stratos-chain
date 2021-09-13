package keeper_test

import (
	"github.com/stratosnet/stratos-chain/x/evm/types"
)

func (suite *KeeperTestSuite) TestParams() {
	params := suite.app.GetEvmKeeper().GetParams(suite.ctx)
	suite.Require().Equal(types.DefaultParams(), params)
	params.EvmDenom = "ara"
	suite.app.GetEvmKeeper().SetParams(suite.ctx, params)
	newParams := suite.app.GetEvmKeeper().GetParams(suite.ctx)
	suite.Require().Equal(newParams, params)
}
