
import { Response, Router } from "express";

import ServerRequest from "../types/ServerRequest";
import BotModel from '../models/Bot';
import { ObjectId } from "mongodb";
import { getStrategyById } from "../services/StrategyService";
import { getBotsByUserId } from "../services/BotService";
import { getExchangeConnectionById } from "../services/ExchangeConnectionService";
import { singleton } from "tsyringe";
import RequestHandler from "../common/RequestHandler";

@singleton()
export default class BotHandler implements RequestHandler {

  route(router: Router): void {
    router.get('/', this.handleGetBots.bind(this));
    router.post('/', this.handleCreateBot.bind(this));
  }

  async handleCreateBot(req: ServerRequest, res: Response) {
    try {
      const { type, strategyId, startBalance, exchangeConnection } = req.body;
      if(!type || !strategyId || !startBalance) return res.sendStatus(400);
  
      const strategy = await getStrategyById(new ObjectId(strategyId));
      if(!strategy) return res.sendStatus(400);
      if(!strategy.creator.equals(req.user._id)) return res.sendStatus(403);
  
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
  
      if(exchangeConnection) {
        const conn = await getExchangeConnectionById(new ObjectId(exchangeConnection));
        if(!conn.userId.equals(req.user._id)) return res.sendStatus(403);
  
        bot.exchangeConnectionId = new ObjectId(exchangeConnection);
      }
  
      await bot.save();
  
      return res.json(bot);
    } catch(err) {
      console.error(err);
      return res.sendStatus(500);
    }
  }
  
  async handleGetBots(req: ServerRequest, res: Response) {
    try {
      const bots = await getBotsByUserId(req.user._id);
      return res.json(bots);
    } catch(err) {
      console.error(err);
      return res.sendStatus(500);
    }
  }

}
