// SPDX-License-Identifier: MIT

pragma solidity ^0.8.18;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

import "./versions/Version0.sol";

contract Prepay is Version0, OwnableUpgradeable {
    event PrepayCreated(
        address indexed sender,
        address beneficiery,
        uint256 amount,
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

    function prepay(address beneficiary, uint256 amount) external payable {
        require(amount != 0, "P: ZERO_AMOUNT");
        require(beneficiary != address(0), "P: ZERO_ADDRESS");

        address sender = _msgSender();

        uint256[3] memory input;
        input[0] = uint256(uint160(address(sender)));
        input[1] = uint256(uint160(address(beneficiary)));
        input[2] = amount;

        uint256[1] memory output;

        assembly {
            // TODO: PROXY: Implement!!!
            if iszero(call(not(0), 0xf1, 0x0, input, 0x20, output, 0x20)) {
                revert(0x0, 0x0)
            }
        }

        emit PrepayCreated(sender, beneficiary, amount, output[0]);
    }
}
