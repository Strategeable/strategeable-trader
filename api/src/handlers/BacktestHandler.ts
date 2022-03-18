import { Response, Router } from "express";
import { ObjectId } from "mongodb";

import ServerRequest from "../types/ServerRequest";
import { createBacktest, deleteBacktestById, getBacktestIdsByStrategyId, getBacktestsById, getBacktestsByStrategyId } from "../services/BacktestService";
import { getStrategyById } from "../services/StrategyService";
import { singleton } from "tsyringe";
import RequestHandler from "../common/RequestHandler";
import AmqpConnection from "../common/AmqpConnection";

@singleton()
export default class BacktestHandler implements RequestHandler {

  constructor(private amqpConnection: AmqpConnection) {}

  route(router: Router): void {
    router.get('/strategy/:id', this.handleGetBacktestsByStrategyId.bind(this));
    router.get('/:backtestId', this.handleGetBacktestsById.bind(this));
    router.post('/:backtestId/stop', this.handleStopBacktest.bind(this));
    router.post('/', this.handleRunBacktest.bind(this));
  }

  async handleStopBacktest(req: ServerRequest, res: Response) {
    const { backtestId } = req.params;
  
    const backtest = await getBacktestsById(new ObjectId(backtestId));
    if(!backtest) return res.sendStatus(404);
    if(backtest.finished) return res.status(412).send('already finished');

    this.amqpConnection.stopBacktest(backtestId);
    await deleteBacktestById([ backtest._id ]);
  
    return res.sendStatus(204);
  }

  async handleGetBacktestsById(req: ServerRequest, res: Response) {
    const { backtestId } = req.params;
  
    const backtest = await getBacktestsById(new ObjectId(backtestId));
    if(!backtest) {
      return res.sendStatus(404);
    }
  
    return res.json(backtest);
  }
  
  async handleGetBacktestsByStrategyId(req: ServerRequest, res: Response) {
    const { id } = req.params;
  
    const backtests = await getBacktestsByStrategyId(new ObjectId(id));
    if(!backtests) {
      return res.sendStatus(500);
    }
  
    return res.json(backtests);
  }
  
  async handleRunBacktest(req: ServerRequest, res: Response) {
    const { strategyId, startBalance, fromDate, toDate } = req.body;
    const strategy = await getStrategyById(strategyId);
    if(!strategy) return res.sendStatus(400);
    if(strategy.creator.toString() !== req.user._id.toString()) return res.sendStatus(403);
  
    const backtest = await createBacktest({
      startedOn: new Date(),
      strategy: strategy,
      fromDate,
      toDate,
      startBalance: Number(startBalance),
      finished: false,
      positions: []
    });
  
    if(!backtest) return res.sendStatus(500);

    const backtestIdsToDelete = await getBacktestIdsByStrategyId(strategy._id, 5);

    if(backtestIdsToDelete.length > 0) {
      await deleteBacktestById(backtestIdsToDelete.map(b => b._id));
    }
  
    this.amqpConnection.queueBacktest(backtest._id.toString());
    
    res.json(backtest);
  }
  
}