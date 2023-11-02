import React, { useState } from 'react';
import styled from 'styled-components';
import { useIsMobile } from 'hooks/uiHooks';
import { nonWidgetConfigs } from 'people/utils/Constants';
import { useStores } from 'store';
import moment from 'moment';
import { satToUsd } from 'helpers';
import { Modal } from '../../../components/common';
import { colors } from '../../../config/colors';
import Invoice from '../summaries/wantedSummaries/Invoice';
import { AddBudgetModalProps } from './interface';
const color = colors['light'];

const ModelWrapper = styled.div`
  width: 100%;
  padding: 3rem;
  display: flex;
  align-items: center;
  justify-content: center;
`;
const InvoiceFormHeading = styled.h2`
  color: #3c3f41;
  leading-trim: both;
  text-edge: cap;
  font-family: Barlow;
  font-size: 1.875rem;
  font-style: normal;
  font-weight: 800;
  line-height: normal;
  margin-bottom: 2.3rem;
`;

const InvoiceForm = styled.div`
  width: 100%;
`;

const InvoiceLabel = styled.label`
  color: #5f6368;
  leading-trim: both;
  text-edge: cap;
  font-family: Barlow;
  font-size: 1rem;
  font-style: normal;
  font-weight: 500;
  line-height: 2.1875rem;
  margin-bottom: 0.75rem;
`;

const InvoiceWrapper = styled.div`
  display: flex;
  flex-direction: column;
`;

const InvoiceInputWrapper = styled.div`
  display: flex;
  width: 100%;
  padding: 0.5rem 1rem;
  justify-content: center;
  align-items: center;
  gap: 0.75rem;
  border-radius: 0.375rem;
  border: 2px solid #86d9b9;
`;

const CurrencyUnit = styled.p`
  color: #b0b7bc;
  text-align: right;
  leading-trim: both;
  text-edge: cap;
  font-family: Barlow;
  font-size: 0.9375rem;
  font-style: normal;
  font-weight: 400;
  line-height: 2.1875rem;
  margin-bottom: 0;
`;

const Input = styled.input`
  border: none;
  outline: none;
  color: #3c3f41;
  leading-trim: both;
  text-edge: cap;
  font-family: Barlow;
  font-size: 1.25rem;
  font-style: normal;
  font-weight: 400;
  line-height: 2.1875rem;
  caret-color: #49c998;
  width: 100%;
  ::placeholder {
    color: #b0b7bc;
  }
`;

const UsdValue = styled.p`
  color: #909baa;
  leading-trim: both;
  text-edge: cap;
  font-family: Barlow;
  font-size: 0.8125rem;
  font-style: normal;
  font-weight: 500;
  line-height: normal;
  margin-top: 0.6rem;
  margin-bottom: 0;
`;

const Button = styled.button`
  width: 100%;
  padding: 1rem;
  border-radius: 0.375rem;
  margin-top: 1.25rem;
  font-family: Barlow;
  font-size: 0.9375rem;
  font-style: normal;
  font-weight: 500;
  letter-spacing: 0.00938rem;
  background: #49c998;
  box-shadow: 0px 2px 10px 0px rgba(73, 201, 152, 0.5);
  border: none;
  color: #fff;
  &:disabled {
    border: 1px solid rgba(0, 0, 0, 0.07);
    background: rgba(0, 0, 0, 0.04);
    color: rgba(142, 150, 156, 0.85);
    cursor: not-allowed;
    box-shadow: none;
  }
`;
const AddBudgetModal = (props: AddBudgetModalProps) => {
  const [amount, setAmount] = useState('');
  const [lnInvoice, setLnInvoice] = useState('');

  const isMobile = useIsMobile();
  const { ui, main } = useStores();
  const { isOpen, close, invoiceStatus, uuid, startPolling } = props;

  const config = nonWidgetConfigs['organizationusers'];

  const pollMinutes = 2;

  const generateInvoice = async () => {
    if (uuid) {
      const data = await main.getBudgetInvoice({
        amount: Number(amount),
        sender_pubkey: ui.meInfo?.owner_pubkey ?? '',
        org_uuid: uuid,
        payment_type: 'deposit'
      });

      const paymentRequest = data.response.invoice;

      if (paymentRequest) {
        setLnInvoice(paymentRequest);
        startPolling(paymentRequest);

        main.setBudgetInvoice(paymentRequest);
      }
    }
  };

  const handleInputAmountChange = (e: any) => {
    const inputValue = e.target.value;
    const numericValue = inputValue.replace(/[^0-9]/g, '');
    setAmount(numericValue);
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
      <ModelWrapper>
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
              <InvoiceFormHeading>Deposit</InvoiceFormHeading>
              <InvoiceWrapper>
                <InvoiceLabel>Amount (in sats)</InvoiceLabel>
                <InvoiceInputWrapper>
                  <Input
                    placeholder="0"
                    type="text"
                    value={amount}
                    onChange={handleInputAmountChange}
                  />
                  <CurrencyUnit>sats</CurrencyUnit>
                </InvoiceInputWrapper>
                <UsdValue>{satToUsd(Number(amount))} USD</UsdValue>
              </InvoiceWrapper>
              <Button disabled={!amount} onClick={generateInvoice}>
                Generate Invoice
              </Button>
            </InvoiceForm>
          </>
        )}
      </ModelWrapper>
    </Modal>
  );
};

export default AddBudgetModal;
