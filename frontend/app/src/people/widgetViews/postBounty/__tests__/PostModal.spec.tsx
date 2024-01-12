import '@testing-library/jest-dom';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import nock from 'nock';
import React from 'react';
import { setupStore } from '../../../../__test__/__mockData__/setupStore';
import { user } from '../../../../__test__/__mockData__/user';
import { mockUsehistory } from '../../../../__test__/__mockFn__/useHistory';
import { PostModal } from '../PostModal';

beforeAll(() => {
  nock.disableNetConnect();
  setupStore();
  mockUsehistory();
});

describe('Post bounty modal', () => {
  nock(user.url).get('/person/id/1').reply(200, {});

  test('placeholder', () => {});

  /*test('Show and close modal', () => {
    const closeHandler = jest.fn();
    render(<PostModal isOpen={true} onClose={closeHandler} widget="wanted" />);
    expect(screen.getByRole('alertdialog')).toBeInTheDocument();
    fireEvent.click(screen.getByTestId('close-btn'));
    expect(closeHandler).toBeCalledTimes(1);
  });
  test('If modal closed it isnt in the DOM', () => {
    const closeHandler = jest.fn();
    render(<PostModal isOpen={false} onClose={closeHandler} widget="wanted" />);
    expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument();
  });

  const formData = {
    organization: 'organization',
    title: 'title',
    category: 'Web development',
    description: 'description',
    price: 1
  };
  test('FillForm', async () => {
    const closeHandler = jest.fn();
    const successHandler = jest.fn(() => {});
    render(
      <PostModal onSucces={successHandler} isOpen={true} onClose={closeHandler} widget="wanted" />
    );

    // 1step
    expect(screen.queryByText('Choose Bounty type')).toBeInTheDocument();
    const button = await screen.findByText('Start');
    fireEvent.click(button);

    // 2 step
    expect(screen.queryByText('Basic info')).toBeInTheDocument();
    expect(screen.queryByText('Next')).toHaveClass('disableText');
    await waitFor(async () => {
      await userEvent.type(screen.getByLabelText('Bounty Title'), formData.title);
      await userEvent.click(screen.getByTestId('Category'));
      await userEvent.click(screen.getByText(formData.category));
      await userEvent.click(screen.getByText('Next'));
    });

    // 3 step
    expect(screen.queryAllByText('Description')[0]).toBeInTheDocument();
    await waitFor(async () => {
      await userEvent.type(screen.getByLabelText('Description'), formData.description);
      await userEvent.click(screen.getByText('Next'));
    });

    //4 step
    expect(screen.queryByText('Price and Estimate')).toBeInTheDocument();
    await waitFor(async () => {
      await userEvent.type(screen.getByLabelText('Price (Sats)*'), String(formData.price));
      await userEvent.click(screen.getByText('Next'));
    });

    //5 step
    expect(screen.queryByText('Assign Developer')).toBeInTheDocument();
    await waitFor(async () => {
      await userEvent.click(await screen.findByText('Decide Later'));
    });
    expect(screen.queryByText('Finish')).toBeInTheDocument();
  });*/
});
