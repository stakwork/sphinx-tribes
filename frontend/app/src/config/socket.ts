import { getHost } from './host';

export const URL =
  process.env.NODE_ENV !== 'development'
    ? `wss://${getHost()}/websocket`
    : `ws://127.0.0.1:5005/websocket`;

export const SOCKET_MSG = {
  keysend_error: 'keysend_error',
  keysend_success: 'keysend_success',
  invoice_success: 'invoice_success',
  assign_success: 'assign_success',
  lnauth_success: 'lnauth_success',
  user_connect: 'user_connect',
  budget_success: 'budget_success'
};

let socket: WebSocket | null = null;

export const createSocketInstance = (): WebSocket => {
  if (!socket || !socket.OPEN) {
    socket = new WebSocket(URL);
  }
  return socket;
};

export const getSocketInstance = (): WebSocket => {
  if (!socket) {
    throw new Error('Socket instance not created. Call createSocketInstance first.');
  }
  return socket;
};
