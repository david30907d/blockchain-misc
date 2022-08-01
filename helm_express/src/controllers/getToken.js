const { TokenAddress } = require("../models/tokenAddress.js");

async function controller(req, res, next) {
  /*
  return all the token addresses in DB
  */
  TokenAddress.find({})
    .then((collections) => {
      let addresses = [];
      for (const collection of collections) {
        addresses.push(collection.address);
      }
      return res.status(200).json({ status: 200, addresses: addresses });
    })
    .catch((err) => {
      return res.status(400).json({ status: 400, message: err });
    });
}

module.exports = { controller };
