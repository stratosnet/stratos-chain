import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployments, getNamedAccounts } = hre;
  const { deploy } = deployments;
  const { deployer } = await getNamedAccounts();

  await deploy("SmartContractAsExecutor", {
    from: deployer,
    skipIfAlreadyDeployed: false,
    args: [],
    log: true,
  });
};
func.tags = ["0002_deploy_sc_as_exec"];
func.dependencies = ["0001_deploy_system_contract"];

export default func;
