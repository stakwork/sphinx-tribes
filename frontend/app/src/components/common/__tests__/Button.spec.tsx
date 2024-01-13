import '@testing-library/jest-dom';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import nock from 'nock';
import React from 'react';
import { setupStore } from '../../../__test__/__mockData__/setupStore';
import { user } from '../../../__test__/__mockData__/user';
import { mockUsehistory } from '../../../__test__/__mockFn__/useHistory';
import Button from '../Button';

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
    const queryByTextValue = 'queryByText';
    const mockOnClick = jest.fn();

    render(
      <Button
        id={'1'}
        icon={''}
        height={10}
        width={10}
        disabled={false}
        color={'green'}
        leadingIcon={''}
        text={queryByTextValue}
        onClick={mockOnClick}
      />
    );

    //Expect text to be there
    expect(screen.queryByText(queryByTextValue)).toBeInTheDocument();

    await userEvent.click(screen.getByText(queryByTextValue));

    //If we click on text it should register click
    expect(mockOnClick).toHaveBeenCalledTimes(1);
  });
});
