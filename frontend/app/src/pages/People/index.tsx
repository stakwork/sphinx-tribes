import { PeopleBody } from 'people/main'
import React from 'react'
import { Route, Switch, useRouteMatch } from 'react-router-dom'
import { PersonPage } from './PersonPage';

export const PeoplePage = () => {
  const { path } = useRouteMatch();
  console.log(path)
  return (
    <Switch>
      {/* <Route path={`${path}:personId/`}>
        <PersonPage />
      </Route> */}
      <Route path={`${path}`}>
        {/* <PersonPage /> */}
      <PeopleBody selectedWidget="people" />
      </Route>
    </Switch>
  )
}
