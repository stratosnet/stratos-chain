// SPDX-License-Identifier: MIT

pragma solidity ^0.8.20;

interface ISystemContract {
    function registeredTypeUrls(
        string memory typeUrl
    ) external view returns (bool);

    struct ProtoMessage {
        address signer;
        string typeUrl;
    }

    event MsgExecuted(
        address indexed signer,
        string indexed typeUrl,
        bytes data
    );

    error ApproveActionTypeUrl(string typeUrl);
    error NotAllowedTypeUrl(string typeUrl);
    error NotAllowedSigner(address shguy, address pbguy);

    function pause() external;

    function unpause() external;

    function approveProto(string memory typeUrl) external;

    function disapproveProto(string memory typeUrl) external;

    function runMsg(bytes memory data) external returns (bytes memory);

    function getProtoMessageFromData(
        bytes memory data
    ) external view returns (ProtoMessage memory);
}
