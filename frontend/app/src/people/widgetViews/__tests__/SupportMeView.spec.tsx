import '@testing-library/jest-dom';
import { render, screen } from '@testing-library/react';
import nock from 'nock';
import React from 'react';
import { setupStore } from '../../../__test__/__mockData__/setupStore';
import { mockUsehistory } from '../../../__test__/__mockFn__/useHistory';
import SupportMeView from '../SupportMeView';
import { user } from '../../../__test__/__mockData__/user';

beforeAll(() => {
  nock.disableNetConnect();
  setupStore();
  mockUsehistory();
});

// Todo : mock api request in usertickets page
describe('SupportMeView Component', () => {
  nock(user.url).get('/person/id/1').reply(200, {});
  const data = {
    title: 'Support Title',
    description: 'Support Description',
    created: 12345678,
    show: true
  };
  test('It showed display description', async () => {
    render(<SupportMeView {...data} />);
    expect(screen.queryByText(data.description)).toBeInTheDocument();
    expect(screen.queryByText('No link')).toBeInTheDocument();
  });
});
