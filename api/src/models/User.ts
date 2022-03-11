import mongoose, { Schema } from "mongoose";
import User from "../types/User";

const schema: Schema<User> = new Schema({
  username: {
    type: String,
    required: true
  },
  password: {
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

const model = mongoose.model('user', schema);

export default model;
