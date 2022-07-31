var mongoose = require("mongoose");

var tokenAddressSchema = mongoose.Schema({
  // address:  { type: String, ref: "Address", required: true, unique: true},
  address: { type: String, required: true, unique: true },
});

module.exports.TokenAddress = mongoose.model(
  "TokenAddress",
  tokenAddressSchema
);
// module.exports.tokenAddressSchema = tokenAddressSchema
