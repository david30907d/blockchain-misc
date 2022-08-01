"use strict";
const express = require("express");
const mongoose = require("mongoose");
// connect to Mongo when the app initializes
mongoose
  .connect("mongodb://mongo:27017")
  .then(() => console.log("DB connection successful!"));

const marketcap = require("./src/controllers/marketcap");
const addToken = require("./src/controllers/addToken");
const getToken = require("./src/controllers/getToken");
const getPriceHistory = require("./src/controllers/getPriceHistory");
// Constants
const PORT = 8080;
const HOST = "0.0.0.0";
// App
const app = express();
app.use(express.json());
app.get("/", async (_req, res, _next) => {res.status(200).send('OK')});
app.get("/marketcap", marketcap.controller);
app.post("/token", addToken.controller);
app.get("/token", getToken.controller);
app.get("/priceHistory", getPriceHistory.controller);
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);
