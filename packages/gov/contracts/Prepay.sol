// SPDX-License-Identifier: MIT

pragma solidity ^0.8.18;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

import "./versions/Version0.sol";

contract Prepay is Version0, OwnableUpgradeable {
    event PrepayCreated(
        address indexed sender,
        address indexed beneficiery,
        uint256 indexed amount,
        uint256 purchased
    );

    // ----- proxy ------
    uint256[49] private __gap;

    // ===== fallbacks =====

    receive() external payable {}

    // Initialize function for proxy constructor. Must be used atomically
    function initialize() public initializer {
        // proxy inits
        __Context_init_unchained();
        __Ownable_init_unchained();
    }

    function prepay() external payable {
        uint256 amount = msg.value;
        require(amount != 0, "P: ZERO_AMOUNT");

        uint256[1] memory input;
        input[0] = uint256(uint160(msg.sender));

        uint256[1] memory output;

        assembly {
            if iszero(call(not(0), 0xf1, amount, input, 0x20, output, 0x20)) {
                revert(0x0, 0x0)
            }
        }

        emit PrepayCreated(address(this), msg.sender, amount, output[0]);
    }
}
