// const ConvertLib = artifacts.require("ConvertLib");
const Oracle = artifacts.require("Oracle");

module.exports = function(deployer) {
  // deployer.deploy(ConvertLib);
  // deployer.link(ConvertLib, Oracle);
  deployer.deploy(Oracle);
};
