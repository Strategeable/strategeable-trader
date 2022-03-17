import Position from '../types/Position';
import PositionModel from '../models/Position';
import { ObjectId } from 'mongodb';

export async function getPositionsByBotId(botId: ObjectId, open?: boolean): Promise<Position[]> {
  try {
    if(open) return await PositionModel.find({ botId, state: { $in: ['OPEN', 'OPENING'] } });
    return await PositionModel.find({ botId });
  } catch(err) {
    console.error(err)
    return undefined
  }
}

export async function getPositionsByBotIds(botIds: ObjectId[], open?: boolean): Promise<Position[]> {
  try {
    if(open) return await PositionModel.find({ botId: { $in: botIds }, state: { $in: ['OPEN', 'OPENING'] } });
    return await PositionModel.find({ botId: { $in: botIds } });
  } catch(err) {
    console.error(err)
    return undefined
  }
}
