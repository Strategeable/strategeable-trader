import { Request, Response } from "express";
import { ObjectId } from "mongodb";
import { getBacktestsByStrategyId } from "../services/BacktestService";

export async function handleGetBacktestsByStrategyId(req: Request, res: Response) {
  const { strategyId } = req.params;

  const backtests = await getBacktestsByStrategyId(new ObjectId(strategyId));
  if(!backtests) {
    return res.sendStatus(500);
  }

  return res.json(backtests);
}

export async function handleRunBacktest(req: Request, res: Response) {
  const { strategyId, startBalance, startTimestamp, endTimestamp, symbols } = req.body;

  return res.sendStatus(200);
}
