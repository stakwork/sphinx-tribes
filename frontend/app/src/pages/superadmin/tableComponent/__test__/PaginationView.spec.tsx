import React from 'react';
import '@testing-library/jest-dom';
import { render, fireEvent } from '@testing-library/react';
import { screen } from '@testing-library/react';
import { MyTable } from '../index';
import { bounties } from '../mockBountyData.ts';

jest.mock('../styles.css');
jest.mock('../index', () => ({
  ...jest.requireActual('../index'), // Use the actual implementation for other exports
  getTotalBounties: jest.fn().mockResolvedValue(50) // Adjust the value as needed
}));

describe('MyTable Component', () => {
  it('should render pagination when bounties length is greater than pageSize', () => {
    //const paginatePrevSpy = jest.fn();
    render(<MyTable bounties={bounties} />);
    screen.debug();
    const paginationArrow1 = screen.queryByAltText('pagination arrow 1');
    const paginationArrow2 = screen.queryByAltText('pagination arrow 2');

    console.log('paginationArrow1:', paginationArrow1);
    console.log('paginationArrow2:', paginationArrow2);

    expect(paginationArrow1).toBeInTheDocument();
    expect(paginationArrow2).toBeInTheDocument();

    if (paginationArrow1) {
      fireEvent.click(paginationArrow1);
    }
  });

  it('should not render pagination when bounties length is less than or equal to pageSize', () => {
    const { queryByAltText } = render(<MyTable bounties={bounties} />);

    const paginationArrow1 = queryByAltText('pagination arrow 1');
    const paginationArrow2 = queryByAltText('pagination arrow 2');

    expect(paginationArrow1).not.toBeInTheDocument();
    expect(paginationArrow2).not.toBeInTheDocument();
  });
});
