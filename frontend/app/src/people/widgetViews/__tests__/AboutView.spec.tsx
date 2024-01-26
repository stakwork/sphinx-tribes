import '@testing-library/jest-dom';
import { render, screen } from '@testing-library/react';
import nock from 'nock';
import React from 'react';
import { setupStore } from '../../../__test__/__mockData__/setupStore';
import { user } from '../../../__test__/__mockData__/user';
import { mockUsehistory } from '../../../__test__/__mockFn__/useHistory';
import { AboutView } from '../AboutView';

beforeAll(() => {
  nock.disableNetConnect();
  setupStore();
  mockUsehistory();
});

/**
 * @jest-environment jsdom
 */
describe('AboutView Component', () => {
  nock(user.url).get('/person/id/1').reply(200, {});
  test('display about view with extras', () => {
    const description = 'test description';
    const extras = {
      email: [{ value: 'testEmail@sphinx.com' }],
      twitter: { handle: 'twitterHandle' },
      wanted: []
    };

    render(<AboutView extras={extras} description={description} />);
    expect(screen.queryByText(extras.email[0].value)).toBeInTheDocument();
    // expect(screen.queryByText(`@${extras.twitter.handle}`)).toBeInTheDocument();
    expect(screen.queryByText(description)).toBeInTheDocument();
  });
});
