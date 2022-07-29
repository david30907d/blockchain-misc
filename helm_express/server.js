'use strict';
const express = require('express');
const priceController = require("./src/controllers/priceController");
const Web3 = require('web3');
// const web3 = new Web3("https://mainnet.infura.io/v3/18f80354699e46188ac9b12df50f9296");
const web3 = new Web3("https://rinkeby.infura.io/v3/18f80354699e46188ac9b12df50f9296");

// Constants
const PORT = 8080;
const HOST = 'localhost';
// App
const app = express();
  
let correctABI = [{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"payable":true,"stateMutability":"payable","type":"fallback"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]
app.get('/', async function (req, res, next) {
    const tokenRouter = await new web3.eth.Contract( correctABI, '0x83b835c3200147e74bddf3571fde527647b126c7' );
    const totalSupply = await tokenRouter.methods.totalSupply().call();
    const decimals = await tokenRouter.methods.decimals().call();
    const tokenName = await tokenRouter.methods.name().call();
    const correctTotalSupply = totalSupply/(10**decimals)
	
    const tokenPrice = await priceController.getPriceFromChainLink("0x3539F2E214d8BC7E611056383323aC6D1b01943c")
    const marketCap = tokenPrice * correctTotalSupply
    res.send(`${tokenName} - Market Cap: ${marketCap}, price: ${tokenPrice}, `);
   });
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);