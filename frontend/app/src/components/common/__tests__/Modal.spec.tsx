import React from 'react';
import { render, fireEvent } from '@testing-library/react';
import Modal from '../Modal';

describe('Modal', () => {
  it('handles prevArrowNew click event', () => {
    const mockPrevArrowNew = jest.fn();
    const { getByText } = render(<Modal prevArrowNew={mockPrevArrowNew} />);
    fireEvent.click(getByText('chevron_left'));
    expect(mockPrevArrowNew).toHaveBeenCalledTimes(1);
  });

  it('handles nextArrowNew click event', () => {
    const mockNextArrowNew = jest.fn();
    const { getByText } = render(<Modal nextArrowNew={mockNextArrowNew} />);
    fireEvent.click(getByText('chevron_right'));
    expect(mockNextArrowNew).toHaveBeenCalledTimes(1);
  });
});
