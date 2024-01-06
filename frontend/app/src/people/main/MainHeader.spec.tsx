import React from 'react';
import { render, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
import Header from './Header';

describe('Header Component', () => {
  test('Clicking on "Get Sphinx" button should call the clickHandler function', () => {
    const mockHistoryPush = jest.fn();
    const mockUseHistory = jest.fn(() => ({ push: mockHistoryPush }));

    jest.mock('react-router-dom', () => ({
      ...jest.requireActual('react-router-dom'),
      useHistory: mockUseHistory
    }));

    const { getByText } = render(<Header />);

    const getSphinxButton = getByText('Get Sphinx');
    fireEvent.click(getSphinxButton);
    expect(mockHistoryPush).toHaveBeenCalledWith('/');
  });
});
