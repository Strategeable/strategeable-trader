import { Request, Response } from "express";
import bcrypt from 'bcryptjs';
import { sign } from 'jsonwebtoken';

import { createUser, getUserByUsername } from "../services/UserService";

export async function handleLogin(req: Request, res: Response) {
  const { username, password } = req.body;
  if(!username || !password) return res.sendStatus(400);

  try {
    const user = await getUserByUsername(username);
    if(!user) return res.sendStatus(401);

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

export async function handleRegistration(req: Request, res: Response) {
  const { username, password } = req.body;
  if(!username || !password) return res.sendStatus(400);

  try {
    const existingUser = await getUserByUsername(username)
    if(existingUser) return res.status(409).json({ error: 'Username already taken' })

    const salt = await bcrypt.genSalt(10);
    const hash = await bcrypt.hash(password, salt);

    const user = await createUser(username, hash)
    const token = sign({ userId: user._id }, process.env.JWT_SECRET, {
      expiresIn: '1d'
    });

    return res.json({ token });
  } catch(err) {
    console.error(err)
    return res.sendStatus(500);
  }
}
