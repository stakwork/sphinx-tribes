import React, { useState } from 'react';
import { Wrap } from 'components/form/style';
import { useIsMobile } from 'hooks/uiHooks';
import { InvoiceForm, InvoiceInput, InvoiceLabel } from 'people/utils/style';
import { useStores } from 'store';
import styled from 'styled-components';
import { Button, Modal } from '../../../components/common';
import { colors } from '../../../config/colors';
import successIcon from '../../../public/static/withdraw_success.svg';
import { ModalProps } from './interface';
import { Grey } from './style';

const color = colors['light'];

const WithdrawModalTitle = styled.h3`
  font-size: 1.9rem;
  font-weight: bolder;
  margin-bottom: 20px;
`;

const PaymentDetailsWrap = styled.div`
`;

const WithdrawText = styled.p`
  color: #3C3F41;
  font-size: 1.2rem;
`;

const WithdrawAmount = styled.h3`
  font-size: 1.5rem;
  color: #292C33;
`;

const DollarValue = styled.p`
  color: #8E969C;
  font-size: 0.9rem;
`;

const ActionButtonWrap = styled.div`
  display: flex;
`;

const SuccessImage = styled.img`
  width: 56px;
  height: 56px;
`;

const SuccessText = styled.p`
  color: #8E969C;
  font-size: 1.1rem;
`;

const WithdrawBudgetModal = (props: ModalProps) => {
    const [paymentRequest, setPaymentRequest] = useState('');
    const [paymentDetails, setPaymentDetails] = useState(true);
    const [paymentSuccess, setPaymentSuccess] = useState(true)

    const isMobile = useIsMobile();
    const { ui } = useStores();
    const { isOpen, close } = props;

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
                {paymentSuccess && (
                    <PaymentDetailsWrap>
                        <SuccessImage src={successIcon} />
                        <WithdrawAmount>{100250} <Grey>SATS</Grey></WithdrawAmount>
                        <SuccessText>Successfully Withdraw</SuccessText>
                    </PaymentDetailsWrap>
                )}
                {!paymentSuccess && paymentDetails && ui.meInfo?.owner_pubkey && (
                    <PaymentDetailsWrap>
                        <WithdrawText>
                            You are about to withdraw
                        </WithdrawText>
                        <WithdrawAmount>{100250} <Grey>SATS</Grey></WithdrawAmount>
                        <DollarValue>100.96 USD</DollarValue>
                        <ActionButtonWrap>
                            <Button
                                text={'Cancel'}
                                style={{
                                    borderRadius: '10px',
                                    background: '#FFFFFF',
                                    border: '1px solid #D0D5D8',
                                }}
                                height={48}
                                onClick={withdrawBudget}
                            />
                            <Button
                                text={'Withdraw'}
                                style={{
                                    borderRadius: '10px',
                                    background: '#9157F6'
                                }}
                                height={48}
                                onClick={withdrawBudget}
                            />
                        </ActionButtonWrap>
                    </PaymentDetailsWrap>
                )}
                {!paymentSuccess && !paymentDetails && ui.meInfo?.owner_pubkey && (
                    <>
                        <WithdrawModalTitle>Withdraw</WithdrawModalTitle>
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