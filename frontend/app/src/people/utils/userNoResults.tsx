import React from 'react';
import { useStores } from '../../store';
import PageLoadSpinner from './pageLoadSpinner';
import { observer } from 'mobx-react-lite';
import NoneSpace from './noneSpace';
import { widgetConfigs } from '../utils/constants';

export default observer(NoResults);
function NoResults(props) {
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
