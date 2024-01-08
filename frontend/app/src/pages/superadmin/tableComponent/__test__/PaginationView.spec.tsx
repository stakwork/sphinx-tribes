import React from 'react';
import '@testing-library/jest-dom';
import { render } from '@testing-library/react';
import { MyTable } from '../index';
import { bounties } from '../mockBountyData.ts';

jest.mock('../styles.css');

describe('MyTable Component', () => {
  it('should render pagination when bounties length is greater than pageSize', () => {
    const { queryByAltText } = render(
      <MyTable bounties={bounties} />
    );

    const paginationArrow1 = queryByAltText('pagination arrow 1');
    const paginationArrow2 = queryByAltText('pagination arrow 2');

    expect(paginationArrow1).toBeInTheDocument();
    expect(paginationArrow2).toBeInTheDocument();
  });

  it('should not render pagination when bounties length is less than or equal to pageSize', () => {
    const { queryByAltText } = render(
      <MyTable bounties={bounties} />
    );

    const paginationArrow1 = queryByAltText('pagination arrow 1');
    const paginationArrow2 = queryByAltText('pagination arrow 2');

    expect(paginationArrow1).not.toBeInTheDocument();
    expect(paginationArrow2).not.toBeInTheDocument();
  });
});
