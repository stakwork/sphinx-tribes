import React, { useState, useEffect } from 'react';
import { useIsMobile } from 'hooks/uiHooks';
import { nonWidgetConfigs } from 'people/utils/Constants';
import moment from 'moment';
import styled from 'styled-components';
import { PaymentHistory, OrgTransactionType } from 'store/main';
import { useStores } from 'store';
import { formatSat } from '../../../helpers';
import { Modal } from '../../../components/common';
import { colors } from '../../../config/colors';
import ArrowRight from '../../../public/static/arrow-right.svg';
import LinkIcon from '../../../public/static/link.svg';
import { PaymentHistoryModalProps } from './interface';
import UserInfo from './UserInfo';

const HistoryWrapper = styled.div`
  width: 61.125rem;
  height: 100%;
  overflow: hidden;
`;

const ModalHeaderWrapper = styled.div`
  display: flex;
  align-items: center;
  margin-top: 2.31rem;
  padding: 0rem 2.9rem;
  gap: 10.25rem;

  @media only screen and (max-width: 500px) {
    gap: 1rem;
    flex-direction: column;
    align-items: center;
    padding: 0rem 0.2rem;
    margin-top: 1rem;
  }
`;

const ModalTitle = styled.h2`
  color: #3c3f41;
  font-family: 'Barlow';
  font-size: 1.625rem;
  font-style: normal;
  font-weight: 800;
  line-height: 1.875rem;
  margin-bottom: 0;

  @media only screen and (max-width: 500px) {
    font-size: 1rem;
    font-weight: 600;
    line-height: 1rem;
  }
`;

const PaymentFilterWrapper = styled.div`
  display: flex;
  align-items: center;
  gap: 3rem;

  @media only screen and (max-width: 500px) {
    gap: 1rem;
    flex-direction: row;
    align-items: start;
  }
`;

const PaymentType = styled.div`
  display: flex;
  align-items: center;
  gap: 0.75rem;
`;

const Label = styled.label`
  margin-bottom: 0;
  color: #1e1f25;
  font-family: 'Barlow';
  font-size: 0.8125rem;
  font-style: normal;
  font-weight: 500;
  line-height: 2rem;
`;

const TableWrapper = styled.div`
  overflow: auto;
  width: 100%;
  height: 100%;
  margin-bottom: 20rem;
  margin-top: 2.56rem;

  @media only screen and (max-width: 500px) {
    margin-top: 1rem;
  }
`;

const Table = styled.table`
  border-collapse: collapse;
  margin-bottom: 6rem;
`;

const THead = styled.thead`
  position: sticky;
  top: 0;
  z-index: 100;
`;

const THeadRow = styled.tr`
  border-bottom: 1px solid #dde1e5;
  background: #fff;
  box-shadow: 0px 1px 4px 0px rgba(0, 0, 0, 0.1);
  padding: 1rem 2.9rem;

  @media only screen and (max-width: 500px) {
    padding: 0.3rem 0.5rem;
  }
`;

const TH = styled.th`
  color: #8e969c;
  font-family: 'Barlow';
  font-size: 0.625rem;
  font-style: normal;
  font-weight: 500;
  line-height: 1rem;
  letter-spacing: 0.0375rem;
  text-transform: uppercase;
  padding-top: 0.5rem;
  padding-bottom: 0.5rem;
`;

const ThLeft = styled(TH)`
  padding-left: 2.9rem;
  @media only screen and (max-width: 500px) {
    padding-left: 1.2rem;
  }
`;

const ThRight = styled(TH)`
  padding-right: 2.9rem;
  @media only screen and (max-width: 500px) {
    padding-right: 1.2rem;
  }
`;

const TR = styled.tr<{ type: OrgTransactionType }>`
  border-bottom: 1px solid
    ${(props: any) =>
      props.type === 'deposit'
        ? 'rgba(73, 201, 152, 0.20)'
        : props.type === 'withdraw'
        ? 'rgba(145, 87, 246, 0.15)'
        : 'rgba(0, 0, 0, 0.07)'};
  background-color: ${(props: any) =>
    props.type === 'deposit'
      ? '#F6FFFA'
      : props.type === 'withdraw'
      ? 'rgba(167, 108, 243, 0.05)'
      : ''};
`;

const TD = styled.td`
  color: #5f6368;
  font-family: 'Barlow';
  font-size: 0.8125rem;
  font-style: normal;
  font-weight: 400;
  line-height: 1rem;
  padding-top: 1.5rem;
  padding-bottom: 1.5rem;
  padding-right: 2.7rem;
  @media only screen and (max-width: 500px) {
    padding-left: 1.2rem;
    padding-right: 1.2rem;
  }
`;

const AmountSpan = styled.span`
  font-weight: 600;
`;

const TdLeft = styled(TD)<{ type: OrgTransactionType }>`
  color: ${(props: any) =>
    props.type === 'deposit' ? '#49C998' : props.type === 'withdraw' ? '#A76CF3' : '#3C3F41'};
  font-style: normal;
  font-family: ${(props: any) => (props.type === 'payment' ? 'Barlow' : 'Roboto')};
  font-weight: 600;
  line-height: 1rem;
  padding-left: 2.9rem;
  text-transform: capitalize;
  @media only screen and (max-width: 500px) {
    padding-left: 1.2rem;
  }
`;

