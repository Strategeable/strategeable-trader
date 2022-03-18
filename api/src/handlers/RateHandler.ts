
import { Response, Router } from "express";

import ServerRequest from "../types/ServerRequest";
import { getExchangeImplementation } from "../exchanges";
import { singleton } from "tsyringe";
import RequestHandler from "../common/RequestHandler";

@singleton()
export default class RateHandler implements RequestHandler {

  route(router: Router): void {
    router.get('/', this.handleGetRates.bind(this));
  }

  async handleGetRates(req: ServerRequest, res: Response) {
    try {
      const { exchange, coins } = req.query;
      if(!exchange || !coins) return res.sendStatus(400);

      const impl = getExchangeImplementation(exchange as string);
      if(!impl) return res.sendStatus(400);

      const rates = await impl.getRates((coins as string).split(','));
      return res.json(rates)
    } catch(err) {
      console.error(err);
      return res.sendStatus(500);
    }
  }
}
