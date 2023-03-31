import '@testing-library/jest-dom';
import { render, screen } from '@testing-library/react';
import nock from 'nock';
import React from 'react';
import { setupStore } from '__test__/__mockData__/setupStore';
import { user } from '__test__/__mockData__/user';
import { mockUsehistory, mockUseParams, mockUseRouterMatch } from '__test__/__mockFn__/useHistory';
import UserTicketsView from '../userTicketsView';
import routeData from 'react-router';
import { match } from 'react-router-dom';

beforeAll(() => {
    nock.disableNetConnect();
    setupStore();
    mockUsehistory();
    mockUseParams();
    mockUseRouterMatch();
});

describe('UserTicketsView Component', () => {
    nock(user.url);
    test('display user assigned tickets', () => {
        jest.spyOn(routeData, 'useParams').mockReturnValue({ personPubKey: 'test_pub_key_2' });
        jest.spyOn(routeData, 'useRouteMatch').mockReturnValue({ "url": "hello", "path": "jell", "params": {}, "isExact": true });

        const description = 'Just for terst';
        const title = 'Trying this';

        render(<UserTicketsView />);
        expect(screen.queryByText(title)).toBeInTheDocument();
        expect(screen.queryByText(description)).toBeInTheDocument();
    });
});