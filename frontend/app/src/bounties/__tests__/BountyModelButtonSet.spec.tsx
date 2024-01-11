import React from 'react';
import { render } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';
import BountyModalButtonSet from '../BountyModalButtonSet';

describe('BountyModalButtonSet Component', () => {
  it('renders the tribe button when tribe is provided', () => {
    const { queryByText } = render(<BountyModalButtonSet tribe="W3schools" />);
    const tribeButton = queryByText('W3schools');
    expect(tribeButton).toBeInTheDocument();
  });

  it('does not render the tribe button when no tribe is associated', () => {
    const { queryByText } = render(<BountyModalButtonSet />);
    const noTribeText = queryByText('No Tribe');
    expect(noTribeText).toBeInTheDocument();

    const tribeButton = queryByText('W3schools');
    expect(tribeButton).not.toBeInTheDocument(); // Ensure specific tribe button is not rendered

    // @ts-ignore
    const noTribeButtonParent = noTribeText.parentElement;
    expect(noTribeButtonParent).toHaveStyle('pointerEvents: none');
    expect(noTribeButtonParent).toHaveStyle('opacity: 0.5');
  });
});
