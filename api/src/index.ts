import express from 'express'
import mongoose from 'mongoose';
import dotenv from 'dotenv';
dotenv.config();

import { handleGetBacktestsByStrategyId, handleRunBacktest } from './handlers/BacktestHandler';
import { handleCreateStrategy, handleGetStrategiesByUser, handleUpdateStrategy } from './handlers/StrategyHandler';
import { handleLogin } from './handlers/AuthHandler';
import auth from './middleware/auth';
import UserModel from './models/User';

const app = express();
const port = process.env.PORT || 3000;

(async() => {
  await mongoose.connect(process.env.DATABASE_URL);

  app.use(express.json())

  app.post('/login', handleLogin)

  app.use(auth);

  app.get('/backtest/:strategyId', handleGetBacktestsByStrategyId);
  app.post('/backtest', handleRunBacktest);

  app.get('/strategy', handleGetStrategiesByUser);
  app.post('/strategy', handleCreateStrategy);
  app.put('/strategy', handleUpdateStrategy);

  app.listen(port, () => console.log(`Server running on port ${port}`));
})();
