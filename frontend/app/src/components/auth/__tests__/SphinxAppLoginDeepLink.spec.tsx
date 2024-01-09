import '@testing-library/jest-dom';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import nock from 'nock';
import React from 'react';
import { setupStore } from '../../../__test__/__mockData__/setupStore';
import { user } from '../../../__test__/__mockData__/user';
import { mockUsehistory } from '../../../__test__/__mockFn__/useHistory';
import SphinxAppLoginDeepLink from '../SphinxAppLoginDeepLink';

beforeAll(() => {
  nock.disableNetConnect();
  setupStore();
  mockUsehistory();
});

/**
 * @jest-environment jsdom
 */
describe('SphinxAppLoginDeepLink', () => {
  nock(user.url).get('/person/id/1').reply(200, {});
  nock(user.url).get('/ask').reply(200, {});
  test('display text', async () => {
    render(<SphinxAppLoginDeepLink />);

    //Expect text to be there
    expect(screen.queryByText('Opening Sphinx...')).toBeInTheDocument();
  });
});
