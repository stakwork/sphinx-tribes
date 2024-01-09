import * as React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import BountyHeader from '../BountyHeader';
import { BountyHeaderProps } from '../../interfaces';
import * as hooks from '../../../hooks';

const mockProps: BountyHeaderProps = {
  selectedWidget: 'people',
  scrollValue: false,
  onChangeStatus: jest.fn(),
  onChangeLanguage: jest.fn(),
  checkboxIdToSelectedMap: {},
  checkboxIdToSelectedMapLanguage: {}
};

jest.mock('../../../hooks', () => ({
  useIsMobile: jest.fn()
}));

describe('BountyHeader Component Tests', () => {
  beforeEach(() => {
    (hooks.useIsMobile as jest.Mock).mockReturnValue(true);
  });

  test('renders filters', () => {
    render(<BountyHeader {...mockProps} />);
    expect(screen.getByText(/Filter/i)).toBeInTheDocument();
  });
});
