import mongoose, { Schema } from "mongoose";
import { ExchangeConnection } from "../types/Exchange";

const schema: Schema<ExchangeConnection> = new Schema({
  exchange: {
    type: String,
    enum: ['binance', 'kucoin'],
    required: true
  },
  name: {
    type: String,
    required: true
  },
  createdOn: {
    type: Date,
    required: true
  },
  apiKey: {
    type: String,
    required: true
  },
  apiSecret: {
    type: String,
    required: true
  },
  passPhrase: {
    type: String
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

const model = mongoose.model('exchangeconnection', schema);

export default model;
