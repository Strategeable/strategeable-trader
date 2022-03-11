import { Request, Response } from "express";
import { ObjectId } from "mongodb";
import { createStrategy, getStrategies, updateStrategy } from "../services/StrategyService";

export async function handleCreateStrategy(req: Request, res: Response) {
  const { strategy } = req.body;
  if(!strategy) return res.sendStatus(400);

  strategy.creator = new ObjectId();

  const strat = await createStrategy(strategy);
  if(!strat) {
    return res.sendStatus(400);
  }

  return res.json(strat);
}

export async function handleUpdateStrategy(req: Request, res: Response) {
  const { strategy } = req.body;
  if(!strategy || !strategy.id) return res.sendStatus(400);

  const strat = await updateStrategy(strategy.id, strategy);
  if(!strat) {
    return res.sendStatus(400);
  }

  return res.json(strat);
}

export async function handleGetStrategiesByUser(req: Request, res: Response) {
  try {
    const strategies = await getStrategies();
    return res.json(strategies);
  } catch(err) {
    return res.sendStatus(500);
  }
}
