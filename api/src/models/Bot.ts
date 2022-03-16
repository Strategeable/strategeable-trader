import mongoose, { Schema } from "mongoose";
import Bot from "../types/Bot";

const schema: Schema<Bot> = new Schema({
  exchangeConnectionId: {
    type: mongoose.Schema.Types.ObjectId
  },
  startBalance: {
    type: Number,
    required: true
  },
  currentBalance: {
    type: Number,
    required: true
  },
  startDate: {
    type: Date,
    required: true
  },
  endDate: {
    type: Date
  },
  status: {
    type: String,
    enum: ['online', 'offline', 'ended'],
    required: true
  },
  strategy: {
    type: mongoose.Schema.Types.Mixed,
    required: true
  },
  type: {
    type: String,
    enum: ['TEST', 'LIVE'],
    required: true
  },
  userId: {
    type: mongoose.Schema.Types.ObjectId,
    required: true
  },
  quoteCurrency: {
    type: String,
    required: true
  }
});

schema.set('toJSON', {
  virtuals: true,
  versionKey: false,
  transform: (doc, ret) => {
    ret.id = ret._id.toString();
    delete ret._id;
  }
});

const model = mongoose.model('bot', schema);

export default model;
