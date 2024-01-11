import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import { ViewTribe } from '../Components.tsx';

describe('ViewTribe Component', () => {
  it('should display the tribe button when a tribe is associated', () => {
    const mockProps = {
      tribe: 'W3schools',
      tribeInfo: { img: 'https://www.w3schools.com/html/pic_trulli.jpg' }
    };

    render(<ViewTribe {...mockProps} />);
    expect(screen.getByText('View Tribe')).toBeInTheDocument();
    expect(screen.getByRole('button')).not.toHaveAttribute('disabled');
  });

  it('should not display the tribe button when no tribe is associated', () => {
    const mockProps = {
      tribe: 'none',
      tribeInfo: null
    };

    render(<ViewTribe {...mockProps} />);
    expect(screen.getByText('View Tribe')).toBeInTheDocument();
    expect(screen.getByRole('button')).toHaveAttribute('disabled');
  });
});
