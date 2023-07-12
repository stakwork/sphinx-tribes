import React from 'react';
import { observer } from 'mobx-react-lite';
import { NoResultProps } from 'people/interfaces';
import { widgetConfigs } from './Constants';
import PageLoadSpinner from './PageLoadSpinner';
import NoneSpace from './NoneSpace';

function NoResults(props: NoResultProps) {
  const tabs = widgetConfigs;

  if (props.loading) {
    return <PageLoadSpinner show={true} />;
  } else {
    return (
      <NoneSpace
        small
        style={{
          margin: 'auto',
          marginTop: '25%'
        }}
        {...tabs['usertickets']?.noneSpace['noResult']}
      />
    );
  }
}
export default observer(NoResults);
