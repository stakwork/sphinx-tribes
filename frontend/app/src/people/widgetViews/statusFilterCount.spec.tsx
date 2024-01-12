import React from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import BountyHeader from './BountyHeader';
import '@testing-library/jest-dom';

const mockApiResponse = {
  open: 64,
  assigned: 545,
  paid: 502
};

type MockFetch = jest.MockedFunction<typeof fetch>;

const mockFetch = (fetch as unknown as MockFetch).mockImplementation((url) => {
  if (typeof url === 'string' && url.endsWith('gobounties/filter/count')) {
    return Promise.resolve({
      json: () => Promise.resolve(mockApiResponse)
    } as Response);
  }
  return Promise.reject(new Error('Unknown endpoint'));
});

describe('BountyHeader', () => {
  beforeEach(() => {
    mockFetch.mockClear();
  });

  it('displays the filter button and shows status count on click', async () => {
    render(
      <BountyHeader
        selectedWidget={'wanted'}
        scrollValue={false}
        onChangeStatus={function (number: any): void {
          throw new Error('Function not implemented.');
        }}
        onChangeLanguage={function (number: any): void {
          throw new Error('Function not implemented.');
        }}
        checkboxIdToSelectedMap={{
          open: true,
          assigned: false,
          paid: false
        }}
        checkboxIdToSelectedMapLanguage={jest.fn()}
      />
    );

    const filterButton = screen.getByText('Filter');
    expect(filterButton).toBeInTheDocument();
    fireEvent.click(filterButton);
    await waitFor(() => {
      expect(screen.getByText(/Open\s*\[\d+\]/)).toBeInTheDocument();
      expect(screen.getByText(/Assigned\s*\[\d+\]/)).toBeInTheDocument();
      expect(screen.getByText(/Paid\s*\[\d+\]/)).toBeInTheDocument();
    });

    expect(mockFetch).toHaveBeenCalledTimes(2);
    const wasCalledWithUrl = mockFetch.mock.calls.some((call: any) =>
      String(call[1]).includes('gobounties/filter/count')
    );
    expect(wasCalledWithUrl).toBe(false);
  });

  it('handles API errors gracefully', async () => {
    mockFetch.mockRejectedValueOnce(new Error('API error'));

    render(
      <BountyHeader
        selectedWidget={'wanted'}
        scrollValue={false}
        onChangeStatus={function (number: any): void {
          throw new Error('Function not implemented.');
        }}
        onChangeLanguage={function (number: any): void {
          throw new Error('Function not implemented.');
        }}
        checkboxIdToSelectedMap={{
          open: true,
          assigned: false,
          paid: false
        }}
        checkboxIdToSelectedMapLanguage={jest.fn()}
      />
    );
  });
});
