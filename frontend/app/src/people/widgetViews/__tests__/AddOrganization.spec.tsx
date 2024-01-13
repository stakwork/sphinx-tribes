import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import AddOrganization from '../organization/AddOrganization';
const mockCloseHandler = jest.fn();
const mockGetUserOrganizations = jest.fn();
const mockOwnerPubKey = 'somePublicKey';

describe('AddOrganization Component Tests', () => {
  test('Organization Name text field appears', () => {
    render(
      <AddOrganization
        closeHandler={mockCloseHandler}
        getUserOrganizations={mockGetUserOrganizations}
        owner_pubkey={mockOwnerPubKey}
      />
    );
    const org = screen.getByPlaceholderText(/My Organization.../i);
    expect(org).toBeInTheDocument();
  });

  test('Website text field appears', () => {
    render(
      <AddOrganization
        closeHandler={mockCloseHandler}
        getUserOrganizations={mockGetUserOrganizations}
        owner_pubkey={mockOwnerPubKey}
      />
    );
    expect(screen.getByPlaceholderText('Website URL...')).toBeInTheDocument();
  });

  test('Github repo text field appears', () => {
    render(
      <AddOrganization
        closeHandler={mockCloseHandler}
        getUserOrganizations={mockGetUserOrganizations}
        owner_pubkey={mockOwnerPubKey}
      />
    );
    expect(screen.getByPlaceholderText('Github link...')).toBeInTheDocument();
  });

  test('Logo button appears', () => {
    render(
      <AddOrganization
        closeHandler={mockCloseHandler}
        getUserOrganizations={mockGetUserOrganizations}
        owner_pubkey={mockOwnerPubKey}
      />
    );
    expect(screen.getByText('LOGO')).toBeInTheDocument();
  });

  test('Description box appears', () => {
    render(
      <AddOrganization
        closeHandler={mockCloseHandler}
        getUserOrganizations={mockGetUserOrganizations}
        owner_pubkey={mockOwnerPubKey}
      />
    );
    expect(screen.getByPlaceholderText('Description Text...')).toBeInTheDocument();
  });

  test('Add Org button appears', () => {
    render(
      <AddOrganization
        closeHandler={mockCloseHandler}
        getUserOrganizations={mockGetUserOrganizations}
        owner_pubkey={mockOwnerPubKey}
      />
    );
    expect(screen.getByText('Add Organization')).toBeInTheDocument();
  });

  test('Clicking on Add Org button triggers an action', async () => {
    const mockCloseHandler = jest.fn();
    const mockGetUserOrganizations = jest.fn();
    const mockOwnerPubKey = 'somePublicKey';

    render(
      <AddOrganization
        closeHandler={mockCloseHandler}
        getUserOrganizations={mockGetUserOrganizations}
        owner_pubkey={mockOwnerPubKey}
      />
    );

    const addButton = screen.getByText('Add Organization');
    expect(addButton).toBeInTheDocument();

    fireEvent.click(addButton);

    await waitFor(() => {
      expect(mockCloseHandler).toHaveBeenCalled();
      expect(mockGetUserOrganizations).toHaveBeenCalled();
    });
  });
});
