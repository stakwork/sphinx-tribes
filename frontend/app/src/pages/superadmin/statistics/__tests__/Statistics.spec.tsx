import '@testing-library/jest-dom';
import { render } from '@testing-library/react';
import nock from 'nock';
import React from 'react';
import { setupStore } from '../../../../__test__/__mockData__/setupStore';
import { user } from '../../../../__test__/__mockData__/user';
import { mockUsehistory } from '../../../../__test__/__mockFn__/useHistory';
import { Statistics } from '../';

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
  bounties_paid_average: 78,
  sats_paid_percentage: 50,
  average_paid: 10,
  average_completed: 1,
};

const mockHunter = {
  hunters_total_paid: 125,
  hunters_first_bounty_paid: 16,
}

describe('Statistics Component', () => {
  nock(user.url).get('/person/id/1').reply(200, {});

  it('renders without crashing', () => {
    const { container } = render(<Statistics metrics={mockMetrics} />);
    expect(container.firstChild).toBeInTheDocument();
  });

  it('renders bounties metrics correctly', () => {
    const { getByText, getByTestId } = render(<Statistics metrics={mockMetrics} />);
    expect(getByText('Bounties')).toBeInTheDocument();
    expect(getByText('100')).toBeInTheDocument();
    expect(getByTestId('total_bounties_posted')).toHaveTextContent('Total Posted');
    expect(getByText('50')).toBeInTheDocument();
    expect(getByText('Assigned')).toBeInTheDocument();
    expect(getByText('78')).toBeInTheDocument();
    expect(getByTestId('total_bounties_paid')).toHaveTextContent('Paid');
    expect(getByText('Completed')).toBeInTheDocument();
  });

  it('renders satoshis metrics correctly', () => {
    const { getByText, getByTestId } = render(<Statistics metrics={mockMetrics} />);
    expect(getByText('Satoshis')).toBeInTheDocument();
    expect(getByTestId('total_satoshis_posted')).toHaveTextContent('Total Posted');
    expect(getByTestId('total_satoshis_paid')).toHaveTextContent('Paid');
    expect(getByText('Avg Time to Paid')).toBeInTheDocument();
    expect(getByText('5000')).toBeInTheDocument();
    expect(getByText('2500')).toBeInTheDocument();
  });

  it('renders hunters metrics correctly', () => {
    const { getByText } = render(<Statistics metrics={mockMetrics} mockHunter={mockHunter} />);
    expect(getByText('Hunters')).toBeInTheDocument();
    expect(getByText('Total Paid')).toBeInTheDocument();
    expect(getByText('First Bounty Paid')).toBeInTheDocument();
    expect(getByText('125')).toBeInTheDocument();
    expect(getByText('16')).toBeInTheDocument();
  });
});
