import { useState, useEffect } from 'react';
import { getHost } from './host';

const URL =
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

const useSocket = () => {
    const [data, setData] = useState<WebSocket>();

    useEffect(() => {
        const socket: WebSocket = new WebSocket(URL);
        setData(socket);
    }, []);

    return [data];
};

export default useSocket;
