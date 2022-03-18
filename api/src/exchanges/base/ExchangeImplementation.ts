import { ExchangeBalance, Rate } from "../../types/Exchange";

export default interface ExchangeImplementation {
  getBalances: () => Promise<ExchangeBalance[]>
  getRates: (assets: string[]) => Promise<Rate[]>
}
