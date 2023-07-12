import '@testing-library/jest-dom';
import { render, screen } from '@testing-library/react';
import nock from 'nock';
import React from 'react';
import { setupStore } from '__test__/__mockData__/setupStore';
import { mockUsehistory } from '__test__/__mockFn__/useHistory';
import BlogView from '../BlogView';
import { blowViews } from '__test__/__mockData__/blogViews';

beforeAll(() => {
    nock.disableNetConnect();
    setupStore();
    mockUsehistory();
});

// Todo : mock api request in usertickets page
describe('BlogView Component', () => {
    test('It showed display blog title', async () => {
        render(<BlogView {...blowViews} />);
        expect(screen.queryByText(blowViews.title)).toBeInTheDocument();
    });
});
