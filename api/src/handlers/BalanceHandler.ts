
import { Response, Router } from "express";
import crypto from 'crypto-js';

import { getExchangeConnections } from "../services/ExchangeConnectionService";
import ServerRequest from "../types/ServerRequest";
import { ExchangeBalance } from "../types/Exchange";
import { getExchangeImplementation } from "../exchanges";
import { singleton } from "tsyringe";
import RequestHandler from "../common/RequestHandler";

@singleton()
export default class BalanceHandler implements RequestHandler {

  route(router: Router): void {
    router.get('/', this.handleGetExchangeBalances.bind(this));
  }

  async handleGetExchangeBalances(req: ServerRequest, res: Response) {
    try {
      const exchangeConnections = await getExchangeConnections(req.user._id);
      let balances: ExchangeBalance[] = []

      for(const conn of exchangeConnections) {
        const secretKey = crypto.AES.decrypt(conn.apiSecret, process.env.ENCRYPTION_KEY).toString(crypto.enc.Utf8);

        let passPhrase: string | undefined
        if(conn.passPhrase) passPhrase = crypto.AES.decrypt(conn.passPhrase, process.env.ENCRYPTION_KEY).toString(crypto.enc.Utf8);

        const impl = getExchangeImplementation(conn.exchange, conn.apiKey, secretKey, passPhrase);

        balances = [...balances, ...await impl.getBalances()];
      }

      return res.json(balances)
    } catch(err) {
      console.error(err);
      return res.sendStatus(500);
    }
  }
}
