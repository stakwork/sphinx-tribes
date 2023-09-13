import React from 'react';
import { observer } from 'mobx-react-lite';
import { widgetConfigs } from './Constants';
import NoneSpace from './NoneSpace';

function OrgNoResults() {
  const tabs = widgetConfigs;

  return (
    <NoneSpace
      small
      style={{
        margin: 'auto',
        marginTop: '10%'
      }}
      {...tabs['organizations']?.noneSpace['noResult']}
    />
  );
}
export default observer(OrgNoResults);
