

import React from 'react'
import { Button } from '../../sphinxUI';
import QR from './QR';

export interface TorSaveQRProps {
    url: string
    goBack: Function
}

export default function TorSaveQR(props: TorSaveQRProps) {

    const { url, goBack } = props

    return <div
        style={{
            display: "flex",
            flexDirection: "column",
            justifyContent: "center",
            alignItems: "center",
            padding: "10px 20px",
            width: "100%",
        }}
    >
        <div style={{ height: 40 }} />

        <QR size={220} value={url} />

        <Button
            text={"Save on Sphinx"}
            height={60}
            style={{ marginTop: 30 }}
            width={"100%"}
            color={"primary"}
            onClick={() => {
                let el = document.createElement("a");
                el.href = url;
                el.click();
            }}
        />

        <Button
            text={"Dismiss"}
            height={60}
            style={{ marginTop: 20 }}
            width={"100%"}
            color={"action"}
            onClick={() => {
                goBack();
            }}
        />
    </div>

}