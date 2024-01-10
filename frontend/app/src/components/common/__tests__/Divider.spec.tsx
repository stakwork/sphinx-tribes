import '@testing-library/jest-dom';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import nock from 'nock';
import React from 'react';
import { setupStore } from '../../../__test__/__mockData__/setupStore';
import { user } from '../../../__test__/__mockData__/user';
import { mockUsehistory } from '../../../__test__/__mockFn__/useHistory';
import Divider from '../Divider.tsx';

beforeAll(() => {
  nock.disableNetConnect();
  setupStore();
  mockUsehistory();
});

/**
 * @jest-environment jsdom
 */
describe('Button Component', () => {
  nock(user.url).get('/person/id/1').reply(200, {});
  test('display text and click button', async () => {
    render(<Divider style={{ color: 'blue' }} />);

    expect(screen.getByTestId('testid-divider')).toBeInTheDocument();
    expect(screen.getByTestId('testid-divider')).toHaveAttribute('style', 'color: blue;');
  });
});
