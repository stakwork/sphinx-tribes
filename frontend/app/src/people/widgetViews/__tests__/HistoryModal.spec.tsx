import React from 'react';
import { PaymentHistory } from 'store/main';
import { render, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
import { act } from 'react-dom/test-utils';
import HistoryModal from '../organization/HistoryModal';

jest.mock('hooks/uiHooks', () => ({
  useIsMobile: jest.fn().mockReturnValue(false)
}));

jest.mock('store', () => ({
  useStores: jest.fn().mockReturnValue({
    ui: {
      setBountyPerson: jest.fn(),
      meInfo: {
        id: 1
      }
    }
  })
}));

jest.mock('moment', () => () => ({ format: jest.fn(() => '01/01/2022') }));

jest.mock('../../../helpers', () => ({
  formatSat: jest.fn((amount: any) => `${amount} sats`)
}));

jest.mock('../../../public/static/link.svg', () => 'LinkIcon');
jest.mock('../../../public/static/arrow-right.svg', () => 'ArrowRight');

describe('HistoryModal component', () => {
  const mockPaymentsHistory: PaymentHistory[] = [
    {
      id: 1,
      org_uuid: 'organization_uuid',
      status: true,
      updated: '2022-01-01T00:00:00.000Z',
      payment_type: 'payment',
      created: '2022-01-01T00:00:00.000Z',
      amount: 100,
      sender_img: 'senderImage',
      sender_pubkey: 'senderPubkey',
      sender_name: 'Sender Name',
      receiver_img: 'receiverImage',
      receiver_pubkey: 'receiverPubkey',
      receiver_name: 'Receiver Name',
      bounty_id: 123
    }
  ];
  it('renders with mock data', () => {
    const { getByText, getByTestId } = render(
      <HistoryModal
        isOpen
        paymentsHistory={mockPaymentsHistory}
        close={() => {
          jest.fn();
        }}
        url=""
      />
    );

    expect(getByTestId('modal')).toBeInTheDocument();
    expect(getByText('Payment History')).toBeInTheDocument();

    mockPaymentsHistory.forEach((payment: any) => {
      expect(getByText(payment.payment_type)).toBeInTheDocument();
      expect(getByText('01/01/2022')).toBeInTheDocument();
      expect(getByText(`${payment.amount} sats`)).toBeInTheDocument();
      expect(getByText(payment.sender_name)).toBeInTheDocument();
    });

    expect(getByText('Payments')).toBeInTheDocument();
    expect(getByText('Deposit')).toBeInTheDocument();
    expect(getByText('Withdrawals')).toBeInTheDocument();
  });

  it('filters payments based on checkboxes', () => {
    const { getByText, getByLabelText } = render(
      <HistoryModal
        isOpen
        paymentsHistory={mockPaymentsHistory}
        close={() => {
          jest.fn();
        }}
        url=""
      />
    );

    act(() => {
      fireEvent.click(getByLabelText('Payments'));
    });

    expect(getByText('Deposit')).toBeInTheDocument();
    expect(getByText('Withdraw')).toBeInTheDocument();
    expect(getByText('Payment')).not.toBeInTheDocument();
  });
});
