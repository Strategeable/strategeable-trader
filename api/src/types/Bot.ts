import { ObjectId } from "mongodb"
import Strategy from "./Strategy"

export default interface Bot {
  id?: ObjectId
  exchangeConnectionId?: ObjectId
  userId: ObjectId
  type: 'TEST' | 'LIVE'
  strategy: Strategy
  startBalance: number
  startDate: Date
  status: 'online' | 'offline' | 'ended'
}
