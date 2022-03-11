import { ObjectId } from 'mongodb';
import UserModel from '../models/User';
import User from "../types/User";

export async function getUserByUsername(username: string): Promise<User | null> {
  try {
    return await UserModel.findOne({ username });
  } catch(err) {
    console.error(err);
    return null;
  }
}

export async function getUserById(id: string): Promise<User | null> {
  try {
    return await UserModel.findOne({ _id: new ObjectId(id) });
  } catch(err) {
    console.error(err);
    return null;
  }
}
