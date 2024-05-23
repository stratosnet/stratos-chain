import { ethers, deployments, network, getChainId } from "hardhat"
import { SystemContract, SmartContractAsExecutor } from "../typechain-types";
import { expect } from "chai";

import {
  MsgSend,
} from "cosmjs-types/cosmos/bank/v1beta1/tx";
import {
  MsgDelegate,
} from "cosmjs-types/cosmos/staking/v1beta1/tx";
import { Any } from 'cosmjs-types/google/protobuf/any';
import { coins, coin } from "@cosmjs/amino";
import {
  blockAwaiter,
  ethToBech32Address,
  getWallet,
  timeout,
} from "./helpers/ethUtils";
import { getEVMSdkMsgTestCases } from "./helpers/data";
import { ContractTransactionResponse } from "ethers";


const addExecutorIfNotAdded = async (contract: SystemContract, addr: string): Promise<void> => {
  const isAllowed = await contract.allowedExecutors.staticCall(addr);
  if (!isAllowed) {
    console.log(`Adding executor "${addr}" as allowed...`)
    await contract.addExecutor(addr);
    console.log(`Approved executor "${addr}".`)
  }
}

const removeExecutorIfAdded = async (contract: SystemContract, addr: string): Promise<void> => {
  const isAllowed = await contract.allowedExecutors.staticCall(addr);
  if (isAllowed) {
    console.log(`Removing executor "${addr}" as allowed...`)
    await contract.removeExecutor(addr);
    console.log(`Removed executor "${addr}".`)
  }
}

const allowProtoIfNotEnabled = async (contract: SystemContract, typeUrl: string): Promise<void> => {
  const protoAllowed = await contract.registeredTypeUrls.staticCall(typeUrl);
  if (!protoAllowed) {
    console.log(`Adding type url "${typeUrl}" as allowed...`)
    await contract.approveProto(typeUrl);
    console.log(`Approved type url "${typeUrl}".`)
  }
}

const disallowProtoIfNotEnabled = async (contract: SystemContract, typeUrl: string): Promise<void> => {
  const protoAllowed = await contract.registeredTypeUrls.staticCall(typeUrl);
  if (protoAllowed) {
    console.log(`Remove type url "${typeUrl}"...`)
    await contract.disapproveProto(typeUrl);
    console.log(`Removed type url "${typeUrl}".`)
  }
}


const setupTest = deployments.createFixture(
  async ({ deployments, getNamedAccounts, ethers }) => {
    const signer = await ethers.getNamedSigner("deployer");
    const signerWrong = await ethers.getNamedSigner("proxyAdmin");

    const validatorAddress = 'stvaloper1yamhfcvt09n8ctau9vzv4ad6dnmlpvfax86jk0';

    expect(
      await ethers.provider.getBalance(signer.address)
    ).to.be.not.eq(0);

    expect(
      await ethers.provider.getBalance(signerWrong.address)
    ).to.be.not.eq(0);

    // could be changed
    // let systemContractAddress = '0x1000000000000000000000000000000000010101';
    let systemContractAddress: string = undefined;
    let executorAddress: string = undefined;

    if (!systemContractAddress) {
      systemContractAddress = await (await ethers.getContract("SystemContract", signer)).getAddress();
    }

    if (!executorAddress) {
      executorAddress = await (await ethers.getContract("SmartContractAsExecutor", signer)).getAddress();
    }
    const contract = await ethers.getContractAt("SystemContract", systemContractAddress, signer);
    const scAsExecutorContract = await ethers.getContractAt("SmartContractAsExecutor", executorAddress, signer);

    await allowProtoIfNotEnabled(contract, MsgSend.typeUrl);

    await addExecutorIfNotAdded(contract, executorAddress);

    return {
      systemContractAddress,
      signer,
      signerWrong,
      contract,
      scAsExecutorContract,
      validatorAddress,
    }
  }
);

