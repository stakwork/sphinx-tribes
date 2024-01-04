import '@testing-library/jest-dom';
import { act, render, waitFor } from '@testing-library/react';
import { person } from '__test__/__mockData__/persons';
import { setupStore } from '__test__/__mockData__/setupStore';
import { user } from '__test__/__mockData__/user';
import { mockUsehistory } from '__test__/__mockFn__/useHistory';
import mockBounties from 'bounties/__mock__/mockBounties.data';
import nock from 'nock';
import React from 'react';
import { MemoryRouter, Route } from 'react-router-dom';
import { mainStore } from 'store/main';
import UserTickets from '../UserTicketsView';

beforeAll(() => {
  nock.disableNetConnect();
  setupStore();
  mockUsehistory();
});

/**
 * @jest-environment jsdom
 */
describe('UserTickets', () => {
  nock(user.url).get('/person/id/1').reply(200, { user });
  nock(user.url).get('/ask').reply(200, {});

  it('renders no tickets if none assigned', async () => {
    jest.spyOn(mainStore, 'getPersonAssignedBounties').mockReturnValue(Promise.resolve([]));
    await act(async () => {
      const { getByText } = render(
        <MemoryRouter initialEntries={['/p/1234/usertickets']}>
          <Route path="/p/:personPubkey/usertickets" component={UserTickets} />
        </MemoryRouter>
      );

      await waitFor(() => getByText('No Assigned Tickets Yet'));

      expect(getByText('No Assigned Tickets Yet')).toBeInTheDocument();
    });
  });

  it('renders user tickets', async () => {
    const userBounty = { ...mockBounties[0], body: {} } as any;
    userBounty.body = {
      ...mockBounties[0].bounty,
      owner_id: person.owner_pubkey,
      title: 'test bounty here'
    } as any;
    jest
      .spyOn(mainStore, 'getPersonAssignedBounties')
      .mockReturnValue(Promise.resolve([userBounty]));
    await act(async () => {
      const { getByText } = render(
        <MemoryRouter initialEntries={['/p/1234/usertickets']}>
          <Route path="/p/:personPubkey/usertickets" component={UserTickets} />
        </MemoryRouter>
      );

      await waitFor(() => getByText(userBounty.body.title));

      expect(getByText(userBounty.body.title)).toBeInTheDocument();
    });
  });

  it('renders title, description matching the bounty', async () => {
    const userBounty = { ...mockBounties[0], body: {} } as any;
    userBounty.body = {
      ...mockBounties[0].bounty,
      owner_id: person.owner_pubkey,
      title: 'test bounty here',
      description: 'custom ticket for testing'
    } as any;
    jest
      .spyOn(mainStore, 'getPersonAssignedBounties')
      .mockReturnValue(Promise.resolve([userBounty]));
    await act(async () => {
      const { getByText } = render(
        <MemoryRouter initialEntries={['/p/1234/usertickets']}>
          <Route path="/p/:personPubkey/usertickets" component={UserTickets} />
        </MemoryRouter>
      );

      await waitFor(() => getByText(userBounty.body.title));
      await waitFor(() => getByText(userBounty.body.description));

      expect(getByText(userBounty.body.title)).toBeInTheDocument();
      expect(getByText(userBounty.body.description)).toBeInTheDocument();
    });
  });

  it('renders price matching the bounty', async () => {
    const userBounty = { ...mockBounties[0], body: {} } as any;
    userBounty.body = {
      ...mockBounties[0].bounty,
      owner_id: person.owner_pubkey,
      price: 100
    } as any;
    jest
      .spyOn(mainStore, 'getPersonAssignedBounties')
      .mockReturnValue(Promise.resolve([userBounty]));
    await act(async () => {
      const { getByText } = render(
        <MemoryRouter initialEntries={['/p/1234/usertickets']}>
          <Route path="/p/:personPubkey/usertickets" component={UserTickets} />
        </MemoryRouter>
      );

      await waitFor(() => getByText('100'));

      expect(getByText('100')).toBeInTheDocument();
    });
  });

  it('renders estimated time matching the bounty', async () => {
    const userBounty = { ...mockBounties[0], body: {} } as any;
    userBounty.body = {
      ...mockBounties[0].bounty,
      owner_id: person.owner_pubkey,
      estimated_session_length: '< 3hrs'
    } as any;
    jest
      .spyOn(mainStore, 'getPersonAssignedBounties')
      .mockReturnValue(Promise.resolve([userBounty]));
    await act(async () => {
      const { getByText } = render(
        <MemoryRouter initialEntries={['/p/1234/usertickets']}>
          <Route path="/p/:personPubkey/usertickets" component={UserTickets} />
        </MemoryRouter>
      );

      await waitFor(() => getByText(userBounty.body.estimated_session_length));

      expect(getByText(userBounty.body.estimated_session_length)).toBeInTheDocument();
    });
  });
});
