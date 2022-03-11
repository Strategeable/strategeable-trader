import { Response } from "express";
import ServerRequest from "../types/ServerRequest";
import { createStrategy, getStrategies, getStrategyById, updateStrategy } from "../services/StrategyService";
import { ObjectId } from "mongodb";

export async function handleCreateStrategy(req: ServerRequest, res: Response) {
  const { strategy } = req.body;
  if(!strategy) return res.sendStatus(400);

  const strat = await createStrategy(strategy, req.user._id);
  if(!strat) {
    return res.sendStatus(400);
  }

  return res.json(strat);
}

export async function handleUpdateStrategy(req: ServerRequest, res: Response) {
  const { strategy } = req.body;
  if(!strategy || !strategy.id) return res.sendStatus(400);

  const existingStrategy = await getStrategyById(new ObjectId(strategy.id));
  if(!existingStrategy) return res.sendStatus(400);
  if(!existingStrategy.creator.equals(req.user._id)) return res.sendStatus(401);

  const strat = await updateStrategy(strategy.id, strategy);
  if(!strat) {
    return res.sendStatus(400);
  }

  return res.json(strat);
}

export async function handleGetStrategies(req: ServerRequest, res: Response) {
  try {
    const strategies = await getStrategies(req.user._id);
    return res.json(strategies);
  } catch(err) {
    return res.sendStatus(500);
  }
}

export async function handleGetStrategyById(req: ServerRequest, res: Response) {
  try {
    const { id } = req.params
    const strategy = await getStrategyById(new ObjectId(id));
    return res.json(strategy);
  } catch(err) {
    return res.sendStatus(500);
  }
}
