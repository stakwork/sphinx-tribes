import { usePerson } from 'hooks';
import { observer } from 'mobx-react-lite';
import { BountyModal } from 'people/main/bountyModal';
import { widgetConfigs } from 'people/utils/Constants';
import NoneSpace from 'people/utils/NoneSpace';
import { PostBounty } from 'people/widgetViews/postBounty';
import WantedView from 'people/widgetViews/WantedView';
import PageLoadSpinner from 'people/utils/PageLoadSpinner';
import React, { useEffect, useState } from 'react';
import { Route, Switch, useHistory, useRouteMatch, useParams } from 'react-router-dom';
import { useStores } from 'store';
import styled from 'styled-components';
import { LoadMoreContainer } from '../../../people/widgetViews/WidgetSwitchViewer';
import { colors } from '../../../config/colors';
const config = widgetConfigs.wanted;

const Container = styled.div`
  display: flex;
  flex-flow: row wrap;
  gap: 1rem;
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

export const Wanted = observer(() => {
  const { ui, main } = useStores();
  const { person, canEdit } = usePerson(ui.selectedPerson);
  const { path, url } = useRouteMatch();
  const history = useHistory();
  const { personPubkey } = useParams<{ personPubkey: string }>();
  const [loading, setIsLoading] = useState<boolean>(false);
  const [currentItems] = useState<number>(20);
  const [visibleItems, setVisibleItems] = useState<number>(currentItems);
  const loadMore = () => {
    setVisibleItems((prevVisibleItems: number) => prevVisibleItems + currentItems);
  };

  async function getUserTickets() {
    setIsLoading(true);
    const tickets = await main.getPersonCreatedBounties({}, personPubkey);
    main.setCreatedBounties(tickets);
    await main.getPersonAssignedBounties({ sortBy: 'paid' }, personPubkey);
    setIsLoading(false);
  }

  useEffect(() => {
    getUserTickets();
  }, [main]);

  if (!main.createdBounties?.length) {
    return (
      <NoneSpace
        style={{
          margin: 'auto'
        }}
        small
        Button={
          canEdit && (
            <PostBounty
              title={config.noneSpace.me.buttonText}
              buttonProps={{
                leadingIcon: config.noneSpace.me.buttonIcon,
                color: 'secondary'
              }}
              widget={'wanted'}
              onSucces={() => {
                history.goBack();
                window.location.reload();
              }}
              onGoBack={() => {
                history.goBack();
              }}
            />
          )
        }
        {...(canEdit ? config.noneSpace.me : config.noneSpace.otherUser)}
      />
    );
  }
  return (
    <Container>
      <PageLoadSpinner show={loading} />
      <Switch>
        <Route path={`${path}/:wantedId/:wantedIndex`}>
          <BountyModal basePath={url} />
        </Route>
      </Switch>
      <div
        style={{
          width: '100%',
          display: 'flex',
          justifyContent: 'flex-end',
          paddingBottom: '16px'
        }}
      >
        {canEdit && <PostBounty widget="wanted" />}
      </div>
      {main.createdBounties.slice(0, visibleItems).map((w: any, i: any) => {
        if (w.body.owner_id === person?.owner_pubkey) {
          return (
            <Panel
              href={`${url}/${w.body.id}/${i}`}
              key={w.body.id}
              isMobile={false}
              onClick={(e: any) => {
                e.preventDefault();
                ui.setBountyPerson(person?.id);
                history.push({
                  pathname: `${url}/${w.body.id}/${i}`
                });
              }}
            >
              <WantedView {...w.body} person={person} />
            </Panel>
          );
        }
        return null;
      })}
      {main.createdBounties.length > visibleItems && (
        <LoadMoreContainer
          color={colors['light']}
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
    </Container>
  );
});
