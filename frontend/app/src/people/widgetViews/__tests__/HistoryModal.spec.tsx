import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';
import HistoryModal from '../organization/HistoryModal';

const mockPaymentHistories = [
  { id: 1, name: 'Payment 1', amount: 100 },
  { id: 2, name: 'Payment 2', amount: 200 }
];

jest.mock('store/main', () => ({
  getPaymentHistories: jest.fn().mockResolvedValue(mockPaymentHistories)
}));

describe('HistoryModal', () => {
  it('fetches and displays the payment history', async () => {
    render(<HistoryModal isOpen={true} close={jest.fn()} paymentsHistory={[]} url="" />);

    // Wait for the asynchronous operation (API request) to complete
    await waitFor(() => {
      expect(screen.getByText('Payment 1')).toBeInTheDocument();
    });
    expect(screen.getByText('$100')).toBeInTheDocument();
    expect(screen.getByText('Payment 2')).toBeInTheDocument();
    expect(screen.getByText('$200')).toBeInTheDocument();
  });
});
