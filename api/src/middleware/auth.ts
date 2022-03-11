import { NextFunction, Response } from "express";
import { verify } from 'jsonwebtoken';
import ServerRequest from "../types/ServerRequest";
import { getUserById } from "../services/UserService";

export default async function (req: ServerRequest, res: Response, next: NextFunction) {
  const authHeader = req.headers['authorization'];
  if(!authHeader) return res.sendStatus(401);
  if(!authHeader.startsWith('Bearer ')) return res.sendStatus(401);

  const token = authHeader.substring(7);
  try {
    const decoded: any = verify(token, process.env.JWT_SECRET);
    const user = await getUserById(decoded.userId);
    if(!user) return res.sendStatus(401);

    req.user = user;
    next();
  } catch(err) {
    return res.sendStatus(401);
  }
}
