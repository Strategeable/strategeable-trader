import ExchangeImplementation from "./base/ExchangeImplementation";
import BinanceImplementation from "./BinanceImplementation";
import KucoinImplementation from "./KucoinImplementation";

export function getExchangeImplementation(exchange: string, apiKey?: string, apiSecret?: string, passphrase?: string): ExchangeImplementation {
  if(exchange === 'binance') return new BinanceImplementation(apiKey, apiSecret);
  if(exchange === 'kucoin') return new KucoinImplementation(apiKey, apiSecret, passphrase);
  return undefined;
}
