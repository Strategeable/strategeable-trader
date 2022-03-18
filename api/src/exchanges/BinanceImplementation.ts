import axios from 'axios';
import Client, { Binance } from 'binance-api-node'

import { Exchange, ExchangeBalance, Rate } from "../types/Exchange";
import ExchangeImplementation from "./base/ExchangeImplementation";

export default class BinanceImplementation implements ExchangeImplementation {

  private client: Binance

  constructor(apiKey: string, apiSecret: string) {
    this.client = Client({
      apiKey,
      apiSecret
    })
  }

  async getBalances(): Promise<ExchangeBalance[]> {
    try {
      const account = await this.client.accountInfo();
      return account.balances
        .map(b => ({ exchange: 'binance' as Exchange, asset: b.asset, amount: Number(b.free) + Number(b.locked) }))
        .filter(b => b.amount > 0);
    } catch(err) {
      console.error(err)
      return []
    }
  }

  async getRates(assets: string[]): Promise<Rate[]> {
    const prices = await axios.get('https://api.binance.com/api/v3/ticker/price');
    return []
  }
}
