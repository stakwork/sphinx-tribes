import React from 'react';
import { render, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
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

console.log (createdMockBounties)

jest.mock('../../../bounties/__mock__/mockBounties.data', () => ({
  createdMockBounties
}));

describe('Wanted component', () => {
  it('displays "Load More" button when scrolling down', async () => {
    const { getByText } = render(
      <MemoryRouter>
        <Wanted />
      </MemoryRouter>
    );

    fireEvent.scroll(window, { target: { scrollY: 1000 } });

    // Wait for the component to re-render (assuming there's an API call)
    await waitFor(() => {
      // Check if the "Load More" button is displayed
      if (createdMockBounties.length > 20) {
        expect(getByText('Load More')).toBeInTheDocument();
      } else {
        // If there are not enough bounties, you might want to handle this case
        // or skip the expectation.
        console.warn('Not enough bounties for "Load More" button.');
      }
    });
  });
});
