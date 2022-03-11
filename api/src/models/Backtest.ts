import mongoose, { Schema } from "mongoose";
import Backtest from "../types/Backtest";

const schema: Schema<Backtest> = new Schema({
  strategy: {
    type: mongoose.Schema.Types.Mixed,
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
  },
  trades: {
    type: [
      {
        symbol: {
          type: String,
          required: true
        },
        amountIn: {
          type: Number,
          required: true
        },
        amountOut: {
          type: Number,
          required: true
        },
        entryPrice: {
          type: Number,
          required: true
        },
        exitPrice: {
          type: Number,
          required: true
        },
        entryDate: {
          type: Date,
          required: true
        },
        exitDate: {
          type: Date,
          required: true
        },
        fees: {
          type: Number,
          required: true
        },
        buyPathName: {
          type: String,
          required: true
        },
        sellPathName: {
          type: String,
          required: true
        }
      }
    ]
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
