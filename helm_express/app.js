'use strict';
const express = require('express');
const marketcap = require("./src/controllers/marketcap");
// const addToken = require("./src/controllers/addToken");

// Constants
const PORT = 8080;
const HOST = '0.0.0.0';
// App
const app = express();

app.get('/marketcap/:tokenAddress', marketcap.controller);
// app.get('/token', addToken.controller);
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);