import { ObjectId } from "mongodb"

export default interface User {
  _id?: ObjectId
  username: string
  password: string
}
