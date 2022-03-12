import { Response } from "express";
import { ObjectId } from "mongodb";
import { JsonRpc } from 'node-jsonrpc-client';

import ServerRequest from "../types/ServerRequest";
import { createBacktest, getBacktestsById, getBacktestsByStrategyId } from "../services/BacktestService";
import { getStrategyById } from "../services/StrategyService";

export async function handleGetBacktestsById(req: ServerRequest, res: Response) {
  const { backtestId } = req.params;

  const backtest = await getBacktestsById(new ObjectId(backtestId));
  if(!backtest) {
    return res.sendStatus(500);
  }

  return res.json(backtest);
}

export async function handleGetBacktestsByStrategyId(req: ServerRequest, res: Response) {
  const { id } = req.params;

  const backtests = await getBacktestsByStrategyId(new ObjectId(id));
  if(!backtests) {
    return res.sendStatus(500);
  }

  return res.json(backtests);
}

export async function handleRunBacktest(req: ServerRequest, res: Response) {
  const { strategyId, startBalance, fromDate, toDate } = req.body;
  const strategy = await getStrategyById(strategyId);
  if(!strategy) return res.sendStatus(400);

  const backtest = await createBacktest({
    startedOn: new Date(),
    strategy: strategy,
    fromDate,
    toDate,
    startBalance: Number(startBalance),
    finished: false
  });

  if(!backtest) return res.sendStatus(500);

  const client = new JsonRpc(process.env.RPC_ENDPOINT);

  try {
    const result = await client.call('Backtest.Backtest', [ backtest._id.toString() ]);
    if(result.error) {
      console.error(result.error);
      return res.sendStatus(500);
    }
    res.json({ backtestId: backtest._id });
  } catch(err) {
    console.error(err);
    res.sendStatus(500);
  }
}
