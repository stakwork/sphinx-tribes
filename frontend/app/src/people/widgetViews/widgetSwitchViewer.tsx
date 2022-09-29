import React from 'react';
import OfferView from '../widgetViews/offerView';
import WantedView from '../widgetViews/wantedView';
import PostView from '../widgetViews/postView';
import styled from 'styled-components';
import { useIsMobile } from '../../hooks';
import { useStores } from '../../store';
import { useObserver } from 'mobx-react-lite';
import { widgetConfigs } from '../utils/constants';
import { Spacer } from '../main/body';
import NoResults from '../utils/noResults';

export default function WidgetSwitchViewer(props) {
  const { main } = useStores();
  const isMobile = useIsMobile();

  const panelStyles = isMobile
    ? {
        minHeight: 132
      }
    : {
        maxWidth: 291,
        minWidth: 291,
        marginRight: 20,
        marginBottom: 20,
        minHeight: 472
      };

  return useObserver(() => {
    const { peoplePosts, peopleWanteds, peopleOffers } = main;

    let { selectedWidget, onPanelClick } = props;

    if (!selectedWidget) {
      return <div style={{ height: 200 }} />;
    }

    const listSource = {
      post: peoplePosts,
      wanted: peopleWanteds,
      offer: peopleOffers
    };

    const activeList = listSource[selectedWidget];

    let searchKeys: any = widgetConfigs[selectedWidget]?.schema?.map((s) => s.name) || [];
    let foundDynamicSchema = widgetConfigs[selectedWidget]?.schema?.find((f) => f.dynamicSchemas);
    // if dynamic schema, get all those fields
    if (foundDynamicSchema) {
      let dynamicFields: any = [];
      foundDynamicSchema.dynamicSchemas?.forEach((ds) => {
        ds.forEach((f) => {
          if (!dynamicFields.includes(f.name)) dynamicFields.push(f.name);
        });
      });
      searchKeys = dynamicFields;
    }

    const listItems =
      activeList && activeList.length ? (
        activeList.map((item, i) => {
          const { person, body } = item;

          const conditionalStyles = body?.paid
            ? {
                border: isMobile ? '2px 0 0 0 solid #dde1e5' : '1px solid #dde1e5',
                boxShadow: 'none'
              }
            : {};

          // if this person has entries for this widget
          return (
            <Panel
              isMobile={isMobile}
              key={person?.owner_pubkey + i + body?.created}
              onClick={() => {
                if (onPanelClick) onPanelClick(person, body);
              }}
              style={{
                ...panelStyles,
                ...conditionalStyles,
                cursor: 'pointer',
                padding: 0,
                overflow: 'hidden'
              }}>
              {selectedWidget === 'post' ? (
                <PostView
                  showName
                  key={i + person.owner_pubkey + 'pview'}
                  person={person}
                  {...body}
                />
              ) : selectedWidget === 'offer' ? (
                <OfferView
                  showName
                  key={i + person.owner_pubkey + 'oview'}
                  person={person}
                  {...body}
                />
              ) : selectedWidget === 'wanted' ? (
                <WantedView
                  showName
                  key={i + person.owner_pubkey + 'wview'}
                  person={person}
                  {...body}
                />
              ) : null}
            </Panel>
          );
        })
      ) : (
        <NoResults />
      );

    return (
      <>
        {listItems}
        <Spacer key={'spacer'} />
      </>
    );
  });
}

interface PanelProps {
  isMobile?: boolean;
}

const Panel = styled.div<PanelProps>`
  position: relative;
  background: #ffffff;
  color: #000000;
  padding: 20px;
  box-shadow: ${(p) => (p.isMobile ? 'none' : '0px 0px 6px rgb(0 0 0 / 7%)')};
  border-bottom: ${(p) => (p.isMobile ? '2px solid #EBEDEF' : 'none')};
`;
