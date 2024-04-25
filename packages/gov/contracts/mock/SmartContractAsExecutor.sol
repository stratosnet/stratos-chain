// SPDX-License-Identifier: MIT

pragma solidity ^0.8.20;

import "../SystemContract.sol";

contract Failurer {
    function fail() public {
        revert("never");
    }
}

contract SmartContractAsExecutor {
    receive() external payable {}

    Failurer public fr;

    constructor() {
        fr = new Failurer();
    }

    event CallResult(bool success, bytes result);

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

    function runMsgWithSkipRevert(
        SystemContract systemContract,
        bytes memory data
    ) public payable returns (bytes memory) {
        systemContract.runMsg(data);
        (bool success, bytes memory result) = address(fr).call(
            abi.encodeWithSelector(Failurer.fail.selector)
        );
        emit CallResult(success, result);
    }
}