const ArrowImage = styled.img`
  width: 1.25rem;
  height: 1.25rem;
`;

const LinkImage = styled.img`
  width: 1.25rem;
  height: 1.25rem;
`;

const ViewBountyContainer = styled.div`
  display: flex;
  align-items: center;
  gap: 0.375rem;
  cursor: pointer;
`;

const ViewBounty = styled.p`
  color: #5f6368;
  font-family: 'Barlow';
  font-size: 0.8125rem;
  font-style: normal;
  font-weight: 500;
  letter-spacing: 0.00813rem;
  margin-bottom: 0;
`;

const color = colors['light'];

const HistoryModal = (props: PaymentHistoryModalProps) => {
  const isMobile = useIsMobile();
  const { isOpen, close } = props;
  const [filter, setFilter] = useState({ payment: true, deposit: true, withdraw: true });
  const [currentPaymentsHistory, setCurrentPaymentHistory] = useState(props.paymentsHistory);

  const { ui } = useStores();

  const config = nonWidgetConfigs['organizationusers'];

  const viewBounty = async (bountyId: number) => {
    ui.setBountyPerson(ui.meInfo?.id);

    window.open(`/bounty/${bountyId}`, '_blank');
  };

  const handleFilter = (txn: OrgTransactionType) => {
    setFilter((value: any) => ({ ...value, [`${txn}`]: !value[txn] }));
  };

  useEffect(() => {
    const paymentsHistory = [...props.paymentsHistory];
    setCurrentPaymentHistory(
      paymentsHistory.filter((history: PaymentHistory) => filter[history.payment_type])
    );
  }, [filter, props.paymentsHistory]);

  return (
    <Modal
      visible={isOpen}
      style={{
        height: '100%',
        flexDirection: 'column',
        width: '100%',
        alignItems: `${isMobile ? '' : 'center'}`,
        justifyContent: `${isMobile ? '' : 'center'}`,
        overflowY: 'hidden'
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
        top: '1.6rem',
        right: `${isMobile ? '0rem' : '-1.25rem'}`,
        background: '#000',
        borderRadius: '50%'
      }}
    >
      <HistoryWrapper>
        <ModalHeaderWrapper>
          <ModalTitle>Payment History</ModalTitle>
          <PaymentFilterWrapper>
            <PaymentType>
              <input
                id="payment"
                type={'checkbox'}
                checked={filter.payment}
                onChange={() => handleFilter('payment')}
              />
              <Label htmlFor="payment">Payments</Label>
            </PaymentType>
            <PaymentType>
              <input
                id="deposit"
                type={'checkbox'}
                checked={filter.deposit}
                onChange={() => handleFilter('deposit')}
              />
              <Label htmlFor="deposit">Deposit</Label>
            </PaymentType>
            <PaymentType>
              <input
                id="withdraw"
                type={'checkbox'}
                checked={filter.withdraw}
                onChange={() => handleFilter('withdraw')}
              />
              <Label htmlFor="withdraw">Withdrawals</Label>
            </PaymentType>
          </PaymentFilterWrapper>
        </ModalHeaderWrapper>
        <TableWrapper>
          <Table>
            <THead>
              <THeadRow>
                <ThLeft>Type</ThLeft>
                <TH>Date</TH>
                <TH>Amount</TH>
                <TH>Sender</TH>
                <TH />
                <TH>Receiver</TH>
                <ThRight>Bounty</ThRight>
              </THeadRow>
            </THead>
            <tbody>
              {currentPaymentsHistory.map((pay: PaymentHistory, i: number) => (
                <TR type={pay.payment_type || 'payment'} key={i}>
                  <TdLeft type={pay.payment_type}>{pay.payment_type || 'Payment'}</TdLeft>
                  <TD>{moment(pay.created).format('MM/DD/YY')}</TD>
                  <TD>
                    <AmountSpan>{formatSat(pay.amount)}</AmountSpan> sats
                  </TD>
                  <TD>
                    <UserInfo
                      image={pay.sender_img}
                      pubkey={pay.sender_pubkey}
                      name={pay.sender_name}
                    />
                  </TD>
                  <TD>{pay.payment_type === 'payment' ? <ArrowImage src={ArrowRight} /> : null}</TD>
                  <TD>
                    {pay.payment_type === 'payment' ? (
                      <UserInfo
                        image={pay.receiver_img}
                        pubkey={pay.receiver_pubkey}
                        name={pay.receiver_name}
                      />
                    ) : null}
                  </TD>
                  <TD>
                    {pay.payment_type === 'payment' ? (
                      <ViewBountyContainer>
                        <ViewBounty onClick={() => viewBounty(pay.bounty_id)}>
                          View bounty
                        </ViewBounty>
                        <LinkImage src={LinkIcon} />
                      </ViewBountyContainer>
                    ) : null}
                  </TD>
                </TR>
              ))}
            </tbody>
          </Table>
        </TableWrapper>
      </HistoryWrapper>
    </Modal>
  );
};

export default HistoryModal;
