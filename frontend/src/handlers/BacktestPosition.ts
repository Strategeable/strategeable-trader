import { Position } from '@/types/Backtest'

export default class BacktestPosition {
  public position: Position

  constructor (position: Position) {
    this.position = position
  }

  // TODO: change percent should include the fees
  getChangePercent (): number {
    return (this.position.exitValue.rate - this.position.entryValue.rate) / this.position.entryValue.rate * 100
  }

  getQuoteDifference (): number {
    const quoteEntrySize = this.position.entryValue.rate * this.position.entryValue.baseSize
    const quoteExitSize = this.position.exitValue.rate * this.position.exitValue.baseSize
    return quoteExitSize - quoteEntrySize - this.position.exitValue.quoteFees
  }

  getFees (): number {
    return this.position.entryValue.quoteFees + this.position.exitValue.quoteFees
  }
}
