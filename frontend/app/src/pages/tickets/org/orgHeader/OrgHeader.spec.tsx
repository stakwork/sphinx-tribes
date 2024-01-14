import React from 'react';
import { render, fireEvent, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import { OrgHeader } from '.';

// Mocking external dependencies
jest.mock('people/widgetViews/postBounty/PostModal', () => ({
  __esModule: true,
  default: () => <div>PostModalMock</div>
}));

// Mocking SVG imports
jest.mock('./Icons/addBounty.svg', () => {
  const AddBountyIcon = () => <div>AddBountyIcon</div>;
  AddBountyIcon.displayName = 'AddBountyIcon';
  return AddBountyIcon;
});

jest.mock('./Icons/searchIcon.svg', () => {
  const SearchIcon = () => <div>SearchIcon</div>;
  SearchIcon.displayName = 'AddBountyIcon';
  return SearchIcon;
});

jest.mock('./Icons/file.svg', () => {
  const FileIcon = () => <div>FileIcon</div>;
  FileIcon.displayName = 'AddBountyIcon';
  return FileIcon;
});

describe('OrgHeader Component', () => {
  // Test for component rendering
  test('renders without crashing', () => {
    render(<OrgHeader />);
    expect(screen.getByText('Post a Bounty')).toBeInTheDocument();
  });

  // UI structure tests
  test('contains the necessary UI elements', () => {
    render(<OrgHeader />);
    expect(screen.getByText('Post a Bounty')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('Search')).toBeInTheDocument();
    expect(screen.getByLabelText('Status')).toBeInTheDocument();
    expect(screen.getByLabelText('Skill')).toBeInTheDocument();
    expect(screen.getByLabelText('Sort by:Newest First')).toBeInTheDocument();
  });

  // Interaction test
  test('opens PostModal on button click', () => {
    render(<OrgHeader />);
    fireEvent.click(screen.getByText('Post a Bounty'));
    expect(screen.getByText('PostModalMock')).toBeInTheDocument();
  });

  // State change test
  test('changes state when Post a Bounty is clicked', () => {
    render(<OrgHeader />);
    fireEvent.click(screen.getByText('Post a Bounty'));
    expect(screen.getByText('PostModalMock')).toBeInTheDocument();
  });

  // Snapshot test
  test('matches the snapshot', () => {
    const { asFragment } = render(<OrgHeader />);
    expect(asFragment()).toMatchSnapshot();
  });
});
