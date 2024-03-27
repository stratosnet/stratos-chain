// SPDX-License-Identifier: MIT

pragma solidity ^0.8.20;

import "../SystemContract.sol";

contract SmartContractAsExecutor {
    receive() external payable {}

    function getProtoMessageFromData(
        SystemContract systemContract,
        bytes memory data
    ) public view returns (SystemContract.ProtoMessage memory) {
        return systemContract.getProtoMessageFromData(data);
    }

    function runMsg(
        SystemContract systemContract,
        bytes memory data
    ) public payable returns (bytes memory) {
        return systemContract.runMsg(data);
    }

    function runMsgWithRevert(
        SystemContract systemContract,
        bytes memory data
    ) public payable returns (bytes memory) {
        systemContract.runMsg(data);
        revert("never");
    }
}
