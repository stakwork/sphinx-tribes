import React from 'react';
import { render, fireEvent, cleanup, waitFor, screen } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';
import HistoryModal from '../organization/HistoryModal';

const props = {
  isOpen: true,
  close: jest.fn(),
  paymentsHistory: [],
  url: ''
};

jest.mock('store/main', () => ({
  getPaymentHistories: jest.fn().mockResolvedValue([
    { id: 1, name: 'Payment 1', amount: 100 },
    { id: 2, name: 'Payment 2', amount: 200 }
  ])
}));

describe('HistoryModal', () => {
  afterEach(() => {
    cleanup();
  });

  it('fetches and displays the payment history', async () => {
    render(<HistoryModal {...props} />);

    await waitFor(() => screen.getByText('Payment 1'));

    expect(screen.getByText('Payment 1')).toBeInTheDocument();
    expect(screen.getByText('$100')).toBeInTheDocument();
    expect(screen.getByText('Payment 2')).toBeInTheDocument();
    expect(screen.getByText('$200')).toBeInTheDocument();
  });
  it('renders correctly', () => {
    const { getByText } = render(<HistoryModal {...props} />);
    expect(getByText('Payment History')).toBeInTheDocument();
  });

  it('changes the text after click', () => {
    const { getByLabelText } = render(<HistoryModal {...props} />);
    expect(getByLabelText(/off/i)).toBeTruthy();
    fireEvent.click(getByLabelText(/off/i));
    expect(getByLabelText(/on/i)).toBeTruthy();
  });
});
