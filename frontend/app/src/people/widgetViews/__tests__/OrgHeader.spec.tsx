import React from 'react';
import { render, fireEvent, screen, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import { OrgHeader } from 'pages/tickets/org/orgHeader';

const mockOnChangeLanguage = jest.fn();
const selectedWidget = 'wanted';
describe('OrgHeader Component', () => {
  it('renders the component correctly', () => {
    render(
      <OrgHeader
        onChangeLanguage={mockOnChangeLanguage}
        checkboxIdToSelectedMapLanguage={{}}
        selectedWidget={selectedWidget}
        onChangeStatus={jest.fn()}
        checkboxIdToSelectedMap={{}}
        scrollValue={true}
      />
    );
    expect(screen.getByText('Post a Bounty')).toBeInTheDocument();
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText('Skill')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('Search')).toBeInTheDocument();
    expect(screen.getByText(/Bounties/i)).toBeInTheDocument();
  });

  it('opens the PostModal on "Post a Bounty" button click', async () => {
    render(
      <OrgHeader
        onChangeLanguage={mockOnChangeLanguage}
        checkboxIdToSelectedMapLanguage={{}}
        selectedWidget={selectedWidget}
        onChangeStatus={jest.fn()}
        checkboxIdToSelectedMap={{}}
        scrollValue={true}
      />
    );
    fireEvent.click(screen.getByText('Post a Bounty'));
    // You can add further assertions here to check the modal is open
  });

  it('displays the correct number of bounties', () => {
    render(
      <OrgHeader
        onChangeLanguage={mockOnChangeLanguage}
        checkboxIdToSelectedMapLanguage={{}}
        selectedWidget={selectedWidget}
        onChangeStatus={jest.fn()}
        scrollValue={true}
        checkboxIdToSelectedMap={{}}
      />
    );
    expect(screen.getByText('284')).toBeInTheDocument();
    expect(screen.getByText('Bounties')).toBeInTheDocument();
  });

  it('toggles the SkillFilter on "Skill" dropdown button click', async () => {
    render(
      <OrgHeader
        onChangeLanguage={mockOnChangeLanguage}
        checkboxIdToSelectedMapLanguage={{}}
        selectedWidget={selectedWidget}
        onChangeStatus={jest.fn()}
        checkboxIdToSelectedMap={{}}
        scrollValue={true}
      />
    );
    fireEvent.click(screen.getByTestId('skillDropdown'));
    await waitFor(() => {
      expect(screen.getByTestId('skill-filter')).toBeInTheDocument();
    });
  });
});
