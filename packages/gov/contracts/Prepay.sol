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

    event TrackGas(uint256 indexed gas, uint256 indexed delta);

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

    function prepay(address beneficiary) external payable {
        uint256 amount = msg.value;
        require(amount != 0, "P: ZERO_AMOUNT");
        require(beneficiary != address(0), "P: ZERO_ADDRESS");

        address sender = _msgSender();

        uint256[2] memory input;
        input[0] = uint256(uint160(address(sender)));
        input[1] = uint256(uint160(address(beneficiary)));

        uint256[1] memory output;

        uint256 beforeGas = gasleft();

        assembly {
            if iszero(call(not(0), 0xf1, amount, input, 0x40, output, 0x20)) {
                revert(0x0, 0x0)
            }
        }

        uint256 afterGas = gasleft();

        emit TrackGas(beforeGas, 0);
        emit TrackGas(afterGas, beforeGas - afterGas);

        emit PrepayCreated(sender, beneficiary, amount, output[0]);
    }
}
