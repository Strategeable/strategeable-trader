import amqp, { ChannelWrapper } from "amqp-connection-manager";
import { ConfirmChannel } from "amqplib";
import { singleton } from "tsyringe";

@singleton()
export default class AmqpConnection {

    private channel: ChannelWrapper;

    constructor() {
        const connection = amqp.connect(process.env.AMQP_URL);
        this.channel = connection.createChannel({
            setup: async (channel: ConfirmChannel) => {
                await channel.assertExchange('backtest_x', 'topic', {
                    durable: false
                });
                await channel.assertQueue('backtest_q', {
                    durable: true
                });
            }
        });
    }

    getChannel(): ChannelWrapper {
        return this.channel;
    }

    queueBacktest(id: string) {
        this.channel.sendToQueue('backtest_q', Buffer.from(JSON.stringify({
            id
        })));
    }

    stopBacktest(id: string) {
        this.channel.publish('backtest_x', `backtests.${id}.control`, Buffer.from(JSON.stringify({
            action: 'STOP',
            backtestId: id
        })));
    }

}