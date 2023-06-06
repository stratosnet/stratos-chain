// SPDX-License-Identifier: MIT

pragma solidity ^0.8.18;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

import "../versions/Version1.sol";

contract PrepayMockOverride is Version1, OwnableUpgradeable {
    event PrepayCreated(
        address indexed sender,
        address indexed beneficiery,
        uint256 indexed amount,
        uint256 purchased
    );

    address public EXECUTOR;
    uint8 public IS_TEST;

    // ----- proxy ------
    uint256[47] private __gap;

    // ===== fallbacks =====

    receive() external payable {}

    function upgrade() public {
        require(
            _msgSender() == 0x1000000000000000000000000000000000000000,
            "PMO: WRONG_SENDER"
        );
        IS_TEST = 1;
        EXECUTOR = _msgSender();
    }

    // Initialize function for proxy constructor. Must be used atomically
    function initialize() public initializer {
        // proxy inits
        __Context_init_unchained();
        __Ownable_init_unchained();
    }

    function prepay() external {
        emit PrepayCreated(address(this), msg.sender, 0, 0);
    }
}
