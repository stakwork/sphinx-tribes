import { PeopleBody } from 'people/main';
import React from 'react';
import { Route, Switch, useRouteMatch } from 'react-router-dom';
import { PersonPage } from './PersonPage';
import { PeoplePage } from './PeoplePage';

export const People = () => {
  const { path } = useRouteMatch();
  console.log(path);
  return (
    <Switch>
      <Route path={`${path}:personPubkey/`}>
        <PersonPage />
      </Route>
      <Route path={`${path}`}>
        {/* <PeoplePage /> */}
        <PeopleBody selectedWidget="people" />
      </Route>
    </Switch>
  );
};
