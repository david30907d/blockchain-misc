// SPDX-License-Identifier: MIT
// Tells the Solidity compiler to compile only from v0.8.13 to v0.9.0
pragma solidity ^0.8.13;

contract Oracle {
  uint public gasPrice;
  function setGasPrice(uint _gasPrice) external {
    gasPrice = _gasPrice;
  }

  function getGasPrice() external view returns (uint) {
    return gasPrice;
  } 
}