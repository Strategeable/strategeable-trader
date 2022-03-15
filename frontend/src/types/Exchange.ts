export type Exchange = 'binance' | 'kucoin'

export const exchanges: Exchange[] = ['binance', 'kucoin']

interface ExchangeDetails {
  icon: string
  logo: string
}

export const exchangeDetails: Record<Exchange, ExchangeDetails> = {
  binance: {
    icon: '/img/binance.png',
    logo: ''
  },
  kucoin: {
    icon: '',
    logo: ''
  }
}

export interface ExchangeConnection {
  id?: string
  exchange: Exchange
  name: string
  createdOn: string
  apiKey: string
}
