// chainlink proxy list: https://docs.chain.link/docs/ethereum-addresses/
// rinkby ethscan: https://rinkeby.etherscan.io/token/0xe43a8b1bfd19721e918b6068c158bf8fab41abee
// atom: curl --location --request GET 'http://localhost:8080/marketcap?tokenContractAddress=0x83b835c3200147e74bddf3571fde527647b126c7&chainLinkProxyAddress=0x3539F2E214d8BC7E611056383323aC6D1b01943c' --header 'Content-Type: application/json'
// BAT: curl --location --request GET 'http://localhost:8080/marketcap?tokenContractAddress=0xE43A8b1bfD19721E918b6068c158bf8faB41aBee&chainLinkProxyAddress=0x031dB56e01f82f20803059331DC6bEe9b17F7fC9' --header 'Content-Type: application/json'
// BNB: curl --location --request GET 'http://localhost:8080/marketcap?tokenContractAddress=0xef75685f001210855327e08af9143d3e1de5c758&chainLinkProxyAddress=0xcf0f51ca2cDAecb464eeE4227f5295F2384F84ED' --header 'Content-Type: application/json'
const Web3 = require("web3");
const RedisClient = require("../models/redis.js").client;
const cache_expire_seconds = require("../models/redis.js").cache_expire_seconds;

async function controller(req, res, next) {
  /*
  main logic is as follow:
  1. check if cache exists, if yes then return. Otherwise starts calculating from scratch
  2. use token's contract to get total supply
  3. use chainLink this oracle to get its latest price
  4. use totalSupply * latest price to calculate marketCap
  5. store the result in cache! with pre-defined cache-expire seconds
  */
  const tokenContractAddress = req.query.tokenContractAddress;
  const chainLinkProxyAddress = req.query.chainLinkProxyAddress;
  try {
    await RedisClient.connect();
    const cacheResult = await RedisClient.get(tokenContractAddress);
    let resp;
    if (cacheResult === null) {
      // cold start
      resp = await calculateMarketCap(
        tokenContractAddress,
        chainLinkProxyAddress
      );
      resp.cacheStatus = false;
    } else {
      // hit cache
      resp = JSON.parse(cacheResult);
    }
    res.status(200).json(resp);
  } catch (error) {
    res.status(500).json({ status: 500, error: `${error}` });
  } finally {
    await RedisClient.disconnect();
  }
}

async function calculateMarketCap(tokenContractAddress, chainLinkProxyAddress) {
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
  let tokenName, marketCap, tokenPrice, resp;

  const tokenRouter = await new web3.eth.Contract(
    correctABI,
    tokenContractAddress
  );
  const totalSupply = await tokenRouter.methods.totalSupply().call();
  const decimals = await tokenRouter.methods.decimals().call();
  tokenName = await tokenRouter.methods.name().call();
  const correctTotalSupply = totalSupply / 10 ** decimals;

  tokenPrice = await chainLinkModel.getPriceFromChainLink(
    chainLinkProxyAddress
  );
  marketCap = tokenPrice * correctTotalSupply;
  resp = {
    status: 200,
    token: tokenName,
    marketCap: marketCap,
    price: tokenPrice,
    cacheStatus: true,
  };
  await RedisClient.set(
    tokenContractAddress,
    JSON.stringify(resp),
    "EX",
    cache_expire_seconds
  );
  return resp;
}

// add the code below
module.exports = { controller };