describe("SystemContract local test", () => {

  it("should execute data tests", async () => {
    const { signer, contract, validatorAddress } = await setupTest();
    const amount = ethers.parseEther('0.001');
    const testCases = await getEVMSdkMsgTestCases({
      senderAddress: ethToBech32Address(signer.address),
      validatorAddress,
      amount,
    });


    for (let i = 0; i <= testCases.length - 1; i++) {
      const testCase = testCases[i];

      console.group(`--- ${testCase.name} ---`);

      const { struct: MsgProtobuf } = testCase.proto;

      await allowProtoIfNotEnabled(contract, MsgProtobuf.typeUrl);

      const payload = MsgProtobuf.fromJSON(testCase.data);
      const anyMsg: Any = {
        typeUrl: MsgProtobuf.typeUrl,
        value: MsgProtobuf.encode(payload).finish(),
      };

      const data = ethers.getBytes(Any.encode(anyMsg).finish());

      const [block, eth_maxPriorityFeePerGas] = await Promise.all([
        await ethers.provider.getBlock("latest"),
        await ethers.provider.send("eth_maxPriorityFeePerGas", []),
      ]);
      if (!block || !eth_maxPriorityFeePerGas) {
        console.log('Block or prior fee not found, skipping...')
        return;
      }

      if (testCase.error !== null) {
        continue;
      }

      const nonceBefore = await ethers.provider.getTransactionCount(signer.address);

      const coro = contract.runMsg(data, {
        gasLimit: 300_000,
        maxPriorityFeePerGas: BigInt(eth_maxPriorityFeePerGas),
        maxFeePerGas: block.baseFeePerGas,
      });
      await expect(coro).to.be.not.revertedWithCustomError(contract, 'NotAllowedTypeUrl');
      await expect(coro).to.be.not.revertedWithCustomError(contract, 'NotAllowedSigner');
      await expect(coro).to.be.not.revertedWithoutReason();
      await expect(coro).to.be.emit(contract, 'MsgExecuted');

      const resp = await timeout(coro, 5000);

      expect(resp).to.be.not.eq(null);

      const tx = await ethers.provider.getTransactionReceipt(resp.hash);
      const nonceAfter = await ethers.provider.getTransactionCount(signer.address);
      const balanceFromAfter = await ethers.provider.getBalance(signer.address);

      expect(nonceAfter).to.be.eq(nonceBefore + 1);

      console.log('--------------------');
      console.groupEnd();
    }
  })

  it("should check proto and ok", async () => {
    const { signer, contract } = await setupTest();
    const amount = ethers.parseEther('0.001');

    const payload: MsgSend = {
      fromAddress: ethToBech32Address(signer.address),
      toAddress: ethToBech32Address(signer.address),
      amount: coins((amount / BigInt(2)).toString(), 'wei'),
    };
    const anyMsg: Any = {
      typeUrl: MsgSend.typeUrl,
      value: MsgSend.encode(payload).finish(),
    };

    const data = ethers.getBytes(Any.encode(anyMsg).finish());

    const pb = await contract.getProtoMessageFromData.staticCall(data);
    expect(pb.typeUrl).to.be.eq(MsgSend.typeUrl);
    expect(pb.signer).to.be.eq(signer.address);
  })

  it("should execute proto and ok", async () => {
    const { signer, contract } = await setupTest();
    const testAddr = '0x000000000000000000000000000000000000dEaD';
    const amount = ethers.parseEther('0.001');

    const payload: MsgSend = {
      fromAddress: ethToBech32Address(signer.address),
      toAddress: ethToBech32Address(testAddr),
      amount: coins(amount.toString(), 'wei'),
    };
    const anyMsg: Any = {
      typeUrl: MsgSend.typeUrl,
      value: MsgSend.encode(payload).finish(),
    };

    const data = ethers.getBytes(Any.encode(anyMsg).finish());

    const [block, eth_maxPriorityFeePerGas] = await Promise.all([
      await ethers.provider.getBlock("latest"),
      await ethers.provider.send("eth_maxPriorityFeePerGas", []),
    ])
    if (!block || !eth_maxPriorityFeePerGas) {
      console.log('Block or prior fee not found, skipping...')
      return;
    }

    const nonceBefore = await ethers.provider.getTransactionCount(signer.address);
    const balanceFromBefore = await ethers.provider.getBalance(signer.address);
    const balanceToBefore = await ethers.provider.getBalance(testAddr);

    const coro = contract.runMsg(data, {
      gasLimit: 150_000,
      maxPriorityFeePerGas: BigInt(eth_maxPriorityFeePerGas),
      maxFeePerGas: block.baseFeePerGas,
    });
    await expect(coro).to.be.not.revertedWithCustomError(contract, 'NotAllowedTypeUrl');
    await expect(coro).to.be.not.revertedWithCustomError(contract, 'NotAllowedSigner');
    await expect(coro).to.be.not.revertedWithoutReason();
    await expect(coro).to.be.emit(contract, 'MsgExecuted');
    const resp = await coro;

    const nonceAfter = await ethers.provider.getTransactionCount(signer.address);
    const balanceFromAfter = await ethers.provider.getBalance(signer.address);
    const balanceToAfter = await ethers.provider.getBalance(testAddr);
    const tx = await ethers.provider.getTransactionReceipt(resp.hash);
    console.log('tx', tx)
    console.log('logs', tx?.logs)

    const txFee = tx?.gasUsed * block.baseFeePerGas;
    expect(nonceAfter).to.be.eq(nonceBefore + 1);
    expect(balanceFromAfter).to.be.eq(balanceFromBefore - amount - txFee);
    expect(balanceToAfter).to.be.eq(balanceToBefore + amount);
  })

  it("should not execute proto with not allowed executor and fail", async () => {
    const { signer, contract, scAsExecutorContract } = await setupTest();
    const testAddr = '0x000000000000000000000000000000000000dEaD';
    const scAsExecutorContractAddress = await scAsExecutorContract.getAddress();
    const systemContractAddress = await contract.getAddress();
    const amount = ethers.parseEther('0.001');
    
    const payload: MsgSend = {
      fromAddress: ethToBech32Address(scAsExecutorContractAddress),
      toAddress: ethToBech32Address(testAddr),
      amount: coins(amount.toString(), 'wei'),
    };
    const anyMsg: Any = {
      typeUrl: MsgSend.typeUrl,
      value: MsgSend.encode(payload).finish(),
    };
    
    const data = ethers.getBytes(Any.encode(anyMsg).finish());
    
    const [block, eth_maxPriorityFeePerGas] = await Promise.all([
      await ethers.provider.getBlock("latest"),
      await ethers.provider.send("eth_maxPriorityFeePerGas", []),
    ])
    if (!block || !eth_maxPriorityFeePerGas) {
      console.log('Block or prior fee not found, skipping...')
      return;
    }
    await removeExecutorIfAdded(contract, await scAsExecutorContract.getAddress());

    expect(
      await scAsExecutorContract.runMsg(systemContractAddress, data, {
        gasLimit: 150_000,
        maxPriorityFeePerGas: BigInt(eth_maxPriorityFeePerGas),
        maxFeePerGas: block.baseFeePerGas,
        value: amount,
      })
    ).to.be.revertedWithCustomError(contract, 'NotAllowedExecutor');
  })

  it("should not execute proto with not registered type and fail", async () => {
    const { signer, contract, validatorAddress } = await setupTest();
    const amount = ethers.parseEther('0.001');

    await disallowProtoIfNotEnabled(contract, MsgDelegate.typeUrl);

    const payload: MsgDelegate = {
      delegatorAddress: ethToBech32Address(signer.address),
      validatorAddress,
      amount: coin(amount.toString(), 'wei'),
    };
    const anyMsg: Any = {
      typeUrl: MsgDelegate.typeUrl,
      value: MsgDelegate.encode(payload).finish(),
    };

    const data = ethers.getBytes(Any.encode(anyMsg).finish());

    const coro = contract.runMsg.staticCall(data, {
      gasLimit: 150_000,
    });
    await expect(coro).to.be.revertedWithCustomError(contract, 'NotAllowedTypeUrl');
  })

  it("should not execute proto with wrong signer and fail", async () => {
    const { signer, signerWrong, contract } = await setupTest();
    const testAddr = '0x000000000000000000000000000000000000dEaD';
    const amount = ethers.parseEther('0.001');

    const payload: MsgSend = {
      fromAddress: ethToBech32Address(signer.address),
      toAddress: ethToBech32Address(testAddr),
      amount: coins(amount.toString(), 'wei'),
    };
    const anyMsg: Any = {
      typeUrl: MsgSend.typeUrl,
      value: MsgSend.encode(payload).finish(),
    };

    const data = ethers.getBytes(Any.encode(anyMsg).finish());

    const coro = contract.connect(signerWrong).runMsg.staticCall(data, {
      gasLimit: 150_000,
    });
    await expect(coro).to.be.revertedWithCustomError(contract, 'NotAllowedSigner');
  })

  it("should execute proto with sc signer and ok", async () => {
    const { signer, contract, scAsExecutorContract } = await setupTest();
    const testAddr = '0x000000000000000000000000000000000000dEaD';
    const scAsExecutorContractAddress = await scAsExecutorContract.getAddress();
    const systemContractAddress = await contract.getAddress();
    const amount = ethers.parseEther('0.001');

    const payload: MsgSend = {
      fromAddress: ethToBech32Address(scAsExecutorContractAddress),
      toAddress: ethToBech32Address(testAddr),
      amount: coins(amount.toString(), 'wei'),
    };
    const anyMsg: Any = {
      typeUrl: MsgSend.typeUrl,
      value: MsgSend.encode(payload).finish(),
    };

    const data = ethers.getBytes(Any.encode(anyMsg).finish());

    const [block, eth_maxPriorityFeePerGas] = await Promise.all([
      await ethers.provider.getBlock("latest"),
      await ethers.provider.send("eth_maxPriorityFeePerGas", []),
    ])
    if (!block || !eth_maxPriorityFeePerGas) {
      console.log('Block or prior fee not found, skipping...')
      return;
    }

    const signerNonceBefore = await ethers.provider.getTransactionCount(signer.address);
    const contractNonceBefore = await ethers.provider.getTransactionCount(scAsExecutorContractAddress);
    const balanceSignerBefore = await ethers.provider.getBalance(signer.address);
    const balanceFromBefore = await ethers.provider.getBalance(scAsExecutorContractAddress);
    const balanceToBefore = await ethers.provider.getBalance(testAddr);

    const coro = scAsExecutorContract.runMsg(systemContractAddress, data, {
      gasLimit: 150_000,
      maxPriorityFeePerGas: BigInt(eth_maxPriorityFeePerGas),
      maxFeePerGas: block.baseFeePerGas,
      value: amount,
    });
    await expect(coro).to.be.not.revertedWithCustomError(contract, 'NotAllowedTypeUrl');
    await expect(coro).to.be.not.revertedWithCustomError(contract, 'NotAllowedSigner');
    await expect(coro).to.be.not.revertedWithoutReason();
    await expect(coro).to.be.emit(contract, 'MsgExecuted');
    const resp = await coro;

    const signerNonceAfter = await ethers.provider.getTransactionCount(signer.address);
    const contractNonceAfter = await ethers.provider.getTransactionCount(scAsExecutorContractAddress);
    const balanceSignerAfter = await ethers.provider.getBalance(signer.address);
    const balanceFromAfter = await ethers.provider.getBalance(scAsExecutorContractAddress);
    const balanceToAfter = await ethers.provider.getBalance(testAddr);
    const tx = await ethers.provider.getTransactionReceipt(resp.hash);

    const txFee = tx?.gasUsed * block.baseFeePerGas;
    expect(signerNonceAfter).to.be.eq(signerNonceBefore + 1, "nonce signer mismatch");
    expect(contractNonceAfter).to.be.eq(contractNonceBefore, "nonce contract mismatch");
    expect(balanceSignerAfter).to.be.eq(balanceSignerBefore - amount - txFee, "wrong signer balance");
    expect(balanceFromAfter).to.be.eq(balanceFromBefore, "wrong from balance");
    expect(balanceToAfter).to.be.eq(balanceToBefore + amount, "wrong to balance");
  })

  it("should execute proto with sc signer and ok", async () => {
    const { signer, contract, scAsExecutorContract } = await setupTest();
    const testAddr = '0x000000000000000000000000000000000000dEaD';
    const scAsExecutorContractAddress = await scAsExecutorContract.getAddress();
    const systemContractAddress = await contract.getAddress();
    const amount = ethers.parseEther('0.001');
    // NOTE: Hardcoded, this is intristic fee which will be charged anyway if tx failed
    const txFee = ethers.parseEther('0.00025');

    const payload: MsgSend = {
      fromAddress: ethToBech32Address(scAsExecutorContractAddress),
      toAddress: ethToBech32Address(testAddr),
      amount: coins(amount.toString(), 'wei'),
    };
    const anyMsg: Any = {
      typeUrl: MsgSend.typeUrl,
      value: MsgSend.encode(payload).finish(),
    };

    const data = ethers.getBytes(Any.encode(anyMsg).finish());

    const [block, eth_maxPriorityFeePerGas] = await Promise.all([
      await ethers.provider.getBlock("latest"),
      await ethers.provider.send("eth_maxPriorityFeePerGas", []),
    ])
    if (!block || !eth_maxPriorityFeePerGas) {
      console.log('Block or prior fee not found, skipping...')
      return;
    }

    const signerNonceBefore = await ethers.provider.getTransactionCount(signer.address);
    const contractNonceBefore = await ethers.provider.getTransactionCount(scAsExecutorContractAddress);
    const balanceSignerBefore = await ethers.provider.getBalance(signer.address);
    const balanceFromBefore = await ethers.provider.getBalance(scAsExecutorContractAddress);
    const balanceToBefore = await ethers.provider.getBalance(testAddr);

    const contractTxData = await scAsExecutorContract.runMsgWithRevert.populateTransaction(systemContractAddress, data, {
      gasLimit: 250_000,
      maxPriorityFeePerGas: BigInt(eth_maxPriorityFeePerGas),
      maxFeePerGas: block.baseFeePerGas,
      value: amount,
    });
    contractTxData.chainId = BigInt(await getChainId());
    const wallet = getWallet(network.name, 0);
    const resp = await wallet.sendTransaction(contractTxData);
    await blockAwaiter(3);
    const tx = await ethers.provider.getTransactionReceipt(resp.hash);
    expect(tx).to.be.eq(null);

    const signerNonceAfter = await ethers.provider.getTransactionCount(signer.address);
    const contractNonceAfter = await ethers.provider.getTransactionCount(scAsExecutorContractAddress);
    const balanceSignerAfter = await ethers.provider.getBalance(signer.address);
    const balanceFromAfter = await ethers.provider.getBalance(scAsExecutorContractAddress);
    const balanceToAfter = await ethers.provider.getBalance(testAddr);

    expect(signerNonceAfter).to.be.eq(signerNonceBefore + 1, "nonce signer mismatch");
    expect(contractNonceAfter).to.be.eq(contractNonceBefore, "nonce contract mismatch");
    expect(balanceSignerAfter).to.be.eq(balanceSignerBefore - txFee, "wrong signer balance");
    expect(balanceFromAfter).to.be.eq(balanceFromBefore, "wrong from balance");
    expect(balanceToAfter).to.be.eq(balanceToBefore, "wrong to balance");
  })
})