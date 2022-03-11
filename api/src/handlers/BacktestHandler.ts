import { Response } from "express";
import { ObjectId } from "mongodb";
import ServerRequest from "../types/ServerRequest";
import { getBacktestsByStrategyId } from "../services/BacktestService";

export async function handleGetBacktestsByStrategyId(req: ServerRequest, res: Response) {
  const { strategyId } = req.params;

  const backtests = await getBacktestsByStrategyId(new ObjectId(strategyId));
  if(!backtests) {
    return res.sendStatus(500);
  }

  return res.json(backtests);
}

export async function handleRunBacktest(req: ServerRequest, res: Response) {
  const { strategyId, startBalance, startTimestamp, endTimestamp, symbols } = req.body;

  return res.sendStatus(200);
}
