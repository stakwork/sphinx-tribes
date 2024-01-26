import React, { useState } from 'react';
import { useIsMobile } from 'hooks/uiHooks';
import { InvoiceForm, InvoiceInput, InvoiceLabel } from 'people/utils/style';
import { useStores } from 'store';
import styled from 'styled-components';
import { BudgetWithdrawSuccess } from 'store/main';
import { satToUsd } from 'helpers';
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import LighningDecoder from 'light-bolt11-decoder';
import { Button, Modal } from '../../../components/common';
import { colors } from '../../../config/colors';
import successIcon from '../../../public/static/withdraw_success.svg';
import errorIcon from '../../../public/static/error.svg';
import { WithdrawModalProps } from './interface';
import { BudgetButton, Grey } from './style';

const color = colors['light'];

const WithdrawModalTitle = styled.h3`
  font-size: 1.9rem;
  font-weight: bolder;
  margin-bottom: 20px;
`;

const PaymentDetailsWrap = styled.div`
  padding: 10px 0px;
  display: flex;
  flex-direction: column;
  align-items: center;
`;

const WithdrawText = styled.p`
  color: #3c3f41;
  font-size: 1.05rem;
  margin-bottom: 30px;
`;

const WithdrawAmount = styled.h3`
  font-size: 1.8rem;
  color: #292c33;
  padding: 0px;
  margin: 0px;
`;

const DollarValue = styled.p`
  color: #8e969c;
  font-size: 0.9rem;
  padding: 0px;
  margin: 0px;
`;

const ActionButtonWrap = styled.div`
  display: flex;
  margin-top: 40px;
`;

const SuccessImage = styled.img`
  width: 56px;
  height: 56px;
  margin-bottom: 30px;
`;

const SuccessText = styled.p`
  color: #8e969c;
  font-size: 0.95rem;
  font-weight: 500;
  margin-top: 10px;
`;

const ErrorText = styled(SuccessText)`
  color: #ed7474;
  text-align: center;
`;

const Wrapper = styled.div`
  width: 100%;
  display: flex;
  flex-direction: column;
  padding: 40px 50px;
`;

const WithdrawBudgetModal = (props: WithdrawModalProps) => {
  const [paymentRequest, setPaymentRequest] = useState('');
  const [amountInSats, setAmountInSats] = useState(0);
  const [paymentSettled, setPaymentSettled] = useState(false);
  const [paymentError, setPaymentError] = useState('');

  const isMobile = useIsMobile();
  const { ui, main } = useStores();
  const { isOpen, close, uuid, getOrganizationBudget } = props;

  const withdrawBudget = async () => {
    const token = ui.meInfo?.websocketToken;
    const body = {
      org_uuid: uuid ?? '',
      payment_request: paymentRequest,
      websocket_token: token
    };

    const isBudgetSuccess = (object: any): object is BudgetWithdrawSuccess => 'response' in object;

    const response = await main.withdrawBountyBudget(body);
    if (isBudgetSuccess(response)) {
      await getOrganizationBudget();
      setPaymentSettled(true);
    } else {
      setPaymentError(response.error);
    }
  };

  const getInvoiceDetails = async (paymentRequest: string) => {
    try {
      const decoded = LighningDecoder.decode(paymentRequest);
      const sats = decoded.sections[2].value / 1000;
      setAmountInSats(sats);
    } catch (e) {
      console.log(`Cannot decode lightning invoice: ${e}`);
    }
  };

  const displayWuthdraw = !amountInSats && ui.meInfo?.owner_pubkey;
  const displayInvoiceSats =
    amountInSats > 0 && !paymentSettled && !paymentError && ui.meInfo?.owner_pubkey;

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
        maxHeight: '100%',
        borderRadius: '10px',
        minWidth: '22rem',
        minHeight: '20rem'
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
      <Wrapper>
        {paymentSettled && (
          <PaymentDetailsWrap>
            <SuccessImage src={successIcon} />
            <WithdrawAmount>
              {amountInSats.toLocaleString()} <Grey>SATS</Grey>
            </WithdrawAmount>
            <SuccessText>Successfully Withdraw</SuccessText>
          </PaymentDetailsWrap>
        )}
        {paymentError && (
          <PaymentDetailsWrap>
            <SuccessImage src={errorIcon} />
            <WithdrawAmount>
              {amountInSats.toLocaleString()} <Grey>SATS</Grey>
            </WithdrawAmount>
            <ErrorText>{paymentError}</ErrorText>
          </PaymentDetailsWrap>
        )}
        {displayInvoiceSats && (
          <PaymentDetailsWrap>
            <WithdrawText>You are about to withdraw</WithdrawText>
            <WithdrawAmount>
              {amountInSats.toLocaleString()} <Grey>SATS</Grey>
            </WithdrawAmount>
            <DollarValue>{satToUsd(Number(amountInSats))} USD</DollarValue>
            <ActionButtonWrap>
              <Button
                text={'Cancel'}
                style={{
                  borderRadius: '6px',
                  background: '#FFFFFF',
                  border: '1px solid #D0D5D8',
                  color: '#5F6368',
                  marginRight: '8px'
                }}
                height={42}
                onClick={close}
              />
              <Button
                text={'Withdraw'}
                style={{
                  borderRadius: '6px',
                  background: '#9157F6',
                  marginLeft: '8px'
                }}
                height={42}
                onClick={withdrawBudget}
              />
            </ActionButtonWrap>
          </PaymentDetailsWrap>
        )}
        {displayWuthdraw && (
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
            <BudgetButton
              disabled={!paymentRequest ? true : false}
              style={{
                borderRadius: '8px',
                marginTop: '12px',
                color: paymentRequest ? '#FFF' : 'rgba(142, 150, 156, 0.85)',
                background: paymentRequest ? '#9157F6' : 'rgba(0, 0, 0, 0.04)'
              }}
              onClick={() => getInvoiceDetails(paymentRequest)}
            >
              Confirm
            </BudgetButton>
          </>
        )}
      </Wrapper>
    </Modal>
  );
};

export default WithdrawBudgetModal;
