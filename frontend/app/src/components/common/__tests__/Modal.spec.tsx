import React from 'react';
import { render, fireEvent } from '@testing-library/react';
import Modal from '../Modal';

describe('Modal', () => {
  it('handles prevArrowNew click event', () => {
    const mockPrevArrowNew = jest.fn();
    const { getByRole } = render(<Modal prevArrowNew={mockPrevArrowNew} />);
    fireEvent.click(getByRole('button'));
    expect(mockPrevArrowNew).toHaveBeenCalledTimes(1);
  });

  it('handles nextArrowNew click event', () => {
    const mockNextArrowNew = jest.fn();
    const { getByRole } = render(<Modal nextArrowNew={mockNextArrowNew} />);
    fireEvent.click(getByRole('button'));
    expect(mockNextArrowNew).toHaveBeenCalledTimes(1);
  });
});
