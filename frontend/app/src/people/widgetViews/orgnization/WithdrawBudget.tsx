import React, { useState } from 'react';
import { Wrap } from 'components/form/style';
import { useIsMobile } from 'hooks/uiHooks';
import { nonWidgetConfigs } from 'people/utils/Constants';
import { InvoiceForm, InvoiceInput, InvoiceLabel } from 'people/utils/style';
import { useStores } from 'store';
import moment from 'moment';
import styled from 'styled-components';
import { Button, Modal } from '../../../components/common';
import { colors } from '../../../config/colors';
import Invoice from '../summaries/wantedSummaries/Invoice';
import { ModalProps } from './interface';


const color = colors['light'];

const WithdrawModalTitle = styled.h3`
  font-size: 1.9rem;
  font-weight: bolder;
  margin-bottom: 20px;
`;

const WithdrawBudgetModal = (props: ModalProps) => {
    const [paymentRequest, setPaymentRequest] = useState('');
    const [lnInvoice, setLnInvoice] = useState('');

    const isMobile = useIsMobile();
    const { ui, main } = useStores();
    const { isOpen, close, uuid } = props;


    const config = nonWidgetConfigs['organizationusers'];

    const pollMinutes = 2;

    const withdrawBudget = async () => {
        const token = ui.meInfo?.websocketToken;
    };

    return (
        <Modal
            visible={isOpen}
            style={{
                height: '100%',
                flexDirection: 'column',
                padding: '10px'
            }}
            envStyle={{
                marginTop: isMobile ? 64 : 0,
                background: color.pureWhite,
                zIndex: 20,
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
                <WithdrawModalTitle>Withdraw</WithdrawModalTitle>
                {lnInvoice && ui.meInfo?.owner_pubkey && (
                    <>
                        <Invoice
                            startDate={new Date(moment().add(pollMinutes, 'minutes').format().toString())}
                            invoiceStatus={true}
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
                                Paste your invoice
                            </InvoiceLabel>
                            <InvoiceInput
                                type="text"
                                style={{
                                    width: '100%'
                                }}
                                value={paymentRequest}
                                onChange={(e: any) => setPaymentRequest(e.target.value)}
                            />
                        </InvoiceForm>
                        <Button
                            text={'Confirm'}
                            style={{
                                borderRadius: '10px',
                                marginTop: '12px',
                                background: '#9157F6'
                            }}
                            height={48}
                            width={'100%'}
                            onClick={withdrawBudget}
                        />
                    </>
                )}
            </Wrap>
        </Modal>
    )

};

export default WithdrawBudgetModal;