import React from 'react'
import { useRouteMatch } from 'react-router-dom';

export const PersonPage = () => {
  const { path } = useRouteMatch();
  return (
    <div>{path}</div>
  )
}
