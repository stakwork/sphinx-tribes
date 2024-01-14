import React from 'react';
import { render, fireEvent, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import { OrgHeader } from '../orgHeader/index';

// Mocking external dependencies
jest.mock('people/widgetViews/postBounty/PostModal', () => ({
  __esModule: true,
  default: () => <div>PostModalMock</div>
}));

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
  it('renders the component', () => {
    render(<OrgHeader />);
    expect(screen.getByText('Post a Bounty')).toBeInTheDocument();
  });

  it('opens the PostBountyModal on button click', () => {
    render(<OrgHeader />);
    fireEvent.click(screen.getByText('Post a Bounty'));
    expect(screen.getByText('ModalTitle')).toBeInTheDocument();
  });

  it('should contain status, skill, and sort by elements', () => {
    render(<OrgHeader />);
    expect(screen.getByLabelText('Status')).toBeInTheDocument();
    expect(screen.getByLabelText('Skill')).toBeInTheDocument();
    expect(screen.getByText('Sort by:Newest First')).toBeInTheDocument();
  });
});
it('matches the snapshot', () => {
  const { asFragment } = render(<OrgHeader />);
  expect(asFragment()).toMatchSnapshot();
});
