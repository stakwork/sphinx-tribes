import React from 'react';
import { observer } from 'mobx-react-lite';
import { widgetConfigs } from './Constants';
import NoneSpace from './NoneSpace';

function OrgNoResults(props: { showAction: boolean, action: () => void }) {
  const tabs = widgetConfigs;

  if (props.showAction) {
    return (
      <NoneSpace
        small
        style={{
          margin: 'auto',
          marginTop: '10%'
        }}
        action={props.action}
        buttonText={tabs['organizations'].action.text}
        buttonIcon={tabs['organizations'].action.icon}
        {...tabs['organizations']?.noneSpace['noUserResult']}
      />
    )
  }

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
