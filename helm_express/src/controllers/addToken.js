const { TokenAddress } = require("../models/tokenAddress.js");

async function controller(req, res, next) {
  try {
    await new TokenAddress({ address: req.body.address }).save();
    const msg = `Successfully saved ${req.body.address} into Mongo`;
    console.log(msg);
    return res.status(200).json({ status: 200, message: msg });
  } catch (error) {
    return res.status(400).json({ status: 400, message: `${error}` });
  }
}

// add the code below
module.exports = { controller };
