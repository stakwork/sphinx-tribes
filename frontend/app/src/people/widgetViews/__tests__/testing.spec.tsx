import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import BountyHeader from '../BountyHeader';
import { BountyHeaderProps } from '../../interfaces';
import * as hooks from '../../../hooks';

const mockProps: BountyHeaderProps = {
    selectedWidget: 'people',
    scrollValue: false,
    onChangeStatus: jest.fn(),
    onChangeLanguage: jest.fn(),
    checkboxIdToSelectedMap: {},
    checkboxIdToSelectedMapLanguage: {}
};

jest.mock('../../../hooks', () => ({
    useIsMobile: jest.fn(),
    useBountyHeaderData: jest.fn(),
}));

describe('BountyHeader Component Tests', () => {

    beforeEach(() => {
        (hooks.useIsMobile as jest.Mock).mockReturnValue(true);
        (hooks.useBountyHeaderData as jest.Mock).mockReset();
    });

    test('renders Post Bounty Button', () => {
        render(<BountyHeader {...mockProps} />);
        expect(screen.getByRole('button', { name: /Post Bounty/i })).toBeInTheDocument();
    });

    test('renders Leaderboard button', () => {
        render(<BountyHeader {...mockProps} />);
        expect(screen.getByRole('button', { name: /Leaderboard/i })).toBeInTheDocument();
    });

    test('renders search bar', () => {
        render(<BountyHeader {...mockProps} />);
        expect(screen.getByRole('searchbox')).toBeInTheDocument();
    });

    test('renders filters', () => {
        render(<BountyHeader {...mockProps} />);
        expect(screen.getByText(/Filter/i)).toBeInTheDocument();
    });

    test('shows total developer count from mock API', async () => {
        const mockDeveloperCount = 100;
        (hooks.useBountyHeaderData as jest.Mock).mockReturnValue({ developer_count: mockDeveloperCount });

        render(<BountyHeader {...mockProps} />);

        await waitFor(() => {
            expect(screen.getByText(Total Developers: ${mockDeveloperCount.toString()})).toBeInTheDocument();
        });
    });
});
