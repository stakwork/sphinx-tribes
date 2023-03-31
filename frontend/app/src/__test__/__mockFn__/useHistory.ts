export const mockUsehistory = () => {
  const mockHistoryPush = jest.fn();
  jest.mock('react-router-dom', () => ({
    ...jest.requireActual('react-router-dom'),
    useHistory: () => ({
      push: mockHistoryPush
    })
  }));
};

export const mockUseParams = () => {
  jest.mock('react-router-dom', () => ({
    ...jest.requireActual('react-router-dom'),
    useParams: ({
      personPubKey: 'test_pub_key_2'
    }),
  }));
};

export const mockUseRouterMatch = () => {
  jest.mock('react-router-dom', () => ({
    ...jest.requireActual('react-router-dom'),
    useRouteMatch: jest.fn(),
  }));
};