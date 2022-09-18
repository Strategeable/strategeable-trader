import axios from 'axios';
import Client, { Binance } from 'binance-api-node'

import { Exchange, ExchangeBalance, Rate } from "../types/Exchange";
import ExchangeImplementation from "./base/ExchangeImplementation";

export default class BinanceImplementation implements ExchangeImplementation {

  private client: Binance

  constructor(apiKey?: string, apiSecret?: string) {
    this.client = Client({
      apiKey,
      apiSecret
    })
  }

  async getBalances(): Promise<ExchangeBalance[]> {
    try {
      const account = await this.client.accountInfo({
        useServerTime: true
      });
      return account.balances
        .map(b => ({ exchange: 'binance' as Exchange, asset: b.asset, amount: Number(b.free) + Number(b.locked) }))
        .filter(b => b.amount > 0);
    } catch(err) {
      console.error(err)
      return []
    }
  }

  async getRates(assets: string[]): Promise<Rate[]> {
    try {
      const prices: { data: { symbol: string, price: string }[] } = await axios.get('https://api.binance.com/api/v3/ticker/price');
      const rates: Rate[] = [];

      for(const asset of assets) {
        const assetRates = prices.data.filter(p => p.symbol.startsWith(asset));
        const rateObj: Rate = {
          asset,
          exchange: 'binance',
          quote: {}
        }
        for(const rate of assetRates) {
          rateObj.quote[rate.symbol.slice(asset.length)] = Number(rate.price)
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
