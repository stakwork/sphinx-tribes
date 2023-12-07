import React, { useState, useEffect } from 'react';
import styled from 'styled-components';
import { useIsMobile } from 'hooks/uiHooks';
import { nonWidgetConfigs } from 'people/utils/Constants';
import { useStores } from 'store';
import moment from 'moment';
import { satToUsd } from 'helpers';
import { EuiGlobalToastList, EuiLoadingSpinner } from '@elastic/eui';
import { Modal } from '../../../components/common';
import { colors } from '../../../config/colors';
import Invoice from '../summaries/wantedSummaries/Invoice';
import { AddBudgetModalProps, InvoiceState, Toast } from './interface';
import ExpiredInvoice from './ExpiredInvoice';
import PaidInvoice from './PaidInvoice';
import { BudgetButton } from './style';

const color = colors['light'];

const ModelWrapper = styled.div`
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
`;
const InvoiceFormHeading = styled.h2`
  color: #3c3f41;
  font-family: 'Barlow';
  font-size: 1.875rem;
  font-style: normal;
  font-weight: 800;
  line-height: normal;
  margin-bottom: 2.3rem;
`;

const InvoiceForm = styled.div`
  width: 16rem;
  margin: 3rem;
`;

const InvoiceLabel = styled.label`
  color: #5f6368;
  font-family: 'Barlow';
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
  font-family: 'Barlow';
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
  font-family: 'Barlow';
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
  font-family: 'Barlow';
  font-size: 0.8125rem;
  font-style: normal;
  font-weight: 500;
  line-height: normal;
  margin-top: 0.6rem;
  margin-bottom: 0;
`;

const InvoiceQrWrapper = styled.div`
  width: 13.5rem;
  margin-left: 4.75rem;
  margin-right: 4.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
`;

const AddBudgetModal = (props: AddBudgetModalProps) => {
  const [amount, setAmount] = useState('');
  const [lnInvoice, setLnInvoice] = useState('');
  const [invoiceState, setInvoiceState] = useState<InvoiceState>(null);

  const isMobile = useIsMobile();
  const { ui, main } = useStores();
  const { isOpen, close, invoiceStatus, uuid, startPolling } = props;
  const [isLoading, setIsLoading] = useState(false);
  const [toasts, setToasts] = useState<Toast[]>([]);

  const config = nonWidgetConfigs['organizationusers'];

  const pollMinutes = 2;

  function addSuccessToast() {
    setToasts([
      {
        id: '1',
        title: 'Create Invoice',
        color: 'success',
        text: 'Invoice Created Successfully'
      }
    ]);
  }

  function addErrorToast(text: string) {
    setToasts([
      {
        id: '2',
        title: 'Create Invoice',
        color: 'danger',
        text
      }
    ]);
  }

  function removeToast() {
    setToasts([]);
  }

  const generateInvoice = async () => {
    if (uuid) {
      try {
        setIsLoading(true);
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
          setInvoiceState('PENDING');
          main.setBudgetInvoice(paymentRequest);
          addSuccessToast();
          return;
        }
        addErrorToast('Error occured while creating invoice');
        setIsLoading(false);
      } catch (error) {
        setIsLoading(false);
        addErrorToast('Error occured while creating invoice');
      }
    }
  };

  const handleInputAmountChange = (e: any) => {
    const inputValue = e.target.value;
    const numericValue = inputValue.replace(/[^0-9]/g, '');
    setAmount(numericValue);
  };

  useEffect(() => {
    if (invoiceStatus) {
      setInvoiceState('PAID');
      props.setInvoiceStatus(false);
    }
  }, [invoiceStatus]);

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
      <ModelWrapper>
        {lnInvoice && invoiceState === 'PENDING' && ui.meInfo?.owner_pubkey && (
          <InvoiceQrWrapper>
            <Invoice
              startDate={new Date(moment().add(pollMinutes, 'minutes').format().toString())}
              invoiceStatus={invoiceStatus}
              lnInvoice={lnInvoice}
              invoiceTime={pollMinutes}
              setInvoiceState={setInvoiceState}
              qrSize={216}
            />
          </InvoiceQrWrapper>
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
              <BudgetButton disabled={!Number(amount) || isLoading} onClick={generateInvoice}>
                {isLoading ? <EuiLoadingSpinner size="m" /> : 'Generate Invoice'}
              </BudgetButton>
            </InvoiceForm>
          </>
        )}
        {invoiceState === 'EXPIRED' && (
          <ExpiredInvoice setInvoiceState={setInvoiceState} setLnInvoice={setLnInvoice} />
        )}
        {invoiceState === 'PAID' && <PaidInvoice amount={Number(amount)} />}
      </ModelWrapper>
      <EuiGlobalToastList toasts={toasts} dismissToast={removeToast} toastLifeTimeMs={3000} />
    </Modal>
  );
};

export default AddBudgetModal;
