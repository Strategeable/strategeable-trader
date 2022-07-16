import mongoose, { Schema } from "mongoose";
import Position from "../types/Position";

const schema: Schema<Position> = new Schema({
  botId: {
    type: mongoose.Schema.Types.ObjectId,
    required: true
  },
  closeTime: {
    type: Date
  },
  openTime: {
    type: Date,
    required: true
  },
  orders: {
    type: [
      {
        orderId: {
          type: String,
          required: true
        },
        time: {
          type: Date,
          required: true
        },
        side: {
          type: String,
          enum: ['BUY', 'SELL'],
          required: true
        },
        active: {
          type: Boolean,
          required: true
        },
        size: {
          type: Number,
          required: true
        },
        rate: {
          type: Number,
          required: true
        },
        fills: [
          {
            time: {
              type: Date,
              required: true
            },
            rate: {
              type: Number,
              required: true
            },
            quantity: {
              type: Number,
              required: true
            },
            quoteFee: {
              type: Number,
              required: true
            }
          }
        ]
      }
    ]
  },
  state: {
    type: String,
    enum: ['OPENING', 'OPEN', 'CLOSING', 'CLOSED'],
    required: true
  },
  symbol: {
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

const model = mongoose.model('position', schema);

export default model;
