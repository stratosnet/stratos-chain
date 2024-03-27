import { toBech32 } from "@cosmjs/encoding";
import { config, ethers } from "hardhat"


export const sleep = (ms) => {
  return new Promise(resolve => setTimeout(resolve, ms));
}

export const timeout = (prom: Promise<T>, time: number): Promise<T> =>
  Promise.race([prom, new Promise((_r, rej) => setTimeout(rej, time))]);

export const blockAwaiter = async (count, awaitTime = 1000) => {
  let block = await ethers.provider.getBlock();
  let startBNumber = block.number;
  const endBNumber = block.number + count;
  if (startBNumber >= endBNumber) {
    console.log("No blocks to await, skipping");
    return block;
  }
  console.log(`Awaiting for block #${endBNumber} (current: #${startBNumber})`)
  while (startBNumber < endBNumber) {
    await sleep(awaitTime);
    block = await ethers.provider.getBlock();
    if (block.number !== startBNumber) {
      startBNumber = block.number;
      console.log(`Awaiting for block #${endBNumber} (current: #${startBNumber})`)
    }
  }
  console.log(`Block #${endBNumber} successfully awaited`)
  return block;
}

export const getWallet = (network, index = 0) => {
  const hhAccounts = config.networks[network].accounts;
  const mnemonic = ethers.Mnemonic.fromPhrase(hhAccounts.mnemonic)
  const wallet = ethers.HDNodeWallet.fromMnemonic(mnemonic, hhAccounts.path + `/${index}`);
  return wallet.connect(ethers.provider);
}

export const ethToBech32Address = (addr) => {
  return toBech32('st', ethers.getBytes(addr))
}
