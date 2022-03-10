import { ObjectId } from "mongodb";
import BacktestModel from '../models/Backtest';
import Backtest from "../types/Backtest";

export async function getBacktestsByStrategyId(id: ObjectId): Promise<Backtest[]> {
  try {
    return BacktestModel.find({ strategyId: id });
  } catch(err) {
    console.error(err);
    return null;
  }
}
