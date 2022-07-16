
import { Response } from "express";

import ServerRequest from "../types/ServerRequest";
import { getPositionsByBotIds } from "../services/PositionService";
import { getBotsByUserId } from "../services/BotService";

export async function handleGetPositions(req: ServerRequest, res: Response) {
  try {
    const bots = await getBotsByUserId(req.user._id);
    if(bots.length === 0) return res.json([]);

    const positions = await getPositionsByBotIds(bots.map(b => b.id));
    return res.json(positions);
  } catch(err) {
    console.error(err);
    return res.sendStatus(500);
  }
}

export async function handleGetOpenPositions(req: ServerRequest, res: Response) {
  try {
    const bots = await getBotsByUserId(req.user._id);
    if(bots.length === 0) return res.json([]);

    const positions = await getPositionsByBotIds(bots.map(b => b.id), true);
    return res.json(positions);
  } catch(err) {
    console.error(err);
    return res.sendStatus(500);
  }
}
