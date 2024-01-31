// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

contract Target {
    address public importantAddress;

    function updateAddress(address _address) public {
        importantAddress = _address;
    }
}

contract Proxy {

    address public stateVariable;
    function delegateCallAddressUpdate(address target, address newAddress) public {
        (bool success, ) = target.delegatecall(
            abi.encodeWithSignature("updateAddress(address)", newAddress)
        );
        require(success, "Delegatecall failed");
    }
}
