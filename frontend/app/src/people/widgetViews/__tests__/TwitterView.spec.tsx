import '@testing-library/jest-dom';
import { render, screen } from '@testing-library/react';
import nock from 'nock';
import React from 'react';
import { setupStore } from '__test__/__mockData__/setupStore';
import { mockUsehistory } from '__test__/__mockFn__/useHistory';
import TwitterView from '../TwitterView';

beforeAll(() => {
  nock.disableNetConnect();
  setupStore();
  mockUsehistory();
});

// Todo : mock api request in usertickets page
describe('TwitterView Component', () => {
  test('It showed display twitter handle', async () => {
    render(<TwitterView handle="test_test" />);
    expect(screen.queryByText('@test_test')).toBeInTheDocument();
  });
});
