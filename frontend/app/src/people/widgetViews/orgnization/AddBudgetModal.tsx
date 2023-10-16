import React, { useState } from 'react';
import { Wrap } from 'components/form/style';
import { useIsMobile } from 'hooks/uiHooks';
import { nonWidgetConfigs } from 'people/utils/Constants';
import { InvoiceForm, InvoiceInput, InvoiceLabel } from 'people/utils/style';
import { useStores } from 'store';
import moment from 'moment';
import { Button, Modal } from '../../../components/common';
import { colors } from '../../../config/colors';
import Invoice from '../summaries/wantedSummaries/Invoice';
import { AddBudgetModalProps } from './interface';
import { ModalTitle } from './style';

const color = colors['light'];

const AddBudgetModal = (props: AddBudgetModalProps) => {
    const [amount, setAmount] = useState(1);
    const [lnInvoice, setLnInvoice] = useState('');

    const isMobile = useIsMobile();
    const { ui, main } = useStores();
    const { isOpen, close, invoiceStatus, uuid } = props;


    const config = nonWidgetConfigs['organizationusers'];

    const pollMinutes = 2;

    const generateInvoice = async () => {
        const token = ui.meInfo?.websocketToken;
        if (token && uuid) {
            const data = await main.getBudgetInvoice({
                amount: amount,
                sender_pubkey: ui.meInfo?.owner_pubkey ?? '',
                org_uuid: uuid,
                websocket_token: token,
                payment_type: 'deposit'
            });

            setLnInvoice(data.response.invoice);
        }
    };

    return (
        <Modal
            visible={isOpen}
            style={{
                height: '100%',
                flexDirection: 'column'
            }}
            envStyle={{
                marginTop: isMobile ? 64 : 0,
                background: color.pureWhite,
                zIndex: 20,
                ...(config?.modalStyle ?? {}),
                maxHeight: '100%',
                borderRadius: '10px'
            }}
            overlayClick={close}
            bigCloseImage={close}
            bigCloseImageStyle={{
                top: '-18px',
                right: '-18px',
                background: '#000',
                borderRadius: '50%'
            }}
        >
            <Wrap newDesign={true}>
                <ModalTitle>Add budget</ModalTitle>
                {lnInvoice && ui.meInfo?.owner_pubkey && (
                    <>
                        <Invoice
                            startDate={new Date(moment().add(pollMinutes, 'minutes').format().toString())}
                            invoiceStatus={invoiceStatus}
                            lnInvoice={lnInvoice}
                            invoiceTime={pollMinutes}
                        />
                    </>
                )}
                {!lnInvoice && ui.meInfo?.owner_pubkey && (
                    <>
                        <InvoiceForm>
                            <InvoiceLabel
                                style={{
                                    display: 'block'
                                }}
                            >
                                Amount (in sats)
                            </InvoiceLabel>
                            <InvoiceInput
                                type="number"
                                style={{
                                    width: '100%'
                                }}
                                value={amount}
                                onChange={(e: any) => setAmount(Number(e.target.value))}
                            />
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
            </Wrap>
        </Modal>
    )

};

export default AddBudgetModal;