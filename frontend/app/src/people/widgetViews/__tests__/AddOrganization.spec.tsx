import '@testing-library/jest-dom';
import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import AddOrganization from '../organization/AddOrganization';
import { mainStore } from 'store/main';
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
    expect(screen.getByPlaceholderText('My Organization...')).toBeInTheDocument();
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
    jest.spyOn(mainStore, 'addOrganization').mockReturnValueOnce(Promise.resolve({
        status: 200,
        json: () => Promise.resolve({})
    }));

    render(
      <AddOrganization
        closeHandler={mockCloseHandler}
        getUserOrganizations={mockGetUserOrganizations}
        owner_pubkey={mockOwnerPubKey}
      />
    );

    const addButton = screen.getByText('Add Organization');
    expect(addButton).toBeInTheDocument();
    const orgNameInput = screen.getByPlaceholderText(/My Organization.../i);
    fireEvent.change(orgNameInput, { target: { value: 'My Org' } });

    fireEvent.click(addButton);

    await waitFor(() => {
      expect(mockCloseHandler).toHaveBeenCalled();
      expect(mockGetUserOrganizations).toHaveBeenCalled();
    });
  });
});
