import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { ConnectCardProps } from 'people/interfaces';
import { useStores } from 'store';
import { EuiGlobalToastList } from '@elastic/eui';
import moment from 'moment';
import { SOCKET_MSG, createSocketInstance } from 'config/socket';
import Invoice from '../widgetViews/summaries/wantedSummaries/Invoice';
import { colors } from '../../config/colors';
import { Button, Modal } from '../../components/common';

interface styledProps {
  color?: any;
}

const B = styled.small`
  font-weight: bold;
  display: block;
  margin-bottom: 10px;
`;
const N = styled.div<styledProps>`
  font-family: Barlow;
  font-style: normal;
  font-weight: 500;
  font-size: 17px;
  line-height: 26px;
  text-align: center;
  margin-bottom: 10px;
  color: ${(p: any) => p?.color && p?.color.grayish.G100};
`;
const ModalBottomText = styled.div<styledProps>`
  position: absolute;
  bottom: -36px;
  width: 310;
  background-color: transparent;
  display: flex;
  justify-content: center;
  .bottomText {
    margin-left: 12px;
    color: ${(p: any) => p?.color && p?.color.pureWhite};
  }
`;
const InvoiceForm = styled.div`
  margin: 10px 0px;
  text-align: left;
`;
const InvoiceLabel = styled.label`
  font-size: 0.9rem;
  font-weight: bold;
`;
const InvoiceInput = styled.input`
  padding: 10px 20px;
  border-radius: 10px;
  border: 0.5px solid black;
`;
export default function AssignBounty(props: ConnectCardProps) {
  const color = colors['light'];
  const { person, created, visible } = props;
  const { main, ui } = useStores();

  const [bountyHours, setBountyHours] = useState(1);
  const [lnInvoice, setLnInvoice] = useState('');
  const [invoiceStatus, setInvoiceStatus] = useState(false);

  const pollMinutes = 2;

  const [toasts, setToasts]: any = useState([]);

  const addToast = () => {
    setToasts([
      {
        id: '1',
        title: 'Bounty has been assigned'
      }
    ]);
  };
  const removeToast = () => {
    setToasts([]);
  };

  const generateInvoice = async () => {
    if (created && ui.meInfo?.websocketToken) {
      const data = await main.getLnInvoice({
        amount: 200 * bountyHours,
        memo: '',
        owner_pubkey: person?.owner_pubkey ?? '',
        user_pubkey: ui.meInfo?.owner_pubkey ?? '',
        created: created ? created.toString() : '',
        type: 'ASSIGN',
        assigned_hours: bountyHours,
        commitment_fee: bountyHours * 200,
        bounty_expires: new Date(
          moment().add(bountyHours, 'hours').format().toString()
        ).toUTCString()
      });

      setLnInvoice(data.response.invoice);
    }
  };

  const onHandle = (event: any) => {
    const res = JSON.parse(event.data);
    if (res.msg === SOCKET_MSG.user_connect) {
      const user = ui.meInfo;
      if (user) {
        user.websocketToken = res.body;
        ui.setMeInfo(user);
      }
    } else if (res.msg === SOCKET_MSG.assign_success && res.invoice === main.lnInvoice) {
      addToast();
      setLnInvoice('');
      setInvoiceStatus(true);
      main.setLnInvoice('');

      // get new wanted list
      main.getPeopleWanteds({ page: 1, resetPage: true });

      props.dismiss();
      if (props.dismissConnectModal) props.dismissConnectModal();
    }
  };

  useEffect(() => {
    const socket: WebSocket = createSocketInstance();

    socket.onopen = () => {
      console.log('Socket connected');
    };

    socket.onmessage = (event: MessageEvent) => {
      onHandle(event);
    };

    socket.onclose = () => {
      console.log('Socket disconnected');
    };
  }, []);

  return (
    <div onClick={(e: any) => e.stopPropagation()}>
      <Modal style={props.modalStyle} overlayClick={() => props.dismiss()} visible={visible}>
        <div style={{ textAlign: 'center', paddingTop: 59, width: 310 }}>
          <div
            style={{ textAlign: 'center', width: '100%', overflow: 'hidden', padding: '0 50px' }}
          >
            <N color={color}>Asign bounty to your self</N>
            <B>Each hour cost 200 sats</B>
            {lnInvoice && ui.meInfo?.owner_pubkey && (
              <Invoice
                startDate={new Date(moment().add(pollMinutes, 'minutes').format().toString())}
                invoiceStatus={invoiceStatus}
                lnInvoice={lnInvoice}
                invoiceTime={pollMinutes}
              />
            )}

            {!lnInvoice && ui.meInfo?.owner_pubkey && (
              <>
                <InvoiceForm>
                  <InvoiceLabel>Number Of Hours</InvoiceLabel>
                  <InvoiceInput
                    type="number"
                    value={bountyHours}
                    onChange={(e: any) => setBountyHours(Number(e.target.value))}
                  />
                </InvoiceForm>
                <Button
                  text={'Generate Invoice'}
                  color={'primary'}
                  style={{ paddingLeft: 25, margin: '12px 0 10px' }}
                  img={'sphinx_white.png'}
                  imgSize={27}
                  height={48}
                  width={'100%'}
                  onClick={generateInvoice}
                />
              </>
            )}
          </div>
        </div>
        <ModalBottomText color={color}>
          <img src="/static/scan_qr.svg" alt="scan" />
          <div className="bottomText">Pay the invoice to assign to your self</div>
        </ModalBottomText>
      </Modal>
      <EuiGlobalToastList toasts={toasts} dismissToast={removeToast} toastLifeTimeMs={3000} />
    </div>
  );
}
