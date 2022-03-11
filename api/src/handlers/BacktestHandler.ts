import { Response } from "express";
import { ObjectId } from "mongodb";
import { JsonRpc } from 'node-jsonrpc-client';

import ServerRequest from "../types/ServerRequest";
import { createBacktest, getBacktestsByStrategyId } from "../services/BacktestService";
import { getStrategyById } from "../services/StrategyService";

export async function handleGetBacktestsByStrategyId(req: ServerRequest, res: Response) {
  const { strategyId } = req.params;

  const backtests = await getBacktestsByStrategyId(new ObjectId(strategyId));
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
    strategy: strategy,
    fromDate,
    toDate,
    startBalance: Number(startBalance),
    trades: [],
    finished: false
  });

  if(!backtest) return res.sendStatus(500);

  const client = new JsonRpc(process.env.RPC_ENDPOINT);

  try {
    const response = await client.call('Backtest.Backtest', backtest._id.toString());
    console.log(response);
    res.sendStatus(200);
  } catch(err) {
    console.error(err);
    res.sendStatus(500);
  }
}
