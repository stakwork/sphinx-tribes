
export const mockUseRouter = () => {
    const mockHistoryPush = jest.fn();
    jest.mock('react-router-dom', () => ({
        ...jest.requireActual('react-router-dom') as {},
        useLocation: jest.fn().mockImplementation(() => ({
            pathname: '/bounties',
        })),
        useHistory: () => ({
            push: mockHistoryPush
        })
    }));
};
