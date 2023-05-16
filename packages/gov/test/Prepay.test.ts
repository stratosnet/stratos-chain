import { ethers } from "hardhat"
import { Prepay } from "../typechain-types";
import { expect } from "chai";

describe("Prepay local test", () => {

  // could be changed
  const prepayContractAddress = '0x1000000000000000000000000000000000010101';
  const totalUnissuedPrepayAddress = '0xa2776119cb0503aaa30f11b6d3fb1b6ba539c61d';

  const sendValue = ethers.utils.parseUnits('1', 'wei');
  const prepayValue = ethers.utils.parseUnits('1', 'gwei');
  const purchased = ethers.utils.parseUnits('999', 'wei')

  it("should prepay be executed", async () => {
    const signer = await ethers.getNamedSigner("tester");

    console.group('1. Send to unissued address if it is empty')
    const balanceBefore = await ethers.provider.getBalance(totalUnissuedPrepayAddress);
    if (balanceBefore.lt(sendValue)) {
      console.log(`Sending "${sendValue.toString()}" ETH to the address "${totalUnissuedPrepayAddress}"...`)
      const sendTx = await signer.sendTransaction({ value: sendValue, to: totalUnissuedPrepayAddress }).then((tx) => tx.wait());
      console.log(`Completed, result: "${sendTx.transactionHash}"`);
      balanceBefore = await ethers.provider.getBalance(totalUnissuedPrepayAddress);
    } else {
      console.log('No need of send, skipping');
    }
    console.log("Checking balance before next step for unissued address:", balanceBefore.toString());
    console.groupEnd();

    console.group('2. Preparations for prepay')
    const contract: Prepay = await ethers.getContractAt("Prepay", prepayContractAddress, signer);
    console.log(`Executing prepay with "${prepayValue.toString()}" ETH on proxy contract "${prepayContractAddress}"...`)
    const tx = contract.prepay({ value: prepayValue });

    await expect(tx).to
      .emit(contract, "PrepayCreated")
      .withArgs(contract.address, signer.address, prepayValue, purchased);

    console.log(`Completed, result: "${(await tx).hash}"`);

    const balanceAfter = await ethers.provider.getBalance(totalUnissuedPrepayAddress);
    console.log("Checking balance after next step for unissued address:", balanceAfter.toString());
    expect(balanceAfter).to.be.eq(balanceBefore.add(purchased));
    console.groupEnd();
  })
})
