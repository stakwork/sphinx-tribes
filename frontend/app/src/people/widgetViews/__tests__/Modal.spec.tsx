import React from 'react';
import { render, fireEvent } from '@testing-library/react';
import { Modal } from 'components/common';

describe('Modal Component', () => {
  it('calls prevArrow callback when clicking the left arrow', () => {
    const prevArrowMock = jest.fn();
    const nextArrowMock = jest.fn();

    const { getByTestId } = render(
      <Modal visible={true} prevArrow={prevArrowMock} nextArrow={nextArrowMock}>
        <div>Modal Content</div>
      </Modal>
    );

    const leftArrowButton = getByTestId('prev-arrow');
    fireEvent.click(leftArrowButton);

    expect(prevArrowMock).toHaveBeenCalledTimes(1);
    expect(nextArrowMock).toHaveBeenCalledTimes(0);
  });

  it('calls nextArrow callback when clicking the right arrow', () => {
    const prevArrowMock = jest.fn();
    const nextArrowMock = jest.fn();

    const { getByTestId } = render(
      <Modal visible={true} prevArrow={prevArrowMock} nextArrow={nextArrowMock}>
        <div>Modal Content</div>
      </Modal>
    );

    const rightArrowButton = getByTestId('next-arrow');
    fireEvent.click(rightArrowButton);

    expect(prevArrowMock).toHaveBeenCalledTimes(0);
    expect(nextArrowMock).toHaveBeenCalledTimes(1);
  });
});
