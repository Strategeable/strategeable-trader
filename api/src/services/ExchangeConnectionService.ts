import { ObjectId } from 'mongodb';
import ExchangeConnectionModel from '../models/ExchangeConnection';
import { ExchangeConnection } from "../types/Exchange";

export async function createExchangeConnection(conn: ExchangeConnection): Promise<ExchangeConnection> {
  try {
    const model = new ExchangeConnectionModel(conn);
    return model.save();
  } catch(err) {
    console.error(err);
    return undefined;
  }
}

export async function getExchangeConnections(userId: ObjectId): Promise<ExchangeConnection[]> {
  try {
    return ExchangeConnectionModel.find({ userId });
  } catch(err) {
    console.error(err);
    return undefined;
  }
}

export async function deleteExchangeConnection(id: ObjectId): Promise<boolean> {
  try {
    await ExchangeConnectionModel.deleteOne({ _id: id });
    return true;
  } catch(err) {
    console.error(err);
    return false;
  }
}
