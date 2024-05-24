// SPDX-License-Identifier: MIT

pragma solidity ^0.8.20;

interface ISystemContract {
    function registeredTypeUrls(
        string memory typeUrl
    ) external view returns (bool);

    function allowedExecutors(address guy) external view returns (bool);

    struct ProtoMessage {
        address signer;
        string typeUrl;
    }

    event MsgExecuted(
        address indexed signer,
        string indexed typeUrl,
        bytes data
    );

    enum Action {
        ADD_EXECUTOR,
        REMOVE_EXECUTOR,
        APPROVE_PROTO,
        DISAPPROVE_PROTO
    }

    error ActionFailed(Action name, bytes data);
    error NotAllowedTypeUrl(string typeUrl);
    error NotAllowedSigner(address shguy, address pbguy);
    error NotAllowedExecutor(address guy);

    function pause() external;

    function unpause() external;

    function addExecutor(address guy) external;

    function removeExecutor(address guy) external;

    function approveProto(string memory typeUrl) external;

    function disapproveProto(string memory typeUrl) external;

    function runMsg(bytes memory data) external returns (bytes memory);

    function getProtoMessageFromData(
        bytes memory data
    ) external view returns (ProtoMessage memory);
}
