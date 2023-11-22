import React from 'react';
import styled from 'styled-components';
import { InvoiceState } from './interface';

const Container = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  padding-top: 3rem;
  padding-bottom: 3rem;
`;

const ImgWrapper = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  width: 7.25rem;
  height: 7.25rem;
`;

const Text = styled.p`
  margin-bottom: 2rem;
  margin-top: 2rem;
  color: #ed7474;
  text-align: center;
  font-family: 'Barlow';
  font-size: 1.375rem;
  font-style: normal;
  font-weight: 500;
  line-height: normal;
  letter-spacing: 0.01375rem;
`;

const Button = styled.button`
  width: 9.25rem;
  height: 2.5rem;
  padding: 0.5rem 1rem;
  border-radius: 0.375rem;
  border: 1px solid #d0d5d8;
  background: #fff;
  box-shadow: 0px 1px 2px 0px rgba(0, 0, 0, 0.06);
  color: #5f6368;
  font-family: 'Barlow';
  font-size: 0.875rem;
  font-style: normal;
  font-weight: 500;
  line-height: 0rem;
  letter-spacing: 0.00875rem;
`;

const ExpiredInvoice = (props: {
  setLnInvoice: (ln: string) => void;
  setInvoiceState: (state: InvoiceState) => void;
}) => {
  const handleTryAgain = () => {
    props.setInvoiceState(null);
    props.setLnInvoice('');
  };
  return (
    <Container>
      <ImgWrapper>
        <img src="/static/expired_invoice.svg" alt="expired" />
      </ImgWrapper>
      <Text>Invoice Expired</Text>
      <Button onClick={handleTryAgain}>Try Again</Button>
    </Container>
  );
};

export default ExpiredInvoice;
