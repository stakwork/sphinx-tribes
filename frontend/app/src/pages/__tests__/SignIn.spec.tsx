import '@testing-library/jest-dom';
import { render, screen } from '@testing-library/react';
import React from 'react';
import Header from 'people/main/Header';

const mockHistoryPush = jest.fn();
jest.mock('react-router', () => ({
    ...jest.requireActual("react-router-dom") as {},
    useLocation: jest.fn().mockImplementation(() => {
        return { pathname: "/bounties" };
    }),
    useHistory: () => ({
        push: mockHistoryPush
    })
}));

/**
 * @jest-environment jsdom
 */
describe('HomePage Component', () => {
    test('check user logged out state', () => {
        const signin = 'Sign in';
        const getSphinx = 'Get Sphinx';

        render(<Header />);
        expect(screen.queryByText(signin)).toBeInTheDocument();
        expect(screen.queryByText(getSphinx)).toBeInTheDocument();
    });
});
