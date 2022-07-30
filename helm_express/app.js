"use strict";
const express = require("express");
const mongoose = require("mongoose");
// connect to Mongo when the app initializes
mongoose.connect('mongodb://mongo:27017').then(() => console.log('DB connection successful!'));


const marketcap = require("./src/controllers/marketcap");
const addToken = require("./src/controllers/addToken");
const getToken = require("./src/controllers/getToken");
// Constants
const PORT = 8080;
const HOST = "0.0.0.0";
// App
const app = express();
app.use(express.json());
app.get("/marketcap/:tokenAddress", marketcap.controller);
app.post("/token", addToken.controller);
app.get("/token", getToken.controller);
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);