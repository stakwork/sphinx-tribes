import React from 'react';
import { render, fireEvent } from '@testing-library/react';
import { BrowserRouter as Router } from 'react-router-dom';
import Header from './Header';

let openSpy;

beforeEach(() => {
  openSpy = jest.spyOn(window, 'open');
});

afterEach(() => {
  openSpy.mockRestore();
});

describe('Header Component', () => {
  test('Clicking on "Get Sphinx" button should open a new window with the given URL', () => {
    const { getByText } = render(
      <Router>
        <Header />
      </Router>
    );
    const getSphinxButton = getByText('Get Sphinx');
    fireEvent.click(getSphinxButton);

    expect(openSpy).toEqual('https://buy.sphinx.chat/');
  });
});
