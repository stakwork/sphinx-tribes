import React from 'react';
import { render, fireEvent, screen } from '@testing-library/react';
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

  it('should retain the same width before and after inputting text', () => {
  
    render(<AddUserModal {...MockProps} />);
   
    const modal = screen.getByRole('dialog');
    
    const initialWidth = modal.offsetWidth;
  
    const input = screen.getByRole('textbox');

    fireEvent.change(input, { target: { value: 'Saif' } });
    
    const widthAfterInput = modal.offsetWidth;
  
    expect(initialWidth).toBe(widthAfterInput);
  });
});
