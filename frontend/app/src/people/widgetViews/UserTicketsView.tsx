import React, { useEffect, useState } from 'react';
import { Route, Switch, useParams, useRouteMatch, Router } from 'react-router-dom';
import { useStores } from 'store';
import NoResults from 'people/utils/UserNoResults';
import { useIsMobile } from 'hooks';
import { Spacer } from 'people/main/Body';
import styled from 'styled-components';
import { BountyModal } from 'people/main/bountyModal/BountyModal';
import PageLoadSpinner from 'people/utils/PageLoadSpinner';
import { Person } from 'store/main';
import history from '../../config/history';
import { colors } from '../../config/colors';
import WantedView from './WantedView';
import DeleteTicketModal from './DeleteModal';
import { LoadMoreContainer } from './WidgetSwitchViewer';

const Container = styled.div`
  display: flex;
  flex-flow: row wrap;
  gap: 1rem;
  align-items: center;
  flex: 1 1 100%;
`;

interface PanelProps {
  isMobile: boolean;
}
const Panel = styled.a<PanelProps>`
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

  &:hover {
    text-decoration: none !important;
  }
`;

const UserTickets = () => {
  const color = colors['light'];
  const { personPubkey } = useParams<{ personPubkey: string }>();
  const { main, ui } = useStores();
  const isMobile = useIsMobile();
  const { path, url } = useRouteMatch();

  const [userTickets, setUserTickets] = useState<any>([]);
  const [currentItems] = useState<number>(20);
  const [visibleItems, setVisibleItems] = useState<number>(currentItems);

  const loadMore = () => {
    setVisibleItems((prevVisibleItems: number) => prevVisibleItems + currentItems);
  };

  const [deletePayload, setDeletePayload] = useState<object>({});
  const [showDeleteModal, setShowDeleteModal] = useState<boolean>(false);
  const closeModal = () => setShowDeleteModal(false);
  const showModal = () => setShowDeleteModal(true);
  const [loading, setIsLoading] = useState<boolean>(false);
  const [bountyOwner, setBountyOwner] = useState<Person>();

  function onPanelClick(id: number, index: number) {
    history.push({
      pathname: `${url}/${id}/${index}`
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
    async function getUserTickets() {
      setIsLoading(true);
      const tickets = await main.getPersonAssignedBounties({}, personPubkey);
      setUserTickets(tickets);
      setIsLoading(false);
    }

    getUserTickets();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [main.getPersonAssignedBounties, personPubkey]);

  const listItems =
    userTickets && userTickets.length ? (
      userTickets.slice(0, visibleItems).map((item: any, i: number) => {
        const person = main.people.find((p: any) => p.owner_pubkey === item.body.owner_id);
        const body = { ...item.body };
        // if this person has entries for this widget
        return (
          <Panel href={`${url}/${body.id}/${i}`} isMobile={isMobile} key={i + body?.created}>
            <WantedView
              colors={color}
              showName
              onPanelClick={(e: any) => {
                e.preventDefault();
                onPanelClick(body.id, i);
                ui.setBountyPerson(person?.id);
                setBountyOwner(person);
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
      <NoResults />
    );

  return (
    <Container data-testid="test">
      <PageLoadSpinner show={loading} />
      <Router history={history}>
        <Switch>
          <Route path={`${path}/:wantedId/:wantedIndex`}>
            <BountyModal fromPage={'usertickets'} bountyOwner={bountyOwner} basePath={url} />
          </Route>
        </Switch>
      </Router>
      {listItems}
      {userTickets.length > visibleItems && (
        <LoadMoreContainer
          color={color}
          style={{
            width: '100%',
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            margin: '20px 0'
          }}
        >
          <div className="LoadMoreButton" onClick={loadMore}>
            Load More
          </div>
        </LoadMoreContainer>
      )}
      <Spacer key={'spacer2'} />
      {showDeleteModal && (
        <DeleteTicketModal closeModal={closeModal} confirmDelete={confirmDelete} />
      )}
    </Container>
  );
};

export default UserTickets;
