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

export async function getBacktestIdsByStrategyId(id: ObjectId, skip: number = 0): Promise<{ _id: ObjectId; }[]> {
  try {
    return BacktestModel.find({ strategyId: id }).sort({
      startedOn: -1
    }).skip(skip).select('id');
  } catch(err) {
    console.error(err);
    return null;
  }
}

export async function deleteBacktestById(ids: ObjectId[]): Promise<void> {
  try {
    await BacktestModel.deleteMany({
      _id: {
        $in: ids
      }
    }).exec();
  } catch(err) {
    console.error(err);
    return null;
  }
}

export async function getBacktestsById(id: ObjectId): Promise<Backtest> {
  try {
    return BacktestModel.findOne({ _id: id });
  } catch(err) {
    console.error(err);
    return null;
  }
}

export async function createBacktest(backtest: Backtest): Promise<Backtest> {
  try {
    const model = new BacktestModel(backtest);
    return model.save();
  } catch(err) {
    console.error(err);
    return undefined;
  }
}
