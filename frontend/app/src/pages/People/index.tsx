import React from 'react';
import { Route, Switch, useRouteMatch, BrowserRouter } from 'react-router-dom';
import { PeoplePage } from './PeoplePage';
import { PersonPage } from './PersonPage';

export const People = () => {
  const { path } = useRouteMatch();
  return (
    <BrowserRouter>
      <Switch>
        <Route path={`${path}:personPubkey/`}>
          <PersonPage />
        </Route>
        <Route path={`${path}`}>
          <PeoplePage />
        </Route>
      </Switch>
    </BrowserRouter>
  );
};
