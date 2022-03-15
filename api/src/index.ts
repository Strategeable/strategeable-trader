import express from 'express'
import mongoose from 'mongoose';
import dotenv from 'dotenv';
import cors from 'cors';
dotenv.config();

import { handleGetBacktestsById, handleGetBacktestsByStrategyId, handleRunBacktest } from './handlers/BacktestHandler';
import { handleCreateStrategy, handleGetStrategies, handleGetStrategyById, handleUpdateStrategy } from './handlers/StrategyHandler';
import { handleLogin, handleRegistration } from './handlers/AuthHandler';
import { handleCreateExchangeConnection, handleGetExchangeConnections } from './handlers/ExchangeHandler';
import auth from './middleware/auth';

const app = express();
const port = process.env.PORT || 3000;

(async() => {
  await mongoose.connect(process.env.DATABASE_URL);

  app.use(express.json());
  app.use(cors());

  app.post('/login', handleLogin);
  app.post('/register', handleRegistration);

  app.use(auth);

  app.get('/backtest/strategy/:id', handleGetBacktestsByStrategyId);
  app.get('/backtest/:backtestId', handleGetBacktestsById);
  app.post('/backtest', handleRunBacktest);

  app.get('/strategy', handleGetStrategies);
  app.get('/strategy/:id', handleGetStrategyById);
  app.post('/strategy', handleCreateStrategy);
  app.put('/strategy', handleUpdateStrategy);

  app.get('/settings/exchange-connection', handleCreateExchangeConnection);
  app.post('/settings/exchange-connection', handleGetExchangeConnections);

  app.listen(port, () => console.log(`Server running on port ${port}`));
})();
