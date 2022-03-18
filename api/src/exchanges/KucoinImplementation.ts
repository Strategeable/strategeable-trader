const api = require('kucoin-node-sdk');

import { ExchangeBalance, Rate } from "../types/Exchange";
import ExchangeImplementation from "./base/ExchangeImplementation";

interface Balance {
  id: string
  currency: string
  type: string
  balance: string
  available: string
  holds: string
}

interface Ticker {
  symbol: string
  buy: string
  sell: string
  last: string
}

export default class KucoinImplementation implements ExchangeImplementation {

  private client: any

  constructor(apiKey?: string, apiSecret?: string, passphrase?: string) {
    this.client = api
    if(apiSecret) {
      this.client.init({
        baseUrl: 'https://api.kucoin.com',
        apiAuth: {
          key: apiKey,
          secret: apiSecret,
          passphrase: passphrase
        },
        authVersion: 2
      });
    } else {
      this.client.init({
        baseUrl: 'https://api.kucoin.com',
        apiAuth: {
          key: '',
          secret: '',
          passphrase: ''
        },
        authVersion: 2
      })
    }
  }

  async getBalances(): Promise<ExchangeBalance[]> {
    try {
      const response = await this.client.rest.User.Account.getAccountsList();
      if(!response.data) return []

      const balances: Balance[] = response.data;
      const mapped: ExchangeBalance[] = balances.map(b => ({
        amount: Number(b.balance),
        asset: b.currency,
        exchange: 'kucoin'
      }));

      return mapped.filter(m => m.amount > 0);
    } catch(err) {
      console.error(err)
      return []
    }
  }

  async getRates(assets: string[]): Promise<Rate[]> {
    try {
      const response = await this.client.rest.Market.Symbols.getAllTickers()
      if(!response.data) return [];

      const tickers: Ticker[] = response.data.ticker;
      const rates: Rate[] = [];

      for(const asset of assets) {
        const assetRates = tickers.filter(t => t.symbol.startsWith(asset));
        const rateObj: Rate = {
          asset,
          exchange: 'kucoin',
          quote: {}
        }
        for(const ticker of assetRates) {
          rateObj.quote[ticker.symbol.split('-')[1]] = Number(ticker.last);
        }

        rates.push(rateObj);
      }

      return rates
    } catch(err) {
      console.error(err);
      return []
    }
  }
}
