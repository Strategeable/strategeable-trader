import { Strategy } from "./Strategy"

export default interface Bot {
  id?: string
  exchangeConnectionId?: string
  type: 'TEST' | 'LIVE'
  strategy: Strategy
  startBalance: number
  currentBalance: number
  startDate: Date
  endDate?: Date
  status: 'online' | 'offline' | 'ended'
  quoteCurrency: string
}

export interface LaunchParameters {
  type: 'LIVE' | 'TEST'
  strategyId: string
  startBalance: number
  exchangeConnection: string
}
