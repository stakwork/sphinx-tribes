/* eslint-disable func-style */
import React from 'react';
import styled from 'styled-components';
import { observer } from 'mobx-react-lite';
import { RenderWidgetsProps } from 'people/interfaces';
import { useStores } from '../../store';

import { meSchema } from '../../components/form/schema';
import { useIsMobile } from '../../hooks';
import Badges from '../utils/Badges';
import { widgetConfigs } from '../utils/Constants';
import NoneSpace from '../utils/NoneSpace';
import OfferView from './OfferView';
import { AboutView } from './AboutView';
import { PostBounty } from './postBounty';
import UserTickets from './UserTicketsView';
import OrganizationView from './OrganizationView';

interface PanelProps {
  isMobile: boolean;
}

const Panel = styled.div<PanelProps>`
  position: relative;
  background: #ffffff;
  color: #000000;
  padding: 20px;
  box-shadow: ${(p: any) => (p.isMobile ? 'none' : '0px 0px 6px rgb(0 0 0 / 7%)')};
  border-bottom: ${(p: any) => (p.isMobile ? '2px solid #EBEDEF' : 'none')};
`;

function RenderWidgets({ widget }: RenderWidgetsProps) {
  const { main, ui } = useStores();
  const { meInfo } = ui || {};
  const personId = ui.selectedPerson;

  let person: any =
    main.people && main.people.length && main.people.find((f: any) => f.id === personId);

  // if i select myself, fill person with meInfo
  if (personId === ui.meInfo?.id) {
    person = {
      ...ui.meInfo
    };
  }

  const { id, extras } = person || {};

  let { description } = person || {};

  // backend is adding 'description' to empty descriptions, short term fix
  if (description === 'description') description = '';

  const canEdit = id === meInfo?.id;
  const isMobile = useIsMobile();

  const selectedWidget = widget;

  let widgetSchemas: any = meSchema.find((f: any) => f.name === 'extras');
  if (widgetSchemas && widgetSchemas.extras) {
    widgetSchemas = widgetSchemas && widgetSchemas.extras;
  }

  const fullSelectedWidget: any = extras && selectedWidget ? extras[selectedWidget] : null;

  // we do this because sometimes the widgets are empty arrays
  const filteredExtras = extras && { ...extras };

  if (filteredExtras) {
    const emptyArrayKeys = [''];

    Object.keys(filteredExtras).forEach((name: any) => {
      const p = extras && extras[name];
      if (Array.isArray(p) && !p.length) {
        emptyArrayKeys.push(name);
      }
      const thisSchema = widgetSchemas && widgetSchemas.find((e: any) => e.name === name);
      if (filteredExtras && thisSchema && thisSchema.single) {
        delete filteredExtras[name];
      }
    });

    emptyArrayKeys.forEach((e: any) => {
      if (filteredExtras && e) delete filteredExtras[e];
    });
  }

  const tabs = widgetConfigs;

  function renderWidgets() {
    if (!selectedWidget) {
      return <div style={{ height: 200 }} />;
    }

    if (selectedWidget === 'badges') {
      return <Badges person={person} />;
    }

    const widgetSchema: any =
      (widgetSchemas && widgetSchemas.find((f: any) => f.name === selectedWidget)) || {};
    const { single } = widgetSchema;

    function wrapIt(child: any) {
      if (single) {
        return <Panel isMobile={isMobile}>{child}</Panel>;
      }

      const elementArray: any = [];

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

      fullSelectedWidget &&
        fullSelectedWidget.forEach((s: any, i: number) => {
          if (!canEdit && 'show' in s && s.show === false) {
            // skip hidden items
            return;
          }

          const conditionalStyles =
            !isMobile && s?.paid
              ? {
                  border: '1px solid #dde1e5',
                  boxShadow: 'none'
                }
              : {};

          elementArray.push(
            <Panel
              isMobile={isMobile}
              key={i}
              style={{
                ...panelStyles,
                ...conditionalStyles,
                cursor: 'pointer',
                padding: 0,
                overflow: 'hidden'
              }}
            >
              {React.cloneElement(child, { ...s })}
            </Panel>
          );
        });
      const noneKey = canEdit ? 'me' : 'otherUser';
      const noneSpaceProps = tabs[selectedWidget]?.noneSpace[noneKey];

      const panels: any = elementArray.length ? (
        <div style={{ width: '100%', display: 'flex', flexDirection: 'column' }}>
          <div style={{ width: '100%', display: 'flex', flexDirection: 'row', flexWrap: 'wrap' }}>
            {elementArray}
          </div>
        </div>
      ) : (
        <NoneSpace
          small
          style={{
            margin: 'auto'
          }}
          Button={
            canEdit && (
              <PostBounty
                title={noneSpaceProps.buttonText}
                buttonProps={{
                  leadingIcon: noneSpaceProps.buttonIcon,
                  color: 'secondary'
                }}
                widget={selectedWidget}
              />
            )
          }
          {...tabs[selectedWidget]?.noneSpace[noneKey]}
        />
      );

      return <>{panels}</>;
    }

    switch (selectedWidget) {
      case 'organizations':
        return <OrganizationView person={person} />;
      case 'badges':
        return <Badges person={person} />;
      case 'about':
        return (
          <Panel isMobile={isMobile}>
            <AboutView {...person} />
          </Panel>
        );
      case 'usertickets':
        return <UserTickets />;
      case 'offer':
        return wrapIt(<OfferView {...fullSelectedWidget} person={person} />);
      default:
        return null;
    }
  }

  return renderWidgets();
}
export default observer(RenderWidgets);
