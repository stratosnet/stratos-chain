// SPDX-License-Identifier: MIT

pragma solidity ^0.8.20;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";

import "./versions/Version0.sol";
import "./interfaces/ISystemContract.sol";

contract SystemContract is
    Version0,
    OwnableUpgradeable,
    AccessControlUpgradeable,
    PausableUpgradeable,
    ReentrancyGuardUpgradeable,
    ISystemContract
{
    bytes32 public constant MAINTAINER_ROLE = keccak256("MAINTAINER_ROLE");
    mapping(string => bool) public override registeredTypeUrls;

    // ===== fallbacks =====

    receive() external payable {}

    // Initialize function for proxy constructor. Must be used atomically
    function initialize() public initializer {
        // proxy inits
        __Context_init();
        __Ownable_init(_msgSender());
        __ReentrancyGuard_init();
        // roles update
        _grantRole(DEFAULT_ADMIN_ROLE, _msgSender());
        _grantRole(MAINTAINER_ROLE, _msgSender());
        // EOA
        _grantRole(MAINTAINER_ROLE, 0x31AB8edC475A73daA2a1E2B2933f9Dd19FA8e6F9);
        // proto perm add
        _initAllowedProtos();
    }

    function _initAllowedProtos() internal {
        registeredTypeUrls["/stratos.sds.v1.MsgPrepay"] = true;
        registeredTypeUrls["/cosmos.staking.v1beta1.MsgDelegate"] = true;
        registeredTypeUrls["/cosmos.staking.v1beta1.MsgUndelegate"] = true;
        registeredTypeUrls["/cosmos.staking.v1beta1.MsgBeginRedelegate"] = true;
        registeredTypeUrls[
            "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward"
        ] = true;
    }

    function pause() external override onlyRole(MAINTAINER_ROLE) {
        _pause();
    }

    function unpause() external override onlyRole(MAINTAINER_ROLE) {
        _unpause();
    }

    function approveProto(
        string memory typeUrl
    ) external override onlyRole(MAINTAINER_ROLE) {
        if (registeredTypeUrls[typeUrl]) {
            revert ApproveActionTypeUrl(typeUrl);
        }
        registeredTypeUrls[typeUrl] = true;
    }

    function disapproveProto(
        string memory typeUrl
    ) external override onlyRole(MAINTAINER_ROLE) {
        if (!registeredTypeUrls[typeUrl]) {
            revert ApproveActionTypeUrl(typeUrl);
        }
        delete registeredTypeUrls[typeUrl];
    }

    function runMsg(
        bytes memory data
    ) public override whenNotPaused nonReentrant returns (bytes memory) {
        ProtoMessage memory pb = getProtoMessageFromData(data);
        if (!registeredTypeUrls[pb.typeUrl]) {
            revert NotAllowedTypeUrl(pb.typeUrl);
        }
        if (msg.sender != pb.signer) {
            revert NotAllowedSigner(msg.sender, pb.signer);
        }
        bytes memory resp = opRunSdkMsg(data);

        emit MsgExecuted(msg.sender, pb.typeUrl, resp);

        return resp;
    }

    function opRunSdkMsg(
        bytes memory data
    ) internal returns (bytes memory _returnData) {
        bytes memory inData = abi.encode(msg.sender, data);

        assembly {
            let result := call(
                gas(),
                0xdd, // RUNSDKMSG
                0,
                add(inData, 0x20), // trim first encoding prefix from abi.encode
                mload(inData), // no need extra 32 bytes as trimmed
                0,
                0
            )
            if iszero(result) {
                revert(0x0, 0x0)
            }
            _returnData := mload(0x40)
            let _returnSize := returndatasize()
            mstore(_returnData, _returnSize)
            returndatacopy(add(_returnData, 0x20), 0, _returnSize)
            mstore(0x40, add(_returnData, add(_returnSize, 0x20)))
        }
    }

    function getProtoMessageFromData(
        bytes memory data
    ) public view override returns (ProtoMessage memory pb) {
        (pb.signer, pb.typeUrl) = abi.decode(opPpfd(data), (address, string));
        return pb;
    }

    // NOTE: Should be moved better in sol protos, currently opcode
    function opPpfd(
        bytes memory data
    ) internal view returns (bytes memory _returnData) {
        assembly {
            let result := staticcall(
                gas(),
                0xdc, // PPFD
                data,
                add(mload(data), 0x20),
                0,
                0
            )
            if iszero(result) {
                revert(0x0, 0x0)
            }
            _returnData := mload(0x40)
            let _returnSize := returndatasize()
            mstore(_returnData, _returnSize)
            returndatacopy(add(_returnData, 0x20), 0, _returnSize)
            mstore(0x40, add(_returnData, add(_returnSize, 0x20)))
        }
    }

    // ----- proxy ------
    uint256[49] private __gap;
}
