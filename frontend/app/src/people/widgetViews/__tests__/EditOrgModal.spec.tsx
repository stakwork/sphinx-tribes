import '@testing-library/jest-dom';
import React from 'react';
import { render, screen } from '@testing-library/react';
import { Organization } from 'store/main';
import EditOrgModal from '../organization/EditOrgModal';

const mockOrganization: Organization = {
  id: '1',
  uuid: 'abc123',
  name: 'Tech Innovators Ltd.',
  website: 'https://test.org',
  github: 'https://github.com/stakwork',
  description: 'Test Descirption',
  owner_pubkey: 'xyz456',
  img: 'https://example.com/logo.png',
  created: '2024-01-15T12:00:00Z',
  updated: '2024-01-15T14:30:00Z',
  show: true,
  bounty_count: 5,
  budget: 100000,
  deleted: false
};

const props = {
  ...mockOrganization,
  isOpen: true,
  onDelete: () => null,
  resetOrg: () => null,
  addToast: () => null,
  close: () => null
};

describe('EditOrgModal Component', () => {
  test('displays the Organization Name text field', () => {
    render(<EditOrgModal {...props} />);
    expect(screen.getAllByText(/Organization Name/i)).toHaveLength(2);
  });

  test('displays the Website text field', () => {
    render(<EditOrgModal {...props} />);
    expect(screen.getAllByText(/Website/i)).toHaveLength(2);
  });

  test('displays the Github repo text field', () => {
    render(<EditOrgModal {...props} />);
    expect(screen.getAllByText(/Github repo/i)).toHaveLength(2);
  });

  test('displays the Description box', () => {
    render(<EditOrgModal {...props} />);
    expect(screen.getAllByText(/Description/i)).toHaveLength(2);
  });

  test('displays the Save changes button', () => {
    render(<EditOrgModal {...props} />);
    expect(screen.getByText(/Save changes/i)).toBeInTheDocument();
  });

  test('displays the Delete button', () => {
    render(<EditOrgModal {...props} />);
    expect(screen.getByText(/Delete Organization/i)).toBeInTheDocument();
  });
});
