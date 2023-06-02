import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployments, getNamedAccounts } = hre;
  const { deploy } = deployments;
  const { tester } = await getNamedAccounts();

  await deploy("Prepay", {
    from: tester,
    skipIfAlreadyDeployed: true,
    args: [],
    log: true,
  });
};
func.tags = ["0001_deploy_prepay"];

export default func;
