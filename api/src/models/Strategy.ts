import mongoose, { Schema } from "mongoose";
import Strategy from "../types/Strategy";

const schema: Schema<Strategy> = new Schema({
  creator: {
    type: mongoose.Schema.Types.ObjectId,
    required: true
  },
  createdAt: {
    type: Date,
    required: true
  },
  exchange: {
    type: String,
    enum: ['binance'],
    required: true
  },
  lastEdited: {
    type: Date,
    required: true
  },
  name: {
    type: mongoose.Schema.Types.String,
    required: true
  },
  symbols: [String],
  chunks: [
    {
      id: {
        type: String,
        required: true
      },
      name: {
        type: String,
        required: true
      },
      steps: [
        {
          id: String,
          type: {
            type: String,
            enum: ['SIGNAL_TILE', 'ANY_SIGNAL_TILE', 'CHUNK_ID']
          },
          data: mongoose.Schema.Types.Mixed
        }
      ]
    }
  ],
  paths: [
    {
      id: {
        type: String,
        required: true
      },
      name: {
        type: String,
        required: true
      },
      type: {
        type: String,
        enum: ['BUY', 'SELL'],
        required: true
      },
      whitelist: [String],
      steps: {
        type: [
          {
            id: {
              type: String,
              required: true
            },
            type: {
              type: String,
              enum: ['SIGNAL_TILE', 'ANY_SIGNAL_TILE', 'CHUNK_ID'],
              required: true
            },
            data: {
              type: mongoose.Schema.Types.Mixed,
              required: true
            }
          }
        ],
        required: true
      }
    }
  ],
  variables: [
    {
      id: {
        type: String,
        required: true
      },
      type: {
        type: String,
        enum: ['number', 'timeframe'],
        required: true
      },
      key: {
        type: String,
        required: true
      },
      value: {
        type: mongoose.Schema.Types.Mixed
      }
    }
  ]
});

schema.set('toJSON', {
  virtuals: true,
  versionKey: false,
  transform: (doc, ret) => {
    ret.id = ret._id.toString();
    delete ret._id;
  }
});

const model = mongoose.model('strategy', schema);

export default model;
