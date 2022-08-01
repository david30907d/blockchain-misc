const redis = require("redis");

const client = redis.createClient({
  socket: {
    host: process.env.REDIS_HOST,
    port: process.env.REDIS_PORT,
  },
});
client.on("error", (err) => console.log("Redis Client Error", err));

const cache_expire_seconds = 3600;
module.exports.client = client;
module.exports.cache_expire_seconds = cache_expire_seconds;
