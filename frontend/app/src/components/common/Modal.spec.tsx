import React from 'react';
import { render, fireEvent, screen } from '@testing-library/react';
import Modal from './Modal';

describe('<Modal /> Arrow Buttons', () => {
  test('renders the previous arrow when prevArrowNew is provided', () => {
    const prevArrowFunction = jest.fn();
    render(<Modal prevArrowNew={prevArrowFunction} visible={true} />);
    const prevArrow = screen.getByTestId('prev-arrow');
    expect(prevArrow).toBeInTheDocument();
  });

  test('renders the next arrow when nextArrowNew is provided', () => {
    const nextArrowFunction = jest.fn();
    render(<Modal nextArrowNew={nextArrowFunction} visible={true} />);
    const nextArrow = screen.getByTestId('next-arrow');
    expect(nextArrow).toBeInTheDocument();
  });

  test('calls prevArrowNew function when previous arrow is clicked', () => {
    const prevArrowFunction = jest.fn();
    render(<Modal prevArrowNew={prevArrowFunction} visible={true} />);
    const prevArrow = screen.getByTestId('prev-arrow');
    fireEvent.click(prevArrow);
    expect(prevArrowFunction).toHaveBeenCalled();
  });

  test('calls nextArrowNew function when next arrow is clicked', () => {
    const nextArrowFunction = jest.fn();
    render(<Modal nextArrowNew={nextArrowFunction} visible={true} />);
    const nextArrow = screen.getByTestId('next-arrow');
    fireEvent.click(nextArrow);
    expect(nextArrowFunction).toHaveBeenCalled();
  });
});
