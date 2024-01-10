import '@testing-library/jest-dom';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import nock from 'nock';
import React from 'react';
import { setupStore } from '../../../__test__/__mockData__/setupStore';
import { user } from '../../../__test__/__mockData__/user';
import { mockUsehistory } from '../../../__test__/__mockFn__/useHistory';
import QR from '../QR';

beforeAll(() => {
  nock.disableNetConnect();
  setupStore();
  mockUsehistory();
});

/**
 * @jest-environment jsdom
 */
describe('QR Component', () => {
  nock(user.url).get('/person/id/1').reply(200, {});
  test('does display qr code', async () => {
    const value = 'test value';
    render(<QR type={'connect'} value={value} size={10} />);

    //Expect text to be there
    expect(screen.getByTestId('testid-qrcode')).toBeInTheDocument();
    expect(screen.getByTestId('testid-connectimg')).toBeInTheDocument();
  });
});
