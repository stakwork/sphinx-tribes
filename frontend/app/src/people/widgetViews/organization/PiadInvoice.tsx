import React from 'react';
import styled from 'styled-components';

const Wrapper = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
`;

const ImgWrapper = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  width: 3.5rem;
  height: 3.5rem;
`;

const AmountText = styled.h2`
  color: #292c33;
  text-align: center;
  leading-trim: both;
  text-edge: cap;
  font-family: Barlow;
  font-size: 1.875rem;
  font-style: normal;
  font-weight: 500;
  line-height: normal;
  letter-spacing: 0.01875rem;
  text-transform: uppercase;
  margin-top: 2.5rem;
  margin-bottom: 1.25rem;
`;

const AmountUnit = styled.span`
  color: #8e969c;
`;

const Text = styled.p`
  color: #8e969c;
  text-align: center;
  leading-trim: both;
  text-edge: cap;
  font-family: Barlow;
  font-size: 1.0625rem;
  font-style: normal;
  font-weight: 500;
  line-height: normal;
  letter-spacing: 0.01063rem;
  margin-bottom: 0;
`;

const PaidInvoice = (props: { amount: number }) => {
  return (
    <>
      <Wrapper>
        <ImgWrapper>
          <img src="/static/success.svg" alt="success" />
        </ImgWrapper>
        <AmountText>
          {props.amount.toLocaleString()} <AmountUnit>SATS</AmountUnit>
        </AmountText>
        <Text>Successfully Deposited</Text>
      </Wrapper>
    </>
  );
};

export default PaidInvoice;
