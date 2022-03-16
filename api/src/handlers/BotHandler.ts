
import { Response } from "express";

import ServerRequest from "../types/ServerRequest";
import BotModel from '../models/Bot';
import { ObjectId } from "mongodb";
import { getStrategyById } from "../services/StrategyService";
import { getBotsByUserId } from "../services/BotService";

export async function handleCreateBot(req: ServerRequest, res: Response) {
  try {
    const { type, strategyId, startBalance, exchangeConnection } = req.body;
    if(!type || !strategyId || !startBalance) return res.sendStatus(400);

    const strategy = await getStrategyById(new ObjectId(strategyId));
    if(!strategy) return res.sendStatus(400);

    const bot = new BotModel({
      type,
      strategy,
      startBalance: Number(startBalance),
      currentBalance: Number(startBalance),
      startDate: new Date(),
      status: 'offline',
      userId: req.user._id,
      quoteCurrency: strategy.quoteCurrency
    });

    if(exchangeConnection) bot.exchangeConnectionId = new ObjectId(exchangeConnection);

    await bot.save();

    return res.json(bot);
  } catch(err) {
    console.error(err);
    return res.sendStatus(500);
  }
}

export async function handleGetBots(req: ServerRequest, res: Response) {
  try {
    const bots = await getBotsByUserId(req.user._id);
    return res.json(bots);
  } catch(err) {
    console.error(err);
    return res.sendStatus(500);
  }
}
