import { ethers } from "hardhat"
import { Prepay } from "../typechain-types";

describe("Prepay local test", () => {
  it("should prepay be executed", async () => {
    const signer = await ethers.getNamedSigner("tester");

    // could be changed
    const totalUnissuedPrepayAddress = '0xa2776119cb0503aaa30f11b6d3fb1b6ba539c61d';

    const balanceBefore = await ethers.provider.getBalance(totalUnissuedPrepayAddress);
    console.log("balanceBefore", balanceBefore.toString());

    const contract: Prepay = await ethers.getContractAt("Prepay", "0x1000000000000000000000000000000000010101", signer);
    const value = ethers.utils.parseUnits('1', 'wei');
    const tx = await contract.prepay({ value }).then((tx) => tx.wait());

    console.log("tx hash", tx.transactionHash);
  })
})
