import React, { useEffect, useState } from 'react';
import { Route, Switch, useParams, useRouteMatch, Router } from 'react-router-dom';
import { useStores } from 'store';
import NoResults from 'people/utils/userNoResults';
import { useIsMobile } from 'hooks';
import { Spacer } from 'people/main/body';
import styled from 'styled-components';
import { BountyModal } from 'people/main/BountyModal/BountyModal';
import { colors } from '../../config/colors';
import { bountyHeaderFilter, bountyHeaderLanguageFilter } from '../utils/filterValidation';
import WantedView from '../widgetViews/wantedView';
import history from '../../config/history';
import DeleteTicketModal from './deleteModal';

const Container = styled.div`
  display: flex;
  flex-flow: row wrap;
  gap: 1rem;
  flex: 1 1 100%;
`;

interface PanelProps {
  isMobile: boolean;
}
const Panel = styled.div<PanelProps>`
  position: relative;
  overflow: hidden;
  cursor: pointer;
  max-width: 300px;
  flex: 1 1 auto;
  background: #ffffff;
  color: #000000;
  padding: 20px;
  box-shadow: ${(p: any) => (p.isMobile ? 'none' : '0px 0px 6px rgb(0 0 0 / 7%)')};
  border-bottom: ${(p: any) => (p.isMobile ? '2px solid #EBEDEF' : 'none')};
`;

const UserTickets = () => {
  const color = colors['light'];
  const { personPubkey } = useParams<{ personPubkey: string }>();
  const { main, ui } = useStores();
  const isMobile = useIsMobile();
  const { path, url } = useRouteMatch();

  const [userTickets, setUserTickets] = useState<any>([]);
  const [checkboxIdToSelectedMap] = useState<any>({});
  const [currentItems] = useState<number>(10);
  const [deletePayload, setDeletePayload] = useState<object>({});
  const [showDeleteModal, setShowDeleteModal] = useState<boolean>(false);
  const closeModal = () => setShowDeleteModal(false);
  const showModal = () => setShowDeleteModal(true);
  const [loading, setIsLoading] = useState<boolean>(false);
  const data = {
    checkboxIdToSelectedMap
  };

  const activeList = userTickets.filter(({ body }: any) => {
    const value = { ...body };
    return (
      bountyHeaderFilter(data.checkboxIdToSelectedMap, value?.paid, !!value?.assignee) &&
      bountyHeaderLanguageFilter(value?.codingLanguage, {})
    );
  });

  async function getUserTickets() {
    setIsLoading(true);

    const tickets = await main.getPersonAssignedWanteds({}, personPubkey);

    setUserTickets(tickets);

    setIsLoading(false);
  }

  function onPanelClick(i: number) {
    history.push({
      pathname: `${url}/${i}`
    });
  }

  const deleteTicket = async (payload: any) => {
    const info = ui.meInfo as any;
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

  useEffect(() => {
    getUserTickets();
  }, []);

  const listItems =
    activeList && activeList.length ? (
      activeList.slice(0, currentItems).map((item: any, i: number) => {
        const { person, body } = item;

        // if this person has entries for this widget
        return (
          <Panel isMobile={isMobile} key={person?.owner_pubkey + i + body?.created}>
            <WantedView
              colors={color}
              showName
              onPanelClick={() => {
                onPanelClick(i);
              }}
              person={person}
              showModal={showModal}
              setDeletePayload={setDeletePayload}
              fromBountyPage={false}
              {...body}
              show={true}
            />
          </Panel>
        );
      })
    ) : (
      <NoResults loading={loading} />
    );

  return (
    <div data-testid="test">
      <Container>
        <Router history={history}>
          <Switch>
            <Route path={`${path}/:wantedId`}>
              <BountyModal basePath={url} />
            </Route>
          </Switch>
        </Router>
        {listItems}
        <Spacer key={'spacer2'} />
        {showDeleteModal && (
          <DeleteTicketModal closeModal={closeModal} confirmDelete={confirmDelete} />
        )}
      </Container>
    </div>
  );
};

export default UserTickets;
