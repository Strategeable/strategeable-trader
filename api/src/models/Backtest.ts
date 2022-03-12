import mongoose, { Schema } from "mongoose";
import Backtest from "../types/Backtest";

const schema: Schema<Backtest> = new Schema({
  strategy: {
    type: mongoose.Schema.Types.Mixed,
    required: true
  },
  startedOn: {
    type: Date,
    required: true
  },
  startBalance: {
    type: Number,
    required: true
  },
  endBalance: {
    type: Number
  },
  fromDate: {
    type: Date,
    required: true
  }  ,
  toDate: {
    type: Date,
    required: true
  },
  finished: {
    type: Boolean,
    default: false
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

const model = mongoose.model('backtest', schema);

export default model;
