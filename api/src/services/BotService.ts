import Bot from '../types/Bot';
import BotModel from '../models/Bot';
import { ObjectId } from 'mongodb';

export async function createBot(bot: Bot): Promise<Bot | undefined> {
  try {
    return await BotModel.create(bot);
  } catch(err) {
    console.error(err)
    return undefined
  }
}

export async function getBotsByUserId(userId: ObjectId): Promise<Bot[]> {
  try {
    return await BotModel.find({ userId });
  } catch(err) {
    console.error(err)
    return undefined
  }
}
