import express from 'express'
import { ObjectId } from 'mongodb';
import { handleGetBacktestsByStrategyId, handleRunBacktest } from './handlers/BacktestHandler';
import { handleCreateStrategy, handleUpdateStrategy } from './handlers/StrategyHandler';
import Strategy from './models/Strategy'

const app = express();
const port = process.env.PORT || 3000;

(async() => {
  app.use(express.json())

  app.get('/backtest/:strategyId', handleGetBacktestsByStrategyId);
  app.post('/backtest', handleRunBacktest);

  app.post('/strategy', handleCreateStrategy);
  app.put('/strategy', handleUpdateStrategy);

  app.listen(port, () => console.log(`Server running on port ${port}`));
})();
