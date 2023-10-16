import React from 'react';
import { observer } from 'mobx-react-lite';
import { widgetConfigs } from './Constants';
import NoneSpace from './NoneSpace';

function OrgNoResults(props: { showAction: boolean; action: () => void }) {
  const { text, icon } = widgetConfigs['organizations'].action;

  if (props.showAction) {
    return (
      <NoneSpace
        small
        style={{
          margin: 'auto',
          marginTop: '10%'
        }}
        action={props.action}
        buttonText={text}
        buttonIcon={icon}
        {...widgetConfigs['organizations']?.noneSpace['noUserResult']}
      />
    );
  }

  return (
    <NoneSpace
      small
      style={{
        margin: 'auto',
        marginTop: '10%'
      }}
      {...widgetConfigs['organizations']?.noneSpace['noResult']}
    />
  );
}
export default observer(OrgNoResults);
