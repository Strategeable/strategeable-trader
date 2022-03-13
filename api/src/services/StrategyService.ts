import { ObjectId } from "mongodb";
import Strategy from "../types/Strategy";
import StrategyModel from '../models/Strategy';

export async function getStrategies(creator: ObjectId): Promise<Strategy[] | null> {
  try {
    return StrategyModel.find({ creator });
  } catch(err) {
    console.error(err);
    return null;
  }
}

export async function createStrategy(strategy: Strategy, creator: ObjectId): Promise<Strategy> {
  try {
    strategy.creator = creator;
    return await StrategyModel.create(strategy);
  } catch(err) {
    console.error(err);
    return undefined;
  }
}

export async function updateStrategy(_id: ObjectId, strategyUpdate: any): Promise<any> {
  try {
    const existing: any = await StrategyModel.findOne({ _id })
    if(!existing) return null;


    for(const [key, val] of Object.entries(strategyUpdate)) {
      existing[key] = val;
    }

    existing.lastEdited = new Date()

    await existing.save();
    return existing;
  } catch(err) {
    console.error(err);
    return undefined;
  }
}

export async function getStrategyById(_id: ObjectId): Promise<Strategy> {
  try {
    return StrategyModel.findOne({ _id });
  } catch(err) {
    return err
  }
}
