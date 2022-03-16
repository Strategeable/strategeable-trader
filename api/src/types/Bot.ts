import { ObjectId } from "mongodb"
import Strategy from "./Strategy"

export default interface Bot {
  id?: ObjectId
  exchangeConnectionId?: ObjectId
  userId: ObjectId
  type: 'TEST' | 'LIVE'
  strategy: Strategy
  startBalance: number
  currentBalance: number
  startDate: Date
  endDate?: Date
  status: 'online' | 'offline' | 'ended'
  quoteCurrency: string
}
