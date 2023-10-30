import React from 'react';
import { useIsMobile } from 'hooks/uiHooks';
import { nonWidgetConfigs } from 'people/utils/Constants';
import moment from 'moment';
import styled from 'styled-components';
import { PaymentHistory } from 'store/main';
import { useStores } from 'store';
import { Modal } from '../../../components/common';
import { colors } from '../../../config/colors';
import history from '../../../config/history';
import ArrowRight from '../../../public/static/arrow-right.svg';
import LinkIcon from '../../../public/static/link.svg';
import { PaymentHistoryModalProps } from './interface';
import UserInfo from './UserInfo';

type OrgTransactionType = 'Deposit' | 'Payment' | 'Withdraw';

const HistoryWrapper = styled.div`
  width: 61.125rem;
`;

const ModalHeaderWrapper = styled.div`
  display: flex;
  align-items: center;
  margin-top: 2.3rem;
  padding: 0rem 2.9rem;
`;

const ModalTitle = styled.h2`
  color: #3c3f41;
  leading-trim: both;
  text-edge: cap;
  font-family: Barlow;
  font-size: 1.625rem;
  font-style: normal;
  font-weight: 800;
  line-height: 1.875rem;
`;

const Table = styled.table`
  margin-top: 2rem;
  border-collapse: collapse;
  margin-bottom: 2rem;
  overflow-x: auto;
`;

const THeadRow = styled.tr`
  border-bottom: 1px solid #dde1e5;
  background: #fff;
  box-shadow: 0px 1px 4px 0px rgba(0, 0, 0, 0.1);
  padding: 1rem 2.9rem;
`;

const TH = styled.th`
  color: #8e969c;
  leading-trim: both;
  text-edge: cap;
  font-family: Barlow;
  font-size: 0.8rem;
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
`;

const ThRight = styled(TH)`
  padding-right: 2.9rem;
`;

const TR = styled.tr<{ type: OrgTransactionType }>`
  border-bottom: 1px solid
    ${(props: any) =>
      props.type === 'Deposit'
        ? 'rgba(73, 201, 152, 0.20)'
        : props.type === 'Withdraw'
        ? 'rgba(145, 87, 246, 0.15)'
        : 'rgba(0, 0, 0, 0.07)'};
  background-color: ${(props: any) =>
    props.type === 'Deposit'
      ? '#F6FFFA'
      : props.type === 'Withdraw'
      ? 'rgba(167, 108, 243, 0.05)'
      : ''};
`;

const TD = styled.td`
  color: #5f6368;
  font-family: Barlow;
  font-size: 1rem;
  font-style: normal;
  font-weight: 400;
  line-height: 1rem;
  padding-top: 1.5rem;
  padding-bottom: 1.5rem;
`;

const AmountSpan = styled.span`
  font-weight: 600;
`;

const TdLeft = styled(TD)<{ type: OrgTransactionType }>`
  color: ${(props: any) =>
    props.type === 'Deposit' ? '#49C998' : props.type === 'Withdraw' ? '#A76CF3' : '#3C3F41'};
  font-size: 1rem;
  font-style: normal;
  font-weight: 600;
  line-height: 1rem;
  padding-left: 2.9rem;
`;

const ArrowImage = styled.img`
  width: 1.25rem;
  height: 1.25rem;
  margin-right: 1.1rem;
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
  leading-trim: both;
  text-edge: cap;
  font-family: Barlow;
  font-size: 0.8125rem;
  font-style: normal;
  font-weight: 500;
  letter-spacing: 0.00813rem;
  margin-bottom: 0;
`;

const color = colors['light'];

const HistoryModal = (props: PaymentHistoryModalProps) => {
  const isMobile = useIsMobile();
  const { isOpen, close, url, paymentsHistory } = props;

  const { ui } = useStores();

  const config = nonWidgetConfigs['organizationusers'];

  const viewBounty = async (bountyId: number) => {
    ui.setBountyPerson(ui.meInfo?.id);

    history.push({
      pathname: `${url}/${bountyId}/${0}`
    });
  };

  console.log(paymentsHistory);

  return (
    <Modal
      visible={isOpen}
      style={{
        height: '100%',
        flexDirection: 'column',
        width: '100%'
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
      <HistoryWrapper>
        <ModalHeaderWrapper>
          <ModalTitle>Payment history</ModalTitle>
        </ModalHeaderWrapper>
        <Table>
          <thead>
            <THeadRow>
              <ThLeft>Type</ThLeft>
              <TH>Date</TH>
              <TH>Amount</TH>
              <TH>Sender</TH>
              <TH />
              <TH>Receiver</TH>
              <ThRight>Bounty</ThRight>
            </THeadRow>
          </thead>
          <tbody>
            {paymentsHistory.map((pay: PaymentHistory, i: number) => (
              <TR type={i === 0 ? 'Payment' : i === 1 ? 'Deposit' : 'Withdraw'} key={i}>
                <TdLeft type={i === 0 ? 'Payment' : i === 1 ? 'Deposit' : 'Withdraw'}>
                  {i === 0 ? 'Payment' : i === 1 ? 'Deposit' : 'Withdraw'}
                </TdLeft>
                <TD>{moment(pay.created).format('DD/MM/YY')}</TD>
                <TD>
                  <AmountSpan>{pay.amount}</AmountSpan> sats
                </TD>
                <TD>
                  <UserInfo
                    image="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQR9dAUM-b34F_a6DMw8D6fQ_Y0LUIAVzvfCw&usqp=CAU"
                    pubkey="030841d1519f19c68e80efc5ef5af3460ca4bfa17486fda9baca878b9ef255358"
                    name={pay.sender_name}
                  />
                </TD>
                <TD>
                  <ArrowImage src={ArrowRight} />
                </TD>
                <TD>
                  <UserInfo
                    image="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQR9dAUM-b34F_a6DMw8D6fQ_Y0LUIAVzvfCw&usqp=CAU"
                    pubkey="030841d1519f19c68e80efc5ef5af3460ca4bfa17486fda9baca878b9ef255358"
                    name={pay.receiver_name}
                  />
                </TD>
                <TD>
                  <ViewBountyContainer>
                    <ViewBounty onClick={() => viewBounty(pay.bounty_id)}>View bounty</ViewBounty>
                    <LinkImage src={LinkIcon} />
                  </ViewBountyContainer>
                </TD>
              </TR>
            ))}
          </tbody>
        </Table>
      </HistoryWrapper>
    </Modal>
  );
};

export default HistoryModal;
