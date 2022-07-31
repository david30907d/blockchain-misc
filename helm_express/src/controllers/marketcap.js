// atom: curl --location --request GET 'http://localhost:8080/marketcap?tokenContractAddress=0x83b835c3200147e74bddf3571fde527647b126c7&chainLinkTokenAddress=0xc751E86208F0F8aF2d5CD0e29716cA7AD98B5eF5' --header 'Content-Type: application/json'
// BAT: curl --location --request GET 'http://localhost:8080/marketcap?tokenContractAddress=0xE43A8b1bfD19721E918b6068c158bf8faB41aBee&chainLinkTokenAddress=0x031dB56e01f82f20803059331DC6bEe9b17F7fC9' --header 'Content-Type: application/json'
const Web3 = require("web3");
const RedisClient = require("../models/redis.js").client;

const web3 = new Web3(
  "https://rinkeby.infura.io/v3/18f80354699e46188ac9b12df50f9296"
);
const chainLinkModel = require("../models/chainLinkModel");
const correctABI = [
  {
    constant: true,
    inputs: [],
    name: "name",
    outputs: [{ name: "", type: "string" }],
    payable: false,
    stateMutability: "view",
    type: "function",
  },
  {
    constant: false,
    inputs: [
      { name: "_spender", type: "address" },
      { name: "_value", type: "uint256" },
    ],
    name: "approve",
    outputs: [{ name: "", type: "bool" }],
    payable: false,
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    constant: true,
    inputs: [],
    name: "totalSupply",
    outputs: [{ name: "", type: "uint256" }],
    payable: false,
    stateMutability: "view",
    type: "function",
  },
  {
    constant: false,
    inputs: [
      { name: "_from", type: "address" },
      { name: "_to", type: "address" },
      { name: "_value", type: "uint256" },
    ],
    name: "transferFrom",
    outputs: [{ name: "", type: "bool" }],
    payable: false,
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    constant: true,
    inputs: [],
    name: "decimals",
    outputs: [{ name: "", type: "uint8" }],
    payable: false,
    stateMutability: "view",
    type: "function",
  },
  {
    constant: true,
    inputs: [{ name: "_owner", type: "address" }],
    name: "balanceOf",
    outputs: [{ name: "balance", type: "uint256" }],
    payable: false,
    stateMutability: "view",
    type: "function",
  },
  {
    constant: true,
    inputs: [],
    name: "symbol",
    outputs: [{ name: "", type: "string" }],
    payable: false,
    stateMutability: "view",
    type: "function",
  },
  {
    constant: false,
    inputs: [
      { name: "_to", type: "address" },
      { name: "_value", type: "uint256" },
    ],
    name: "transfer",
    outputs: [{ name: "", type: "bool" }],
    payable: false,
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    constant: true,
    inputs: [
      { name: "_owner", type: "address" },
      { name: "_spender", type: "address" },
    ],
    name: "allowance",
    outputs: [{ name: "", type: "uint256" }],
    payable: false,
    stateMutability: "view",
    type: "function",
  },
  { payable: true, stateMutability: "payable", type: "fallback" },
  {
    anonymous: false,
    inputs: [
      { indexed: true, name: "owner", type: "address" },
      { indexed: true, name: "spender", type: "address" },
      { indexed: false, name: "value", type: "uint256" },
    ],
    name: "Approval",
    type: "event",
  },
  {
    anonymous: false,
    inputs: [
      { indexed: true, name: "from", type: "address" },
      { indexed: true, name: "to", type: "address" },
      { indexed: false, name: "value", type: "uint256" },
    ],
    name: "Transfer",
    type: "event",
  },
];

async function controller(req, res, next) {
  const tokenContractAddress = req.query.tokenContractAddress;
  const chainLinkTokenAddress = req.query.chainLinkTokenAddress;
  try {
    await RedisClient.connect();
    const cacheResult = await RedisClient.get(tokenContractAddress);
    let tokenName, marketCap, tokenPrice, resp;
    if (cacheResult === null) {
      // cold start
      const tokenRouter = await new web3.eth.Contract(
        correctABI,
        tokenContractAddress
      );
      const totalSupply = await tokenRouter.methods.totalSupply().call();
      const decimals = await tokenRouter.methods.decimals().call();
      tokenName = await tokenRouter.methods.name().call();
      const correctTotalSupply = totalSupply / 10 ** decimals;

      tokenPrice = await chainLinkModel.getPriceFromChainLink(
        chainLinkTokenAddress
      );
      marketCap = tokenPrice * correctTotalSupply;
      resp = {
        status: 200,
        token: tokenName,
        marketCap: marketCap,
        price: tokenPrice,
        cacheStatus: true,
      };
      await RedisClient.set(tokenContractAddress, JSON.stringify(resp));
      resp.cacheStatus = false;
    } else {
      // hit cache
      resp = JSON.parse(cacheResult);
    }

    await RedisClient.disconnect();
    return res.status(200).json(resp);
  } catch (error) {
    return res.status(500).json({ status: 500, error: `${error}` });
  }
}

// add the code below
module.exports = { controller };
