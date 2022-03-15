
import { Response } from "express";
import crypto from 'crypto-js';

import { createExchangeConnection, getExchangeConnections } from "../services/ExchangeConnectionService";
import ServerRequest from "../types/ServerRequest";

export async function handleCreateExchangeConnection(req: ServerRequest, res: Response) {
  try {
    const { exchange, name, apiKey } = req.body;
    if (!exchange || !name || !apiKey) return res.sendStatus(400);

    const connections = await getExchangeConnections(req.user._id);
    if (connections.find(c => c.exchange === exchange && c.name === name)) return res.sendStatus(409);

    const encryptedApiKey = crypto.AES.encrypt(apiKey, process.env.ENCRYPTION_KEY).toString()

    const connection = await createExchangeConnection({
      createdOn: new Date(),
      exchange,
      name,
      apiKey: encryptedApiKey,
      userId: req.user._id
    });

    return res.json(connection);
  } catch(err) {
    console.error(err);
    return res.sendStatus(500);
  }
}
