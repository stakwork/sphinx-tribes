import React from 'react';
import { render, fireEvent } from '@testing-library/react';
import AddUserModal from '../organization/AddUserModal'
import {AddUserModalProps} from '../organization/interface'


describe('Width of Modal remains the same before and after input', () => {

    const MockProps: AddUserModalProps = {
    loading: false,
    onSubmit: jest.fn(),
    disableFormButtons: false,
    setDisableFormButtons: jest.fn(),
    isOpen: true,
    close: jest.fn(),
  };

  const { getByTestId, getByText } = render(<AddUserModal {...MockProps} />);
  const modal = getByText('Add New User');

  const getModalMaxWidth = () => {
    const modalStyle = window.getComputedStyle(modal);
    return modalStyle.getPropertyValue('max-width');
  };

  const maxWidthBeforeInput = getModalMaxWidth();
  const searchInput = getByTestId('search-input');

  fireEvent.change(searchInput, { target: { value: 'Saif' } });

  const maxWidthAfterInput = getModalMaxWidth();
  expect(maxWidthBeforeInput).toEqual(maxWidthAfterInput);
});