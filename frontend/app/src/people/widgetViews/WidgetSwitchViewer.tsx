import React, { useState } from 'react';
import styled from 'styled-components';
import { observer } from 'mobx-react-lite';
import { useIsMobile } from 'hooks/uiHooks';
import { queryLimit } from 'store/main';
import { Spacer } from '../main/Body';
import NoResults from '../utils/NoResults';
import { uiStore } from '../../store/ui';
import { colors } from '../../config/colors';
import { useStores } from '../../store';
import { widgetConfigs } from '../utils/Constants';
import OfferView from './OfferView';
import WantedView from './WantedView';
import PostView from './PostView';
import DeleteTicketModal from './DeleteModal';

interface PanelProps {
  isMobile?: boolean;
  color?: any;
  isAssignee?: boolean;
}

const Panel = styled.a<PanelProps>`
  margin-top: 5px;
  background: ${(p: any) => p.color && p.color.pureWhite};
  color: ${(p: any) => p.color && p.color.pureBlack};
  padding: 20px;
  border-bottom: ${(p: any) => (p.isMobile ? `2px solid ${p.color.grayish.G700}` : 'none')};
  :hover {
    box-shadow: ${(p: any) =>
      p.isAssignee ? `0px 1px 6px ${p.color.black100}` : 'none'} !important;
    text-decoration: none !important;
  }
  :active {
    box-shadow: none !important;
  }
`;

export const LoadMoreContainer = styled.div<PanelProps>`
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
    color: ${(p: any) => p.color && p.color.grayish.G10};
    border: 1px solid ${(p: any) => p.color && p.color.grayish.G600};
    border-radius: 30px;
    background: ${(p: any) => p.color && p.color.pureWhite};
    font-family: 'Barlow';
    font-style: normal;
    font-weight: 500;
    font-size: 14px;
    line-height: 17px;
    cursor: pointer;
    user-select: none;
    :hover {
      border: 1px solid ${(p: any) => p.color && p.color.grayish.G300};
    }
    :active {
      border: 1px solid ${(p: any) => p.color && p.color.grayish.G100};
    }
  }
`;

function WidgetSwitchViewer(props: any) {
  const color = colors['light'];
  const { main } = useStores();
  const isMobile = useIsMobile();
  const [showDeleteModal, setShowDeleteModal] = useState<boolean>(false);
  const [deletePayload, setDeletePayload] = useState<object>({});
  const closeModal = () => setShowDeleteModal(false);
  const showModal = () => setShowDeleteModal(true);
  const { currentItems, setCurrentItems, totalBounties, page: propsPage, setPage } = props;

  const items = currentItems ?? 0;
  const bountiesTotal = totalBounties ?? 0;
  const page = propsPage ?? 0;

  const panelStyles = isMobile
    ? {
        minHeight: 132
      }
    : {
        minWidth: '1100px',
        maxWidth: '1100px',
        marginBottom: 16,
        borderRadius: '10px',
        display: 'flex',
        justifyContent: 'center'
      };

  const { peoplePosts, peopleBounties, peopleOffers } = main;

  const { selectedWidget, onPanelClick } = props;

  if (!selectedWidget) {
    return <div style={{ height: 200 }} />;
  }

  const listSource = {
    post: peoplePosts,
    wanted: peopleBounties,
    offer: peopleOffers
  };

  const activeList = [...listSource[selectedWidget]].filter(({ body }: any) => {
    const value = { ...body };
    return value;
  });

  const foundDynamicSchema = widgetConfigs[selectedWidget]?.schema?.find(
    (f: any) => f.dynamicSchemas
  );
  // if dynamic schema, get all those fields
  if (foundDynamicSchema) {
    const dynamicFields: any = [];
    foundDynamicSchema.dynamicSchemas?.forEach((ds: any) => {
      ds.forEach((f: any) => {
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

  const nextBounties = async () => {
    const currentPage = page + 1;
    if (setPage) {
      setPage(currentPage);
    }

    if (setCurrentItems) {
      setCurrentItems(currentItems + queryLimit);
    }
    await main.getPeopleBounties({
      limit: queryLimit,
      page: currentPage,
      ...props.checkboxIdToSelectedMap
    });
  };

  const listItems =
    activeList && activeList.length ? (
      activeList.slice(0, currentItems).map((item: any, i: number) => {
        const { person, body, organization } = item;
        person.img = person.img || main.getUserAvatarPlaceholder(person.owner_pubkey);
        const conditionalStyles = body?.paid
          ? {
              border: isMobile ? `2px 0 0 0 solid ${color.grayish.G600}` : '',
              boxShadow: 'none'
            }
          : {};

        // if this person has entries for this widget
        return (
          <Panel
            color={color}
            isMobile={isMobile}
            key={person?.owner_pubkey + i + body?.created}
            isAssignee={!!body.assignee}
            style={{
              ...panelStyles,
              ...conditionalStyles,
              cursor: 'pointer',
              padding: 0,
              overflow: 'hidden',
              background: 'transparent',
              minHeight: !isMobile ? '160px' : '',
              maxHeight: 'auto',
              boxShadow: 'none'
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
                person={person}
                showModal={showModal}
                setDeletePayload={setDeletePayload}
                fromBountyPage={props.fromBountyPage}
                {...body}
                {...organization}
              />
            ) : null}
          </Panel>
        );
      })
    ) : (
      <NoResults />
    );
  const showLoadMore = bountiesTotal > items && activeList.length >= queryLimit;
  return (
    <>
      {listItems}
      <Spacer key={'spacer2'} />

      {showDeleteModal && (
        <DeleteTicketModal closeModal={closeModal} confirmDelete={confirmDelete} />
      )}
      {showLoadMore && (
        <LoadMoreContainer
          color={color}
          style={{
            width: '100%',
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center'
          }}
        >
          <div className="LoadMoreButton" onClick={() => nextBounties()}>
            Load More
          </div>
        </LoadMoreContainer>
      )}
      <Spacer key={'spacer'} />
    </>
  );
}
export default observer(WidgetSwitchViewer);
