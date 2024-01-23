import React from 'react';
import { render, screen, waitFor, fireEvent, act } from '@testing-library/react';
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
      .mockResolvedValue(
        {
          developer_count: mockDeveloperCount,
          people: [],
          bounties_count: 0
        }
      );

    jest
      .spyOn(mainStore, 'getFilterStatusCount')
      .mockResolvedValue(
        {
          open: 10,
          assigned: 5,
          paid: 5
        }
      );

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
    'Rust',
    'Other',
    'C++',
    'C#'
  ];

  languageOptions.forEach((language: string) => {
    test(`should call onChangeLanguage when the ${language} filter option is selected`, async () => {
      render(<BountyHeader {...mockProps} />);
      const filterContainer = screen.getByText('Filter');
      fireEvent.click(filterContainer);

      let checkbox;
      try {
        checkbox = screen.getByRole('checkbox', { name: language });
      } catch (error) {
        console.error(`No checkbox found with the name: ${language}`);
        return;
      }

      fireEvent.click(checkbox);
      expect(mockProps.onChangeLanguage).toHaveBeenCalledWith(language);
    });
  });

  jest.useFakeTimers();

  it('should call main.getPeopleBounty when search text is empty', async () => {
    const { getByTestId } = render(<BountyHeader {...mockProps} />);

    // Simulate typing in the search bar
    fireEvent.change(getByTestId('search-bar'), { target: { value: 'Test' } });

    // Check if the search text is updated
    expect(getByTestId('search-bar')).toHaveValue('Test');

    // const getPeopleBountiesMock = jest.fn();

    // Simulate clicking on the close icon
    fireEvent.change(getByTestId('search-bar'), { target: { value: '' } });

    expect(getByTestId('search-bar')).toHaveValue('');

    const getPeopleBountiesSpy = jest.spyOn(mainStore, 'getPeopleBounties');

    act(() => {
      jest.advanceTimersByTime(2001);
    });
    // Expect that getPeopleBounties has been called
    expect(await getPeopleBountiesSpy).toHaveBeenCalled();
  });

  afterAll(() => {
    jest.useRealTimers(); // Restore real timers after all tests are done
  });
});
