import { ObjectId } from "mongodb";
import Strategy from "../types/Strategy";
import StrategyModel from '../models/Strategy';

export async function getStrategies(creator?: ObjectId): Promise<Strategy[] | null> {
  try {
    if(creator) return StrategyModel.find({ creator });
    return StrategyModel.find();
  } catch(err) {
    console.error(err);
    return null;
  }
}

export async function createStrategy(strategy: Strategy): Promise<Strategy> {
  try {
    return await StrategyModel.create(strategy);
  } catch(err) {
    console.error(err);
    return undefined;
  }
}

export async function updateStrategy(id: ObjectId, strategyUpdate: any): Promise<any> {
  try {
    const existing: any = await StrategyModel.findOne({ id: new ObjectId(id) })
    if(!existing) return null;


    for(const [key, val] of Object.entries(strategyUpdate)) {
      existing[key] = val;
    }

    await existing.save();
    return existing.toObject();
  } catch(err) {
    return null;
  }
}

export async function getStrategyById(id: ObjectId): Promise<Strategy> {
  try {
    return StrategyModel.findOne({ id });
  } catch(err) {
    return err
  }
}
