import { Request } from "express";
import User from "./User";

export default interface ServerRequest extends Request {
  user?: User
}
