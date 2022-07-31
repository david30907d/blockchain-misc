const CoinGecko = require("coingecko-api");
const RedisClient = require("../models/redis.js").client;
async function controller(req, res, next) {
  await RedisClient.connect();
  const tokenName = req.query.tokenName;
  const CoinGeckoClient = new CoinGecko();
  const currentUnixTimeStamp = Date.now() / 1000;
  const threeDaysAgo = currentUnixTimeStamp - 86400 * 3;
  let priceHistory;
  try {
    const cacheResult = await RedisClient.get(tokenName);
    if (cacheResult === null) {
      // cache doesn't exist!
      let data = await CoinGeckoClient.coins.fetchMarketChartRange(tokenName, {
        from: threeDaysAgo,
        to: currentUnixTimeStamp,
      });
      priceHistory = data.data.prices;
      cacheResp = {
        status: 200,
        priceHistory: priceHistory,
        cacheStatus: true,
      };
      await RedisClient.set(tokenName, JSON.stringify(cacheResp));
      cacheResp.cacheStatus = false;
      resp = cacheResp;
    } else {
      // hit cache!
      resp = JSON.parse(await RedisClient.get(tokenName));
    }
    await RedisClient.disconnect();
    return res.status(200).json(resp);
  } catch (error) {
    console.log(error);
    return res.status(500).json({ status: 500, error: `${error}` });
  }
}
module.exports = { controller };
