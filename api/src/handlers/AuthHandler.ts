import { Request, Response } from "express";
import bcrypt from 'bcryptjs';
import { sign } from 'jsonwebtoken';

import { getUserByUsername } from "../services/UserService";

export async function handleLogin(req: Request, res: Response) {
  const { username, password } = req.body;
  if(!username || !password) {
    return res.sendStatus(400);
  }

  try {
    const user = await getUserByUsername(username);
    const valid = await bcrypt.compare(password, user.password);
    if(!valid) return res.sendStatus(401);

    const token = sign({ userId: user._id }, process.env.JWT_SECRET, {
      expiresIn: '1d'
    });

    return res.json({ token });
  } catch(err) {
    console.error(err);
    return res.sendStatus(500);
  }
}
