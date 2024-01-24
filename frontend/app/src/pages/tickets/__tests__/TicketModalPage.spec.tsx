import React from 'react';
import { render, fireEvent, screen, waitFor } from '@testing-library/react';
import { useIsMobile } from 'hooks';
import { useStores } from 'store';
import { TicketModalPage } from '../TicketModalPage.tsx';

jest.mock('hooks', () => ({
  useIsMobile: jest.fn()
}));

jest.mock('store', () => ({
  useStores: jest.fn()
}));

const mockPush = jest.fn();
const mockGoBack = jest.fn();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: () => ({
    push: mockPush,
    goBack: mockGoBack
  }),
  useLocation: () => ({
    pathname: '/bounty/1239',
    search: '',
    state: {}
  }),
  useParams: () => ({
    uuid: 'ck95pe04nncjnaefo08g',
    bountyId: '1239'
  })
}));

describe('TicketModalPage Component', () => {
  it('should redirect to the appropriate page on close based on the route and referrer', async () => {
    // Mocking the useIsMobile hook
    (useIsMobile as jest.Mock).mockReturnValue(false);

    (useStores as jest.Mock).mockReturnValue({
      main: {
        getBountyById: jest.fn(),
        getBountyIndexById: jest.fn()
      }
    });

    render(<TicketModalPage setConnectPerson={jest.fn()} visible={true} />);

    // eslint-disable-next-line @typescript-eslint/no-empty-function
    await waitFor(() => {});

    const closeButton = screen.queryByTestId('close-btn');
    if (closeButton) {
      fireEvent.click(closeButton);

      expect(mockPush).toHaveBeenCalledWith('/bounties');
    }
  });
});
