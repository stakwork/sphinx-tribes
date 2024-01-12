import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import BountyHeader from './BountyHeader'; // Update with the correct path

// Define the mock API response
const mockApiResponse = {
  open: 64,
  assigned: 545,
  paid: 502
};

// Type definition for the mocked fetch response
type MockFetch = jest.MockedFunction<typeof fetch>;

// Create a mock implementation of fetch
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

  it('fetches and displays the status counts correctly', async () => {
    render(
      <BountyHeader
        selectedWidget={'people'}
        scrollValue={false}
        onChangeStatus={function (number: any): void {
          throw new Error('Function not implemented.');
        }}
        onChangeLanguage={function (number: any): void {
          throw new Error('Function not implemented.');
        }}
        checkboxIdToSelectedMap={undefined}
        checkboxIdToSelectedMapLanguage={undefined}
      />
    );

    // Verify that the status counts are displayed correctly
    await waitFor(() => {
      expect(screen.getByText('Open [64]')).toBeInTheDocument();
      expect(screen.getByText('Assigned [545]')).toBeInTheDocument();
      expect(screen.getByText('Paid [502]')).toBeInTheDocument();
    });

    // Ensure that the fetch was called correctly
    expect(mockFetch).toHaveBeenCalledTimes(1);
    expect(mockFetch).toHaveBeenCalledWith(expect.stringContaining('gobounties/filter/count'));
  });

  it('handles API errors gracefully', async () => {
    mockFetch.mockRejectedValueOnce(new Error('API error'));

    render(
      <BountyHeader
        selectedWidget={'people'}
        scrollValue={false}
        onChangeStatus={function (number: any): void {
          throw new Error('Function not implemented.');
        }}
        onChangeLanguage={function (number: any): void {
          throw new Error('Function not implemented.');
        }}
        checkboxIdToSelectedMap={undefined}
        checkboxIdToSelectedMapLanguage={undefined}
      />
    );

    // Replace 'Error loading bounty counts' with the actual error message your component displays
    await waitFor(() => {
      expect(screen.queryByText('Open [64]')).not.toBeInTheDocument();
      expect(screen.getByText('Error loading bounty counts')).toBeInTheDocument();
    });

    expect(mockFetch).toHaveBeenCalledWith(expect.stringContaining('gobounties/filter/count'));
  });

  // Additional test cases can be added here as needed
});

export {};
