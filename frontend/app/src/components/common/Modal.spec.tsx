import React from 'react';
import { render, fireEvent, screen } from '@testing-library/react';
import Modal from './Modal';

describe('<Modal /> Arrow Buttons', () => {
  test('calls prevArrowNew function when previous arrow is clicked', () => {
    const prevArrowFunction = jest.fn();
    render(<Modal prevArrowNew={prevArrowFunction} visible={true} />);
    const prevArrow = screen.getByText('chevron_left');
    fireEvent.click(prevArrow);
    expect(prevArrowFunction).toHaveBeenCalled();
  });

  test('calls nextArrowNew function when next arrow is clicked', () => {
    const nextArrowFunction = jest.fn();
    render(<Modal nextArrowNew={nextArrowFunction} visible={true} />);
    const nextArrow = screen.getByText('chevron_right');
    fireEvent.click(nextArrow);
    expect(nextArrowFunction).toHaveBeenCalled();
  });
});
