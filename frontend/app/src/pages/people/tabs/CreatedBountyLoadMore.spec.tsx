import React from 'react';
import { render, fireEvent, waitFor } from '@testing-library/react';
import mockBounties from '../../../bounties/__mock__/mockBounties.data';
import '@testing-library/jest-dom/extend-expect';
import { Wanted } from './Wanted';

// eslint-disable-next-line @typescript-eslint/typedef
const createdMockBounties = Array.from({ length: 25 }, (_, index) => ({
  ...(mockBounties[0] || {}), 
  bounty: {
    ...(mockBounties[0]?.bounty || {}), 
    id: mockBounties[0]?.bounty?.id + index + 1
  }
}));

jest.mock('../../../bounties/__mock__/mockBounties.data', () => ({
  createdMockBounties
}));

describe('Wanted component', () => {
  it('displays "Load More" button when scrolling down', async () => {
    const { getByText } = render(<Wanted />);
    fireEvent.scroll(window, { target: { scrollY: 1000 } });
    await waitFor(() => {
      expect(getByText('Load More')).toBeInTheDocument();
    });
  });
});
