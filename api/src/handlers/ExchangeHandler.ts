
import { Response } from "express";
import crypto from 'crypto-js';

import { createExchangeConnection, deleteExchangeConnection, getExchangeConnections } from "../services/ExchangeConnectionService";
import ServerRequest from "../types/ServerRequest";
import { ObjectId } from "mongodb";

export async function handleCreateExchangeConnection(req: ServerRequest, res: Response) {
  try {
    const { exchange, name, apiKey, apiSecret } = req.body;
    if (!exchange || !name || !apiKey || !apiSecret) return res.sendStatus(400);

    const connections = await getExchangeConnections(req.user._id);
    if (connections.find(c => c.exchange === exchange && c.name === name)) return res.sendStatus(409);

    const encryptedSecretKey = crypto.AES.encrypt(apiSecret, process.env.ENCRYPTION_KEY).toString()

    const connection = await createExchangeConnection({
      createdOn: new Date(),
      exchange,
      name,
      apiKey: apiKey,
      apiSecret: encryptedSecretKey,
      userId: req.user._id
    });

    delete connection.apiKey

    return res.json(connection);
  } catch(err) {
    console.error(err);
    return res.sendStatus(500);
  }
}

export async function handleGetExchangeConnections(req: ServerRequest, res: Response) {
  try {
    const connections = await getExchangeConnections(req.user._id);
    connections.forEach(c => delete c.apiKey);
    return res.json(connections);
  } catch(err) {
    console.error(err);
    res.sendStatus(500);
  }
}

export async function handleDeleteExchangeConnection(req: ServerRequest, res: Response) {
  try {
    const deleted = await deleteExchangeConnection(new ObjectId(req.params.id));
    if(!deleted) return res.sendStatus(400);

    return res.sendStatus(200);
  } catch(err) {
    console.error(err);
    res.sendStatus(500);
  }
}
