import React from 'react';
import { OrgWrap } from 'components/form/style';
import { useIsMobile } from 'hooks/uiHooks';
import { nonWidgetConfigs } from 'people/utils/Constants';
import moment from 'moment';
import { PaymentHistory } from 'store/main';
import { useStores } from 'store';
import { Modal } from '../../../components/common';
import { colors } from '../../../config/colors';
import history from '../../../config/history';
import { ModalTitle, ViewBounty } from './style';
import { PaymentHistoryModalProps } from './interface';

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
      <OrgWrap style={{ width: '300px' }}>
        <ModalTitle>Payment history</ModalTitle>
        <table>
          <thead>
            <tr>
              <th>Sender</th>
              <th>Recipient</th>
              <th>Amount</th>
              <th>Date</th>
              <th />
            </tr>
          </thead>
          <tbody>
            {paymentsHistory.map((pay: PaymentHistory, i: number) => (
              <tr key={i}>
                <td className="ellipsis">{pay.sender_name}</td>
                <td className="ellipsis">{pay.receiver_name}</td>
                <td>{pay.amount} sats</td>
                <td>{moment(pay.created).format('DD/MM/YY')}</td>
                <td>
                  <ViewBounty onClick={() => viewBounty(pay.bounty_id)}>View bounty</ViewBounty>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </OrgWrap>
    </Modal>
  );
};

export default HistoryModal;
