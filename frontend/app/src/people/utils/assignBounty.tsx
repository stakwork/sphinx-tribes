import React, { useState } from 'react';
import styled from 'styled-components';
import { Button, Modal } from '../../components/common';
import { colors } from '../../config/colors';
import { ConnectCardProps } from 'people/interfaces';
import { useStores } from 'store';
import { EuiGlobalToastList } from '@elastic/eui';
import Invoice from '../widgetViews/summaries/wantedSummaries/invoice';
import moment from 'moment';

export default function AssignBounty(props: ConnectCardProps) {
    const color = colors['light'];
    const { visible, person, created } = props;
    const { main, ui } = useStores();

    const [bountyHours, setBountyHours] = useState(1);
    const [bountyInvoice, setBountyInvoice] = useState<string>("");
    const [pollCount, setPollCount] = useState(0);
    const [invoiceData, setInvoiceData] = useState<{ invoiceStatus: boolean; bountyPaid: boolean }>({
        invoiceStatus: false,
        bountyPaid: false
    });

    const [toasts, setToasts]: any = useState([]);

    function addToast() {
        setToasts([
            {
                id: '1',
                title: 'Bounty has been assigned'
            }
        ]);
    }

    function removeToast() {
        setToasts([]);
    }

    const generateInvoice = async () => {
        const data = await main.getLnInvoice({
            amount: 200 * bountyHours, memo: "",
            owner_pubkey: person?.owner_pubkey ?? '',
            user_pubkey: ui.meInfo?.owner_pubkey ?? '',
            created: created ? created.toString() : '',
            type: 'ASSIGN'
        })

        setBountyInvoice(data.response.invoice);

        await pollLnInvoice(pollCount);
    }

    async function pollLnInvoice(count: number) {
        if (main.lnInvoice) {
            const data = await main.getLnInvoiceStatus(main.lnInvoice);

            setInvoiceData(data);

            setPollCount(count);

            const pollTimeout = setTimeout(() => {
                pollLnInvoice(count + 1);
                setPollCount(count + 1);
            }, 2000);

            if (data.invoiceStatus) {
                clearTimeout(pollTimeout);
                setPollCount(0);
                main.setLnInvoice('');

                // display a toast
                addToast();
            }

            if (count >= 29) {
                clearTimeout(pollTimeout);
                setPollCount(0);
                main.setLnInvoice('');
            }
        }
    }

    return (
        <div onClick={(e) => e.stopPropagation()}>
            <Modal
                style={props.modalStyle}
                overlayClick={() => {
                    props.dismiss();
                }}
                visible={visible}
            >
                <div style={{ textAlign: 'center', paddingTop: 59, width: 310 }}>
                    <div
                        style={{ textAlign: 'center', width: '100%', overflow: 'hidden', padding: '0 50px' }}
                    >
                        <N color={color}>Asign bounty to your self</N>
                        <B>Each hour cost 200 sats</B>


                        {
                            bountyInvoice && ui.meInfo?.owner_pubkey && (
                                <Invoice
                                    startDate={new Date(moment().add(1, 'minutes').format().toString())}
                                    count={pollCount}
                                    dataStatus={invoiceData.invoiceStatus}
                                />
                            )
                        }

                        {!bountyInvoice && ui.meInfo?.owner_pubkey && (
                            <>
                                <InvoiceForm>
                                    <InvoiceLabel>Number Of Hours</InvoiceLabel>
                                    <InvoiceInput type="number" value={bountyHours} onChange={(e) => setBountyHours(Number(e.target.value))} />
                                </InvoiceForm>
                                <Button
                                    text={'Generate Invoice'}
                                    color={'primary'}
                                    style={{ paddingLeft: 25, margin: '12px 0 10px' }}
                                    img={'sphinx_white.png'}
                                    imgSize={27}
                                    height={48}
                                    width={'100%'}
                                    onClick={generateInvoice}
                                />
                            </>
                        )}

                    </div>
                </div>
                <ModalBottomText color={color}>
                    <img src="/static/scan_qr.svg" alt="scan" />
                    <div className="bottomText">Pay the invoice to assign to your self</div>
                </ModalBottomText>
            </Modal>
            <EuiGlobalToastList toasts={toasts} dismissToast={removeToast} toastLifeTimeMs={1000} />
        </div>
    );
}

interface styledProps {
    color?: any;
}

const B = styled.small`
  font-weight: bold;
  display: block;
  margin-bottom: 10px;
`;
const N = styled.div<styledProps>`
  font-family: Barlow;
  font-style: normal;
  font-weight: 500;
  font-size: 17px;
  line-height: 26px;
  text-align: center;
  margin-bottom: 10px;
  color: ${(p) => p?.color && p?.color.grayish.G100};
`;
const ModalBottomText = styled.div<styledProps>`
  position: absolute;
  bottom: -36px;
  width: 310;
  background-color: transparent;
  display: flex;
  justify-content: center;
  .bottomText {
    margin-left: 12px;
    color: ${(p) => p?.color && p?.color.pureWhite};
  }
`;
const InvoiceForm = styled.div`
  margin: 10px 0px;
  text-align: left;
`;
const InvoiceLabel = styled.label`
  font-size: 0.9rem;
  font-weight: bold;
`;
const InvoiceInput = styled.input`
  padding: 10px 20px;
  border-radius: 10px;
  border: 0.5px solid black;
`
