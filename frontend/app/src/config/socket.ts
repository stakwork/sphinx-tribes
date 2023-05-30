import { getHost } from './host';

const URL =
  process.env.NODE_ENV === 'production'
    ? `wss://${getHost()}/websocket`
    : `ws://127.0.0.1:5005/websocket`;

export const socket = new WebSocket(URL);
