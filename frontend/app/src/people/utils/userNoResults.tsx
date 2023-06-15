import React from 'react';
import { observer } from 'mobx-react-lite';
import { NoResultProps } from 'people/interfaces';
import { useStores } from '../../store';
import { widgetConfigs } from '../utils/constants';
import PageLoadSpinner from './pageLoadSpinner';
import NoneSpace from './noneSpace';

export default observer(NoResults);
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
