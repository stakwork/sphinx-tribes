import '@testing-library/jest-dom';
import { render, screen, waitFor } from '@testing-library/react';
import nock from 'nock';
import React from 'react';
import { setupStore } from '../../../__test__/__mockData__/setupStore';
import { user } from '../../../__test__/__mockData__/user';
import { mockUsehistory } from '../../../__test__/__mockFn__/useHistory';
import { Statistics } from './';

beforeAll(() => {
  nock.disableNetConnect();
  setupStore();
  mockUsehistory();
});

/**
 * @jest-environment jsdom
 */

const mockMetrics = {
  bounties_posted: 100,
  bounties_paid: 50,
  sats_posted: 5000,
  sats_paid: 2500,
  bounties_paid_average: 75,
  sats_paid_percentage: 50,
  average_paid: 10,
  average_completed: 10
};

describe('Statistics component', () => {
  test('renders with metrics data', () => {
    // Render the component with mock data
    const { getByText, getByAltText } = render(<Statistics metrics={mockMetrics} />);

    // Assertions for Bounties section
    expect(getByText('Bounties')).toBeInTheDocument();
    expect(getByText('Total Bounties Posted')).toBeInTheDocument();
    expect(getByText('78')).toBeInTheDocument(); // Assuming 78 is a static value in your component
    expect(getByText('Bounties Assigned')).toBeInTheDocument();
    expect(getByText('50')).toBeInTheDocument(); // Replace with the actual value from mockMetrics
    expect(getByText('Bounties Paid')).toBeInTheDocument();
    expect(getByText('75%')).toBeInTheDocument(); // Replace with the actual value from mockMetrics
    expect(getByText('Completed')).toBeInTheDocument();

    // Assertions for Satoshis section
    expect(getByText('Satoshis')).toBeInTheDocument();
    expect(getByText('Total Sats Posted')).toBeInTheDocument();
    expect(getByText('5000')).toBeInTheDocument(); // Replace with the actual value from mockMetrics
    expect(getByText('Sats Paid')).toBeInTheDocument();
    expect(getByText('2500')).toBeInTheDocument(); // Replace with the actual value from mockMetrics
    expect(getByText('Avg Time to Paid')).toBeInTheDocument();
    expect(getByText('3 Days')).toBeInTheDocument();
    expect(getByText('Paid')).toBeInTheDocument();
    expect(getByText('50%')).toBeInTheDocument(); // Replace with the actual value from mockMetrics

    // Additional assertions for images
    expect(getByAltText('Bounties')).toBeInTheDocument();
    expect(getByAltText('Satoshie')).toBeInTheDocument();
  });
});
