import ExchangeImplementation from "./base/ExchangeImplementation";
import BinanceImplementation from "./BinanceImplementation";

export function getExchangeImplementation(exchange: string, apiKey?: string, apiSecret?: string): ExchangeImplementation {
  if(exchange === 'binance') return new BinanceImplementation(apiKey, apiSecret);
  return undefined;
}
