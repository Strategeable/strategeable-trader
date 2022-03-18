import AmqpConnection from "../common/AmqpConnection";
import { singleton } from "tsyringe";
import { Server } from "http";
import socketio, { Socket } from 'socket.io';
import { ConfirmChannel } from "amqplib";
import { verify } from "jsonwebtoken";
import { getUserById } from "../services/UserService";
import { getBacktestsById } from "../services/BacktestService";
import { ObjectId } from "mongodb";
import { getStrategyById } from "../services/StrategyService";

@singleton()
export default class Websocket {

    private io: socketio.Server;

    constructor(private amqpConnection: AmqpConnection, server: Server) {
        this.io = new socketio.Server(server, {
            cors: {
                origin: 'http://localhost:8080',
                credentials: false
            }
        });
        this.io.on('connection', this.handleConnection.bind(this));

        amqpConnection.getChannel().addSetup(async (channel: ConfirmChannel) => {
            const queue = await channel.assertQueue('', { durable: false, autoDelete: true, exclusive: true });

            await channel.bindQueue(queue.queue, 'backtest_x', 'backtests.*');

            channel.consume(queue.queue, message => {
                this.io.in(message.fields.routingKey).emit('backtestEvent', {
                    id: message.fields.routingKey.split('.')[1],
                    event: JSON.parse(message.content.toString())
                });
            });
        });
    }

    private handleConnection(client: Socket) {
        client.on('authorization', async token => {
            try {
                const decoded: any = verify(token, process.env.JWT_SECRET);
                const user = await getUserById(decoded.userId);
                if(!user) {
                    client.emit('authResult', false);
                    return;
                }
            
                client.data.user = user;
                client.emit('authResult', true);

                this.initClient(client);

                console.log(`${user.username} connected to websocket.`);
            } catch(err) {
                client.emit('authResult', false);
            }
        });
    }

    private initClient(client: Socket) {
        client.on('subscribeBacktest', async backtestId => {
            const backtest = await getBacktestsById(new ObjectId(backtestId));
            if(!backtest || backtest.finished) return;
            
            const strategy = await getStrategyById(backtest.strategy._id);

            if(!strategy || strategy.creator.toString() !== client.data.user._id.toString()) return;

            client.join(`backtests.${backtestId}`);
        });
    }

}