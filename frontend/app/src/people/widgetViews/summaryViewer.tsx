import React from 'react';
import { useStores } from '../../store';
import styled from 'styled-components';
import PostSummary from './summaries/postSummary';
import WantedSummary from './summaries/wantedSummary';
import OfferSummary from './summaries/offerSummary';
import { useIsMobile } from '../../hooks';

// this is where we see others posts (etc) and edit our own
export default function SummaryViewer(props: any) {
  const { item, config, person } = props;
  const { ui } = useStores();
  const isMobile = useIsMobile();

  // FIXME, make "AND is me"
  const isSelectedView = ui?.selectedPerson ? true : false;
  const thisIsMine = ui?.selectedPerson === ui?.meInfo?.id;

  function wrapIt(child) {
    return (
      <Wrap
        style={{
          maxHeight: config.name === 'post' || isMobile ? '' : '80vh',
          height: (isSelectedView && thisIsMine) || isMobile ? 'calc(100% - 60px)' : '100%'
        }}
      >
        {child}
      </Wrap>
    );
  }

  switch (config.name) {
    case 'post':
      return wrapIt(<PostSummary {...item} person={person} />);
    case 'offer':
      return wrapIt(<OfferSummary {...item} person={person} />);
    case 'wanted':
      return wrapIt(<WantedSummary {...item} person={person} />);
    default:
      return wrapIt(<div>none</div>);
  }
}

const Wrap = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
  min-width: 100%;
`;
