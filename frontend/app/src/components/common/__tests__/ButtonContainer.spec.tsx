import '@testing-library/jest-dom';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import nock from 'nock';
import React from 'react';
import { setupStore } from '../../../__test__/__mockData__/setupStore';
import { user } from '../../../__test__/__mockData__/user';
import { mockUsehistory } from '../../../__test__/__mockFn__/useHistory';
import { ButtonContainer } from '../ButtonContainer';

beforeAll(() => {
  nock.disableNetConnect();
  setupStore();
  mockUsehistory();
});

/**
 * @jest-environment jsdom
 */
describe('ButtonContainer', () => {
  nock(user.url).get('/person/id/1').reply(200, {});
  test('is rendered', async () => {
    const { container } = render(<ButtonContainer />);

    //Expect text to be there
    expect(container).toBeTruthy();
  });
});
