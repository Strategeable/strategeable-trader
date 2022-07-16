import ExchangeManager from '@/types/ExchangeManager'
import EventEmitter from './EventEmitter'

interface Ticker {
  baseAsset: string
  quoteAsset: string
}

export default class BinanceHandler extends EventEmitter implements ExchangeManager {
  private conn: WebSocket
  private tickerMapping: Record<string, Ticker>

  constructor (public tickers: string[]) {
    super()

    this.conn = new WebSocket('wss://stream.binance.com:9443/stream')
    this.tickerMapping = {}

    const streams = tickers.map(ticker => {
      this.tickerMapping[ticker.toLowerCase().replace('/', '')] = {
        baseAsset: ticker.split('/')[0],
        quoteAsset: ticker.split('/')[1]
      }

      return ticker.toLowerCase().replace('/', '') + '@aggTrade'
    })

    this.conn.onopen = () => {
      this.conn.send(JSON.stringify({
        method: 'SUBSCRIBE',
        params: streams,
        id: 1
      })
      )
    }

    this.conn.onmessage = trade => {
      if (JSON.parse(trade.data).id) return

      const data = JSON.parse(trade.data).data
      const ticker = this.tickerMapping[data.s.toLowerCase()]

      this.emit('TRADE', {
        rate: Number(data.p),
        quoteAsset: ticker.quoteAsset,
        baseAsset: ticker.baseAsset
      })
    }
  }

  addTicker (ticker: string) {
    // todo
  }

  removeTicker (ticker: string) {
    // todo
  }

  stop () {
    this.conn.close()
  }
}
