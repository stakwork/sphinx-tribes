/* eslint-disable func-style */
import React from 'react';
import { useStores } from '../../store';
import styled from 'styled-components';
import WantedSummary from './summaries/wantedSummary';
import { useIsMobile } from '../../hooks';
import { observer } from 'mobx-react-lite';

// this is where we see others posts (etc) and edit our own
export default observer(SummaryViewer);

function SummaryViewer(props: any) {
  const { item, config, person } = props;

  const { ui } = useStores();
  const isMobile = useIsMobile();

  // FIXME, make "AND is me"
  const isSelectedView = ui?.selectedPerson ? true : false;
  const thisIsMine = ui?.selectedPerson === ui?.meInfo?.id;


return (      <Wrap
        style={{
          maxHeight: config.name === 'post' || isMobile ? '' : '100vh',
          height: (isSelectedView && thisIsMine) || isMobile ? 'calc(100% - 60px)' : '100%'
        }}
      >
        <WantedSummary
          {...item}
          ReCallBounties={props.ReCallBounties}
          person={person}
          formSubmit={props.formSubmit}
          personBody={props?.personBody}
          fromBountyPage={props?.fromBountyPage}
          extraModalFunction={props?.extraModalFunction}
          deleteAction={props?.deleteAction}
          deletingState={props?.deletingState}
          editAction={props?.editAction}
          setIsModalSideButton={props?.setIsModalSideButton}
          setIsExtraStyle={props?.setIsExtraStyle}
        />
      </Wrap> )
}

const Wrap = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
  min-width: 100%;
  align-items: center;
`;
