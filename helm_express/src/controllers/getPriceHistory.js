const CoinGecko = require("coingecko-api");
const RedisClient = require("../models/redis.js").client;
const cache_expire_seconds = require("../models/redis.js").cache_expire_seconds;

async function controller(req, res, next) {
  /*
  main logic is as follow:
  1. check if cache exists
  2. get price history from CoinGecko
  3. store result in cache with cache expire seconds
  */
  await RedisClient.connect();
  const tokenName = req.query.tokenName;
  try {
    const cacheResult = await RedisClient.get(tokenName);
    if (cacheResult === null) {
      // cache doesn't exist!
      cacheResp = await getPriceHistory(RedisClient, tokenName);
      cacheResp.cacheStatus = false;
      console.log(cacheResp);
      console.log(cacheResp);
      console.log(cacheResp);
      console.log(cacheResp);
      console.log(cacheResp);
      console.log(cacheResp);
      console.log(cacheResp);
      console.log(cacheResp);
      resp = cacheResp;
    } else {
      // hit cache!
      resp = JSON.parse(await RedisClient.get(tokenName));
    }
    res.status(200).json(resp);
  } catch (error) {
    console.log(error);
    res.status(500).json({ status: 500, error: `${error}` });
  } finally {
    await RedisClient.disconnect();
  }
}

async function getPriceHistory(RedisClient, tokenName) {
  const CoinGeckoClient = new CoinGecko();
  const currentUnixTimeStamp = Date.now() / 1000;
  const threeHoursAgo = currentUnixTimeStamp - 3600 * 3;
  let data = await CoinGeckoClient.coins.fetchMarketChartRange(tokenName, {
    from: threeHoursAgo,
    to: currentUnixTimeStamp,
  });
  cacheResp = {
    status: 200,
    priceHistory: data.data.prices,
    cacheStatus: true,
  };
  await RedisClient.set(
    tokenName,
    JSON.stringify(cacheResp),
    "EX",
    cache_expire_seconds
  );
  return cacheResp;
}
module.exports = { controller };
