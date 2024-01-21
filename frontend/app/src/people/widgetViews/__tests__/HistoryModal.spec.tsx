import React from 'react';
import { PaymentHistory } from 'store/main';
import { render } from '@testing-library/react';
import '@testing-library/jest-dom';
import { MemoryRouter } from 'react-router-dom';
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
    },
    {
      id: 2,
      org_uuid: 'organization_uuid',
      status: true,
      updated: '2022-01-02T00:00:00.000Z',
      payment_type: 'deposit',
      created: '2022-01-02T00:00:00.000Z',
      amount: 200,
      sender_img: 'depositorImage',
      sender_pubkey: 'depositorPubkey',
      sender_name: 'Depositor Name',
      receiver_img: '',
      receiver_pubkey: '',
      receiver_name: '',
      bounty_id: 456
    },
    {
      id: 3,
      org_uuid: 'organization_uuid',
      status: true,
      updated: '2022-01-03T00:00:00.000Z',
      payment_type: 'withdraw',
      created: '2022-01-03T00:00:00.000Z',
      amount: 300,
      sender_img: 'depositorImage',
      sender_pubkey: 'depositorPubkey',
      sender_name: 'Depositor Name',
      receiver_img: 'receiverImage',
      receiver_pubkey: 'receiverPubkey',
      receiver_name: 'Receiver Name',
      bounty_id: 789
    }
  ];

  it('renders with mock data for all transaction types', () => {
    const { getByText, getAllByText } = render(
      <MemoryRouter>
        <HistoryModal
          isOpen
          paymentsHistory={mockPaymentsHistory}
          close={() => {
            jest.fn();
          }}
          url=""
        />
      </MemoryRouter>
    );

    // Check for modal title and filter labels
    expect(getByText('Payment History')).toBeInTheDocument();
    expect(getByText('Payments')).toBeInTheDocument();
    expect(getByText('Deposit')).toBeInTheDocument();
    expect(getByText('Withdrawals')).toBeInTheDocument();

    // Check for each transaction type
    mockPaymentsHistory.forEach((payment: PaymentHistory) => {
      expect(getAllByText(`${payment.amount} sats`).length).toBeGreaterThan(0);
      expect(getAllByText('01/01/2022').length).toBeGreaterThan(0);

      const senderNames = getAllByText(payment.sender_name);
      expect(senderNames.length).toBeGreaterThan(0);

      if (payment.payment_type === 'payment') {
        const receiverNames = getAllByText(payment.receiver_name);
        expect(receiverNames.length).toBeGreaterThan(0);
      }
    });
  });

  it('renders deposit transactions without a receiver', () => {
    const { queryByText } = render(
      <MemoryRouter>
        <HistoryModal
          isOpen
          paymentsHistory={mockPaymentsHistory.filter(
            (payment: PaymentHistory) => payment.payment_type === 'deposit'
          )}
          close={() => {
            jest.fn();
          }}
          url=""
        />
      </MemoryRouter>
    );

    // Deposit should have a sender but no receiver
    expect(queryByText('Depositor Name')).toBeInTheDocument();
    expect(queryByText('Receiver Name')).toBeNull();
  });

  it('renders payment transactions with both sender and receiver', () => {
    const { getByText } = render(
      <MemoryRouter>
        <HistoryModal
          isOpen
          paymentsHistory={mockPaymentsHistory.filter(
            (payment: PaymentHistory) => payment.payment_type === 'payment'
          )}
          close={() => {
            jest.fn();
          }}
          url=""
        />
      </MemoryRouter>
    );

    // Payment should have both a sender and a receiver
    expect(getByText('Sender Name')).toBeInTheDocument();
    expect(getByText('Receiver Name')).toBeInTheDocument();
  });
});
