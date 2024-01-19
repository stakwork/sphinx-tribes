import '@testing-library/jest-dom';
import { act, fireEvent, render } from '@testing-library/react';
import { createMemoryHistory } from 'history';
import nock from 'nock';
import React from 'react';
import { Router } from 'react-router-dom';
import { mainStore } from 'store/main';
import { person } from '__test__/__mockData__/persons';
import { setupStore } from '../../../__test__/__mockData__/setupStore';
import { user } from '../../../__test__/__mockData__/user';
import { mockUsehistory } from '../../../__test__/__mockFn__/useHistory';
import Header from '../Header';
import { uiStore } from 'store/ui';

beforeAll(() => {
  nock.disableNetConnect();
  setupStore();
  mockUsehistory();
});

describe('AboutView Component', () => {
  nock(user.url).get('/person/id/1').reply(200, {});
  nock(user.url).get(`/person/${user.pubkey}`).reply(200, {});

  it('should not navigate to profile page if clicked on profile', async () => {
    jest.spyOn(mainStore, 'getIsAdmin').mockReturnValue(Promise.resolve(false));
    jest.spyOn(mainStore, 'getPersonById').mockReturnValue(Promise.resolve(person));
    jest.spyOn(mainStore, 'getSelf').mockReturnValue(Promise.resolve());
    const history = createMemoryHistory();
    history.push('/p');
    await act(async () => {
      const { getByText } = render(
        <Router history={history}>
          <Header />
        </Router>
      );
      const me = getByText(user.alias);
      fireEvent.click(me);
      expect(history.location.pathname).toEqual(`/p/${user.owner_pubkey}/organizations`);
      expect(history.length).toEqual(3);
    });
  });

  it('should not add go to edit user page if already on it', async () => {
    jest.spyOn(mainStore, 'getIsAdmin').mockReturnValue(Promise.resolve(false));
    jest.spyOn(mainStore, 'getPersonById').mockReturnValue(Promise.resolve(person));
    jest.spyOn(mainStore, 'getSelf').mockReturnValue(Promise.resolve());
    const history = createMemoryHistory();
    history.push(`/p/${user.owner_pubkey}/organizations`);
    await act(async () => {
      const { getByText } = render(
        <Router history={history}>
          <Header />
        </Router>
      );
      const me = getByText(user.alias);
      fireEvent.click(me);
      expect(history.location.pathname).toEqual(`/p/${user.owner_pubkey}/organizations`);
      expect(history.length).toEqual(2);
    });
  });

  it('should render get sphinx button', async () => {
    jest.spyOn(mainStore, 'getIsAdmin').mockReturnValue(Promise.resolve(false));
    jest.spyOn(mainStore, 'getPersonById').mockReturnValue(Promise.resolve(person));
    jest.spyOn(mainStore, 'getSelf').mockReturnValue(Promise.resolve());
    const history = createMemoryHistory();
    await act(async () => {
      const { getByText } = render(
        <Router history={history}>
          <Header />
        </Router>
      );
      const getSphinxButton = getByText('Get Sphinx');
      expect(getSphinxButton).toBeInTheDocument();
    });
  });

  it('should render get sphinx button', async () => {
    jest.spyOn(mainStore, 'getIsAdmin').mockReturnValue(Promise.resolve(false));
    jest.spyOn(mainStore, 'getPersonById').mockReturnValue(Promise.resolve(person));
    jest.spyOn(mainStore, 'getSelf').mockReturnValue(Promise.resolve());

    uiStore.setMeInfo(null);
    const history = createMemoryHistory();
    await act(async () => {
      const { getByText } = render(
        <Router history={history}>
          <Header />
        </Router>
      );
      const signInBtn = getByText('Sign in');
      expect(signInBtn).toBeInTheDocument();
    });
  });

  it('should render user image if signed in', async () => {
    jest.spyOn(mainStore, 'getIsAdmin').mockReturnValue(Promise.resolve(false));
    jest.spyOn(mainStore, 'getPersonById').mockReturnValue(Promise.resolve(person));
    jest.spyOn(mainStore, 'getSelf').mockReturnValue(Promise.resolve());

    uiStore.setMeInfo(user);
    const history = createMemoryHistory();
    await act(async () => {
      const { getByTestId } = render(
        <Router history={history}>
          <Header />
        </Router>
      );
      const userImg = getByTestId('userImg');
      expect(userImg).toHaveAttribute('src', user.img);
    });
  });

  it('should not render sign in button for signed in user', async () => {
    jest.spyOn(mainStore, 'getIsAdmin').mockReturnValue(Promise.resolve(false));
    jest.spyOn(mainStore, 'getPersonById').mockReturnValue(Promise.resolve(person));
    jest.spyOn(mainStore, 'getSelf').mockReturnValue(Promise.resolve());

    uiStore.setMeInfo(user);
    const history = createMemoryHistory();
    await act(async () => {
      const { queryByText } = render(
        <Router history={history}>
          <Header />
        </Router>
      );
      const signInBtn = queryByText('Sign in');
      expect(signInBtn).not.toBeInTheDocument();
    });
  });
});
