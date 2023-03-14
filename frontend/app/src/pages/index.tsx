import React from 'react';
/* eslint-disable func-style */
import '@material/react-material-icon/dist/material-icon.css';
import { Route, Switch } from 'react-router-dom';
import BotsBody from '../bots/body';
import PeopleBody from '../people/main/body';
import PeopleHeader from '../people/main/header';
import TokenRefresh from '../people/utils/tokenRefresh';
import Body from '../tribes/body';
import Header from '../tribes/header';
import { MainLayout } from './MainLayout';
import { AppMode } from 'config';
import { TicketsPage } from './Tickets';
import { PeoplePage } from './People';
import { Modals } from './Modals';

const modeDispatchPages: Record<AppMode, () => React.ReactElement> = {
  community: () => (
    <MainLayout header={<PeopleHeader />}>
      <TokenRefresh />
      <Switch>
        <Route path="/t/">
          <Body />
        </Route>
        <Route path="/b/">
          <BotsBody />
        </Route>
        <Route path="/p/">
          <PeoplePage />
        </Route>
        <Route path="/tickets/">
          <TicketsPage />
        </Route>
      </Switch>
    </MainLayout>
  ),
  people: () => <></>,
  tribes: () => (
    <MainLayout header={<Header />}>
      <Body />
    </MainLayout>
  )
};

export const Pages = ({ mode }: { mode: AppMode }) => {
  return (
    <>
      {modeDispatchPages[mode]()}
      <Modals />
    </>
  );
};
