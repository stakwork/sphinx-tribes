import React from 'react';
import { render, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter as Router } from 'react-router-dom';
import '@testing-library/jest-dom';
import Header from './Header';

let openSpy;

beforeEach(() => {
  openSpy = jest.spyOn(window, 'open').mockImplementation(() => {});
});

afterEach(() => {
  openSpy.mockRestore();
});

describe('Header Component', () => {
  test('Clicking on "Get Sphinx" button should open a new window with the given URL', async () => {
    const { getByText } = render(
      <Router>
        <Header />
      </Router>
    );
    const getSphinxButton = getByText('Get Sphinx');
    fireEvent.click(getSphinxButton);

    await waitFor(() => expect(openSpy).toHaveBeenCalledWith('https://buy.sphinx.chat/', '_blank'));
  });
});
