import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
import { MemoryRouter } from 'react-router-dom';
import { observer } from 'mobx-react-lite';
import Header from '../Header';

jest.mock('../../../store', () => ({
  useStores: jest.fn(() => ({
    main: {
      getIsAdmin: jest.fn(),
      getSelf: jest.fn()
    },
    ui: {
      meInfo: null,
      setMeInfo: jest.fn(),
      setShowSignIn: jest.fn(),
      setSelectedPerson: jest.fn(),
      setSelectingPerson: jest.fn(),
      showSignIn: false,
      torFormBodyQR: ''
    }
  }))
}));

describe('Header Component', () => {
  test('renders Header component', () => {
    render(<MemoryRouter>{<Header />}</MemoryRouter>);
  });

  test('clicking on the GetSphinxsBtn calls the correct handler', async () => {
    render(<MemoryRouter>{<Header />}</MemoryRouter>);

    const getSphinxsBtn = screen.getByText('Get Sphinx');
    fireEvent.click(getSphinxsBtn);
  });
});
