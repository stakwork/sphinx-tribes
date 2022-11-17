import React, { useState } from 'react';
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
import { uiStore } from '../../store/ui';
import DeleteTicketModal from './deleteModal';

export default function WidgetSwitchViewer(props) {
  const { main } = useStores();
  const isMobile = useIsMobile();
  const [showDeleteModal, setShowDeleteModal] = useState<boolean>(false);
  const [deletePayload, setDeletePayload] = useState<object>({});
  const closeModal = () => setShowDeleteModal(false);
  const showModal = () => setShowDeleteModal(true);
  const [currentItems, setCurrentItems] = useState<number>(10);

  const panelStyles = isMobile
    ? {
        minHeight: 132
      }
    : {
        minWidth: '1100px',
        maxWidth: '1100px',
        marginBottom: 20,
        borderRadius: '10px',
        display: 'flex',
        justifyContent: 'center',
      };

  return useObserver(() => {
    const { peoplePosts, peopleWanteds, peopleOffers } = main;

    const { selectedWidget, onPanelClick } = props;

    if (!selectedWidget) {
      return <div style={{ height: 200 }} />;
    }

    const listSource = {
      post: peoplePosts,
      wanted: peopleWanteds,
      offer: peopleOffers
    };

    const activeList = listSource[selectedWidget];

    const foundDynamicSchema = widgetConfigs[selectedWidget]?.schema?.find((f) => f.dynamicSchemas);
    // if dynamic schema, get all those fields
    if (foundDynamicSchema) {
      const dynamicFields: any = [];
      foundDynamicSchema.dynamicSchemas?.forEach((ds) => {
        ds.forEach((f) => {
          if (!dynamicFields.includes(f.name)) dynamicFields.push(f.name);
        });
      });
    }

    const deleteTicket = async (payload: any) => {
      const info = uiStore.meInfo as any;
      const URL = info.url.startsWith('http') ? info.url : `https://${info.url}`;
      try {
        await fetch(`${URL}/delete_ticket`, {
          method: 'POST',
          body: JSON.stringify(payload),
          headers: {
            'x-jwt': info.jwt,
            'Content-Type': 'application/json'
          }
        });
      } catch (error) {
        console.log(error);
      }
    };

    const confirmDelete = async () => {
      try {
        if (deletePayload) {
          await deleteTicket(deletePayload);
        }
      } catch (error) {
        console.log(error);
      }
      closeModal();
    };

    const listItems =
      activeList && activeList.length ? (
        activeList.slice(0, currentItems).map((item, i) => {
          const { person, body } = item;

          const conditionalStyles = body?.paid
            ? {
                border: isMobile ? '2px 0 0 0 solid #dde1e5' : '',
                boxShadow: 'none'
              }
            : {};

          // if this person has entries for this widget
          return (
            <Panel
              isMobile={isMobile}
              key={person?.owner_pubkey + i + body?.created}
              // onClick={() => {
              //   if (onPanelClick) onPanelClick(person, body);
              // }}
              style={{
                ...panelStyles,
                ...conditionalStyles,
                cursor: 'pointer',
                padding: 0,
                overflow: 'hidden',
                background: 'transparent',
                minHeight: !isMobile ? '160px' : '',
                boxShadow: 'none',
              }}
            >
              {selectedWidget === 'post' ? (
                <PostView
                  showName
                  key={`${i + person.owner_pubkey}pview`}
                  person={person}
                  {...body}
                />
              ) : selectedWidget === 'offer' ? (
                <OfferView
                  showName
                  key={`${i + person.owner_pubkey}oview`}
                  person={person}
                  {...body}
                />
              ) : selectedWidget === 'wanted' ? (
                <WantedView
                  showName
                  onPanelClick={() => {
                      if (onPanelClick) onPanelClick(person, body);
                    }}
                  key={`${i + person.owner_pubkey}wview`}
                  person={person}
                  showModal={showModal}
                  setDeletePayload={setDeletePayload}
                  fromBountyPage={props.fromBountyPage}
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

        {showDeleteModal && (
          <DeleteTicketModal closeModal={closeModal} confirmDelete={confirmDelete} />
        )}
        {activeList?.length > currentItems && (
          <LoadMoreContainer
            style={{
              width: '100%',
              display: 'flex',
              justifyContent: 'center',
              alignItems: 'center',
            }}
          >
            <div
              className="LoadMoreButton"
              onClick={() => {
                setCurrentItems(currentItems + 10);
              }}
            >
              Load More
            </div>
          </LoadMoreContainer>
        )}
        <Spacer key={'spacer'} />
      </>
    );
  });
}

interface PanelProps {
  isMobile?: boolean;
}

const Panel = styled.div<PanelProps>`
  // position: ;
  background: #ffffff;
  color: #000000;
  padding: 20px;
  box-shadow: ${(p) => (p.isMobile ? 'none' : '0px 0px 6px rgb(0 0 0 / 7%)')};
  border-bottom: ${(p) => (p.isMobile ? '2px solid #EBEDEF' : 'none')};
`;

const LoadMoreContainer = styled.div<PanelProps>`
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  .LoadMoreButton {
    width: 166px;
    height: 48px;
    display: flex;
    justify-content: center;
    align-items: center;
    color: #3c3f41;
    border: 1px solid #dde1e5;
    border-radius: 30px;
    background: #ffffff;
    font-family: 'Barlow';
    font-style: normal;
    font-weight: 500;
    font-size: 14px;
    line-height: 17px;
    cursor: pointer;
    user-select: none;
    :hover {
      border: 1px solid #b0b7bc;
    }
    :active {
      border: 1px solid #8e969c;
    }
  }
`;
