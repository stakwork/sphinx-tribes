import React from 'react';
import { render, fireEvent } from '@testing-library/react';
import Modal from '../Modal.tsx';

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: () => ({
    push: jest.fn(),
    goBack: jest.fn()
  })
}));
describe('Modal Component', () => {
  it('should redirect to bounty home page on direct access', () => {
    const historyPushMock = jest.fn();

    // eslint-disable-next-line @typescript-eslint/no-var-requires
    jest.spyOn(require('react-router-dom'), 'useHistory').mockReturnValue({
      push: historyPushMock,
      goBack: jest.fn()
    });

    // eslint-disable-next-line @typescript-eslint/no-var-requires
    jest.spyOn(require('react-router-dom'), 'useLocation').mockReturnValue({
      pathname: '/bounty/1239'
    });

    const { getByTestId } = render(
      <Modal
        visible={true}
        bigCloseImage={() => {
          jest.fn();
        }}
      />
    );

    const closeButton = getByTestId('close-btn');
    fireEvent.click(closeButton);

    expect(historyPushMock).toHaveBeenCalledWith('/bounties');
  });
});
