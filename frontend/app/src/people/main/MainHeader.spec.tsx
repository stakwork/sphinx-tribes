import React from 'react';
import { render, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import { MemoryRouter } from 'react-router-dom';
import Header from './Header';

describe('Header Component', () => {
  test('Clicking on "Get Sphinx" button should call the clickHandler function', async () => {
    const mockHistoryPush = jest.fn();
    const mockUseHistory = jest.fn(() => ({ push: mockHistoryPush }));
    const mockWindowOpen = jest.spyOn(window, 'open').mockImplementation(() => window);

    jest.mock('react-router-dom', () => ({
      ...jest.requireActual('react-router-dom'),
      useHistory: mockUseHistory
    }));

    const { getByText } = render(
      <MemoryRouter>
        <Header />
      </MemoryRouter>
    );

    const getSphinxButton = getByText('Get Sphinx');
    fireEvent.click(getSphinxButton);

    await waitFor(() => {
      expect(mockHistoryPush).toHaveBeenCalledWith('/');
      expect(mockWindowOpen).toHaveBeenCalledWith('https://buy.sphinx.chat/');
    });

    mockWindowOpen.mockRestore();
  });
});
