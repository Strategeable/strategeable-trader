import { Strategy } from "./Strategy"

export default interface Bot {
  id?: string
  exchangeConnectionId?: string
  type: 'TEST' | 'LIVE'
  strategy: Strategy
  startBalance: number
  startDate: Date
  status: 'online' | 'offline' | 'ended'
}
