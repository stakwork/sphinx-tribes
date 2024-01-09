import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import BountyHeader from '../BountyHeader';
import { BountyHeaderProps } from '../../interfaces';
import { mainStore } from '../../../store/main';

const mockProps: BountyHeaderProps = {
    selectedWidget: 'people',
    scrollValue: false,
    onChangeStatus: jest.fn(),
    onChangeLanguage: jest.fn(),
    checkboxIdToSelectedMap: {},
    checkboxIdToSelectedMapLanguage: {}
};

describe('BountyHeader Component Tests', () => {

    beforeEach(() => {
        jest.spyOn(mainStore, 'getBountyHeaderData').mockReset();
    });

    test('renders Post a Bounty Button', async () => {
        render(<BountyHeader {...mockProps} />);
        expect(await screen.findByRole('button', { name: /Post a Bounty/i })).toBeInTheDocument();
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
        jest.spyOn(mainStore, 'getBountyHeaderData').mockResolvedValue({ developer_count: mockDeveloperCount });

        render(<BountyHeader {...mockProps} />);

        await waitFor(() => {
            expect(screen.getByText(mockDeveloperCount.toString())).toBeInTheDocument();
        });
    });
});
