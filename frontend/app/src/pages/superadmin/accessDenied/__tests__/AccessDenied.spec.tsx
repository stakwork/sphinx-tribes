import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
import AdminAccessDenied from '../index.tsx';

describe('AdminAccessDenied', () => {
  test("renders Access Denied message, You don't have access message, and Go Back button", () => {
    render(<AdminAccessDenied />);

    const accessDeniedText = screen.getByText(/Access Denied/i);
    expect(accessDeniedText).toBeInTheDocument();

    const noAccessText = screen.getByText(/You don't have access to this page/i);
    expect(noAccessText).toBeInTheDocument();

    const goBackButton = screen.getByRole('button', { name: /Go Back/i });
    expect(goBackButton).toBeInTheDocument();

    fireEvent.click(goBackButton);
  });
});
