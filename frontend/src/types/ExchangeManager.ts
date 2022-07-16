export default interface ExchangeManager {
  addTicker: (ticker: string) => void
  removeTicker: (ticker: string) => void
  stop: () => void
}
