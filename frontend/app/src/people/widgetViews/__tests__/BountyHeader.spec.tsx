import React from 'react';
import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
import BountyHeader from '../BountyHeader';
import { BountyHeaderProps } from '../../interfaces';
import { mainStore } from '../../../store/main';
import * as hooks from '../../../hooks';

const mockProps: BountyHeaderProps = {
  selectedWidget: 'wanted',
  scrollValue: false,
  onChangeStatus: jest.fn(),
  onChangeLanguage: jest.fn(),
  checkboxIdToSelectedMap: {},
  checkboxIdToSelectedMapLanguage: {}
};

jest.mock('../../../hooks', () => ({
  useIsMobile: jest.fn()
}));
describe('BountyHeader Component', () => {
  beforeEach(() => {
    jest.spyOn(mainStore, 'getBountyHeaderData').mockReset();
    (hooks.useIsMobile as jest.Mock).mockReturnValue(false);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  test('should render the Post a Bounty button', async () => {
    render(<BountyHeader {...mockProps} />);
    expect(await screen.findByRole('button', { name: /Post a Bounty/i })).toBeInTheDocument();
  });

  test('should render the Leaderboard button', () => {
    render(<BountyHeader {...mockProps} />);
    expect(screen.getByRole('button', { name: /Leaderboard/i })).toBeInTheDocument();
  });

  test('should render the search bar', () => {
    render(<BountyHeader {...mockProps} />);
    expect(screen.getByRole('searchbox')).toBeInTheDocument();
  });

  test('should render the filters', () => {
    render(<BountyHeader {...mockProps} />);
    expect(screen.getByText(/Filter/i)).toBeInTheDocument();
  });

  test('should display the MobileFilterCount with correct number when filters are selected in mobile view', async () => {
    jest.spyOn(hooks, 'useIsMobile').mockReturnValue(true);

    const mockSelectedFilters = {
      checkboxIdToSelectedMap: { filter1: true },
      checkboxIdToSelectedMapLanguage: { lang1: true }
    };

    render(<BountyHeader {...mockProps} {...mockSelectedFilters} />);

    expect(await screen.findByText('2')).toBeInTheDocument();
  });

  test('should display the total developer count from the mock API', async () => {
    jest.setTimeout(20000);
    const mockDeveloperCount = 100;
    jest
      .spyOn(mainStore, 'getBountyHeaderData')
      .mockResolvedValue({ developer_count: mockDeveloperCount });

    render(<BountyHeader {...mockProps} />);

    await waitFor(() => {
      expect(screen.getByText(mockDeveloperCount.toString())).toBeInTheDocument();
    });
  });

  const languageOptions = [
    'Lightning',
    'Typescript',
    'Golang',
    'Kotlin',
    'PHP',
    'C#',
    'Java',
    'Ruby',
    'Python',
    'Postgres',
    'Elastic search',
    'Javascript',
    'Node',
    'Swift',
    'MySQL',
    'R',
    'C++',
    'Rust',
    'Other'
  ];

  languageOptions.forEach((language: string) => {
    test(`should call onChangeLanguage when the ${language} filter option is selected`, async () => {
      render(<BountyHeader {...mockProps} />);
      const filterContainer = screen.getByText('Filter');
      fireEvent.click(filterContainer);

      const checkbox = await screen.findByRole('checkbox', { name: new RegExp(language, 'i') });
      fireEvent.click(checkbox);
      expect(mockProps.onChangeLanguage).toHaveBeenCalledWith(language);
    });
  });
});
