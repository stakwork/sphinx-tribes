import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import { ViewTribe } from '../Components.tsx';
// import { ViewTribe } from '../summaries/wantedSummaries/Components.tsx';

describe('ViewTribe Component', () => {
  it('should display the tribe button enabled when a valid tribe is associated', () => {
    const mockProps = {
      tribe: 'W3schools',
      tribeInfo: { img: 'https://www.w3schools.com/html/pic_trulli.jpg' }
    };

    render(<ViewTribe {...mockProps} />);
    expect(screen.getByText('View Tribe')).toBeInTheDocument();
    expect(screen.getByRole('button')).not.toHaveAttribute('disabled');
    expect(screen.getByRole('button')).toHaveStyle('opacity: 1');
  });

  it('should display the tribe button disabled when no tribe is associated', () => {
    const mockProps = {
      tribe: 'none', // Testing with string 'none' which should be considered invalid
      tribeInfo: null
    };

    render(<ViewTribe {...mockProps} />);
    expect(screen.getByText('View Tribe')).toBeInTheDocument();
    expect(screen.getByRole('button')).toHaveAttribute('disabled');
    expect(screen.getByRole('button')).toHaveStyle('opacity: 0.5');
  });

  it('should display the tribe button disabled when tribe is a falsy value', () => {
    const mockProps = {
      tribe: '', // Testing with an empty string
      tribeInfo: null
    };

    render(<ViewTribe {...mockProps} />);
    expect(screen.getByText('View Tribe')).toBeInTheDocument();
    expect(screen.getByRole('button')).toHaveAttribute('disabled');
    expect(screen.getByRole('button')).toHaveStyle('opacity: 0.5');
  });
});
