import React from 'react';
import { render, fireEvent, screen } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';
import Modal from '../Modal';

describe('Modal Component Tests', () => {
  const mockCloseFunction = jest.fn();
  const mockOverlayClick = jest.fn();
  const mockPrevArrow = jest.fn();
  const mockNextArrow = jest.fn();
  const mockPrevArrowNew = jest.fn();
  const mockNextArrowNew = jest.fn();
  const mockBigClose = jest.fn();

  const defaultProps = {
    visible: true,
    close: mockCloseFunction,
    overlayClick: mockOverlayClick,
    prevArrow: mockPrevArrow,
    nextArrow: mockNextArrow,
    prevArrowNew: mockPrevArrowNew,
    nextArrowNew: mockNextArrowNew,
    bigClose: mockBigClose
  };

  it('renders without crashing', () => {
    render(<Modal {...defaultProps} />);
    expect(screen.getByTestId('modal-env')).toBeInTheDocument();
  });

  it('is visible when visible prop is true', () => {
    render(<Modal {...defaultProps} />);
    expect(screen.getByTestId('modal-fadeleft')).toBeVisible();
  });

  it('is not visible when visible prop is false', () => {
    render(<Modal {...{ ...defaultProps, visible: false }} />);
    expect(screen.queryByTestId('modal-fadeleft')).not.toBeInTheDocument();
  });

  it('calls overlayClick when the overlay is clicked', () => {
    render(<Modal {...defaultProps} />);
    fireEvent.click(screen.getByTestId('modal-fadeleft'));
    expect(mockOverlayClick).toHaveBeenCalled();
  });

  it('calls close function when close button is clicked', () => {
    render(<Modal {...defaultProps} />);
    fireEvent.click(screen.getByTestId('modal-close-button'));
    expect(mockCloseFunction).toHaveBeenCalled();
  });

  it('calls bigClose function when big close button is clicked', () => {
    render(<Modal {...defaultProps} />);
    fireEvent.click(screen.getByTestId('modal-bigclose-button'));
    expect(mockBigClose).toHaveBeenCalled();
  });

  it('calls prevArrow function when previous arrow is clicked', () => {
    render(<Modal {...defaultProps} />);
    fireEvent.click(screen.getByTestId('modal-prev-arrow'));
    expect(mockPrevArrow).toHaveBeenCalled();
  });

  it('calls nextArrow function when next arrow is clicked', () => {
    render(<Modal {...defaultProps} />);
    fireEvent.click(screen.getByTestId('modal-next-arrow'));
    expect(mockNextArrow).toHaveBeenCalled();
  });

  it('calls prevArrowNew function when new previous arrow is clicked', () => {
    render(<Modal {...defaultProps} />);
    fireEvent.click(screen.getByTestId('modal-prev-arrow-new'));
    expect(mockPrevArrowNew).toHaveBeenCalled();
  });

  it('calls nextArrowNew function when new next arrow is clicked', () => {
    render(<Modal {...defaultProps} />);
    fireEvent.click(screen.getByTestId('modal-next-arrow-new'));
    expect(mockNextArrowNew).toHaveBeenCalled();
  });
});

