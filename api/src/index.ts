import 'reflect-metadata';
import express, { Router } from 'express'
import mongoose from 'mongoose';
import dotenv from 'dotenv';
import cors from 'cors';
dotenv.config();

import StrategyHandler from './handlers/StrategyHandler';
import AuthHandler from './handlers/AuthHandler';
import ExchangeHandler from './handlers/ExchangeHandler';
import auth from './middleware/auth';
import { handleGetOpenPositions, handleGetPositions } from './handlers/PositionHandler';
import BotHandler from './handlers/BotHandler';
import RequestHandler from './common/RequestHandler';
import { container } from 'tsyringe';
import BacktestHandler from './handlers/BacktestHandler';
import { createServer, Server } from 'http';
import Websocket from './websocket';

const port = process.env.PORT || 3000;

(async() => {
  await mongoose.connect(process.env.DATABASE_URL);

  const app = express();
  const server = createServer(app);

  container.register(Server, {
    useValue: server
  });

  // Start websocket
  container.resolve(Websocket);

  app.use(express.json());
  app.use(cors());

  route(app, '/auth', container.resolve(AuthHandler));

  app.use(auth);

  route(app, '/backtest', container.resolve(BacktestHandler));
  route(app, '/strategy', container.resolve(StrategyHandler));
  route(app, '/bot', container.resolve(BotHandler));

  const settingsRouter = Router();
  app.use('/settings', settingsRouter);

  route(settingsRouter, '/exchange-connection', container.resolve(ExchangeHandler))

  app.get('/position', handleGetPositions);
  app.get('/position/open', handleGetOpenPositions);

  server.listen(port, () => console.log(`Server running on port ${port}`));
})();

function route(router: Router, endpoint: string, handler: RequestHandler) {
  const newRouter = Router();

  router.use(endpoint, newRouter);

  handler.route(newRouter);
}
