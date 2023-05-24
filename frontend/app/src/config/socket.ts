import { io } from 'socket.io-client';
import { getHost } from './host';

console.log("Hosts ===", getHost())
const URL = process.env.NODE_ENV === 'production' ?  `https://${getHost()}/socket.io` : `http://${getHost()}/socket.io`;

export const socket = io(URL);