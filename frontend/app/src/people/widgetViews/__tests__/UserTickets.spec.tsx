import '@testing-library/jest-dom';
import { render, waitFor } from '@testing-library/react';
import nock from 'nock';
import React from 'react';
import { setupStore } from '../../../__test__/__mockData__/setupStore';
import { user } from '../../../__test__/__mockData__/user';

import { mockUsehistory } from '../../../__test__/__mockFn__/useHistory';
import routeData from 'react-router';
import { people } from '../../../__test__/__mockData__/persons';
import { userAssignedBounties } from '../../../__test__/__mockData__/userTickets';
import UserTicketsView from '../UserTicketsView';

beforeAll(() => {
  nock.disableNetConnect();
  setupStore();
  mockUsehistory();
});

// Todo : mock api request in usertickets page
describe('UserTicketsView Component', () => {
  let originFetch;
  beforeEach(() => {
    originFetch = (global as any).fetch;
  });
  afterEach(() => {
    (global as any).fetch = originFetch;
  });

  nock(user.url).get('/person/id/1').reply(200, {});

  test('placeholder', () => {});

  /*test('display no assigned tickets when the api request fails, or user has no assigned tickets', async () => {
    const [person] = people;

    const mRes = jest.fn().mockResolvedValueOnce(userAssignedTickets);
    const mockedFetch = jest.fn().mockResolvedValueOnce(mRes as any);
    (global as any).fetch = mockedFetch;

    jest.spyOn(routeData, 'useParams').mockReturnValue({ personPubKey: person.owner_pubkey });
    jest.spyOn(routeData, 'useRouteMatch').mockReturnValue({
      url: `/p/${person.owner_pubkey}/usertickets`,
      path: '/p/:personPubkey/usertickets',
      params: {},
      isExact: true
    });

    (global as any).fetch = mockedFetch;
    const { getByTestId } = render(<UserTicketsView />);

    const div = await waitFor(() => getByTestId('test'));

    expect(div).toHaveTextContent('No Assigned Tickets Yet');
    expect(mockedFetch).toBeCalledTimes(1);
  });*/
});
