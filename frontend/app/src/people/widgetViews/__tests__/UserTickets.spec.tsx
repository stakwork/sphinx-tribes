import '@testing-library/jest-dom';
import { render, screen } from '@testing-library/react';
import nock from 'nock';
import React from 'react';
import { setupStore } from '__test__/__mockData__/setupStore';
import { user } from '__test__/__mockData__/user';
import { mockUsehistory } from '__test__/__mockFn__/useHistory';
import UserTicketsView from '../userTicketsView';
import routeData from 'react-router';
import { people } from '__test__/__mockData__/persons';

beforeAll(() => {
    nock.disableNetConnect();
    setupStore();
    mockUsehistory();
});

// Todo : mock api request in usertickets page
describe('UserTicketsView Component', () => {
    nock(user.url);
    test('display user assigned tickets', () => {
        const person = people[1];

        jest.spyOn(routeData, 'useParams').mockReturnValue({ personPubKey: person.owner_pubkey });
        jest.spyOn(routeData, 'useRouteMatch').mockReturnValue({ "url": `/p/${person.owner_pubkey}/usertickets`, "path": "/p/:personPubkey/usertickets", "params": {}, "isExact": true });

        const title = 'Trying this';

        render(<UserTicketsView />);
        expect(screen.queryByText(title)).toBeInTheDocument();
    });
});