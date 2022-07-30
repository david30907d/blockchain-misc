const { TokenAddress } = require("../models/tokenAddress.js");

async function controller(req, res, next) {
  new TokenAddress({address: req.body.address}).save().then(() => {
    const msg = `Successfully saved ${req.body.address} into Mongo`
    console.log(msg)
    return res.status(200).json({ "status": 200, "message": msg });
  })
  .catch((err) => {
    return res.status(400).json({ "status": 400, "message": err });
  });
}

// add the code below
module.exports = { controller };
