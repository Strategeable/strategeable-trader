# Strategeable Trader
![logo](https://i.imgur.com/ebYNGRp.png)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE.txt)

__DISCLAIMER__: This project is a work in progress. Feel free to contribute.

*Advanced (crypto) trading bot that enables you to create complex trading strategies and automate them*

Strategeable Trader is an open source effort to create an advanced, yet relatively easy-to-use, crypto trading bot. Simply build buy and sell paths, backtest the strategy against historical data and launch a bot on your favorite exchange.

![preview](https://i.imgur.com/8DXtjkY.png)

**Current features**
* Build/design strategies
  * Chunks (use a block of conditions in multiple paths)
  * Variables (e.g. you use the same moving average periods, change them in only one place)
  * Buy/sell paths with signal tiles, any signal tiles and chunks
* Backtest strategy
  * Backtest across multiple symbols at once
  * Chart/results in interface when finished

**Upcoming features / in progress**
* Configure exchange API keys in your account
* Start an actual bot
  * Discord/Telegram integration?

## Installation & usage
### Prerequisites
* Golang is installed on device
* Yarn globally installed
### Development environment
1. `git clone <your-fork>` 
2. Install each part of the project
    * `cd ./frontend && yarn install`
    * `cd ./api && yarn install`
3. Setup environment variables using the `.env.example` files
    * frontend: `.env.local`
    * api: `.env`
    * bot: `.env`
4. Run each component
    * frontend: `yarn serve`
    * api: `yarn dev`
    * bot: `go run main.go`

## Contributing
Pull requests are welcome. Please feel free to create issues for actual issues or potential additions to the project.