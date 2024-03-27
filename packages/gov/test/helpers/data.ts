import {
  MsgBeginRedelegate,
  MsgDelegate, MsgUndelegate,
} from "cosmjs-types/cosmos/staking/v1beta1/tx";
import {
  MsgWithdrawDelegatorReward,
} from "cosmjs-types/cosmos/distribution/v1beta1/tx";
import { BigNumberish } from "ethers";
import { coin } from "@cosmjs/amino";


export const getEVMSdkMsgTestCases = async ({
  senderAddress,
  validatorAddress,
  amount,
}: {
  senderAddress: string
  validatorAddress: string
  amount: BigNumberish
}) => {
  return [
    {
      name: 'Test simple delegate and ok',
      proto: {
        struct: MsgDelegate,
      },
      data: {
        delegatorAddress: senderAddress,
        validatorAddress,
        amount: coin((amount / BigInt(2)).toString(), 'wei'),
      } as MsgDelegate,
      error: null,
      result: null,
    },
    {
      name: 'Test simple redelegate and ok',
      proto: {
        struct: MsgBeginRedelegate,
      },
      data: {
        delegatorAddress: senderAddress,
        validatorSrcAddress: validatorAddress,
        validatorDstAddress: validatorAddress,
        amount: coin((amount / BigInt(2)).toString(), 'wei'),
      } as MsgBeginRedelegate,
      error: null,
      result: null,
    },
    {
      name: 'Test simple undelegate and ok',
      proto: {
        struct: MsgUndelegate,
      },
      data: {
        delegatorAddress: senderAddress,
        validatorAddress,
        amount: coin((amount / BigInt(2)).toString(), 'wei'),
      } as MsgUndelegate,
      error: null,
      result: null,
    },
    {
      name: 'Test simple withdraw reward and ok',
      proto: {
        struct: MsgWithdrawDelegatorReward,
      },
      data: {
        delegatorAddress: senderAddress,
        validatorAddress,
      } as MsgWithdrawDelegatorReward,
      error: null,
      result: null,
    },
  ];
};