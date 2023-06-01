import { getHost } from './host';

export const URL =
  process.env.NODE_ENV === 'production'
    ? `wss://${getHost()}/websocket`
    : `ws://127.0.0.1:5005/websocket`;

export const socket = new WebSocket(URL);

export const SOCKET_MSG = {
   keysend_error: 'keysend_error',
   keysend_success: 'keysend_success',
   invoice_success: 'invoice_success',
   assign_success: 'assign_success'
};

