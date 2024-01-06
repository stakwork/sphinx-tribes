import React from 'react';
import { render, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
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
    const { getByText } = render(<Header />);
    const getSphinxButton = getByText('Get Sphinx');
    fireEvent.click(getSphinxButton);

    expect(openSpy).toHaveBeenCalledWith('https://buy.sphinx.chat/', '_blank');
  });
});
