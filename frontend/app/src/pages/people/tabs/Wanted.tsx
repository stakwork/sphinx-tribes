import { usePerson } from 'hooks';
import { observer } from 'mobx-react-lite';
import { BountyModal } from 'people/main/bountyModal';
import { widgetConfigs } from 'people/utils/Constants';
import NoneSpace from 'people/utils/NoneSpace';
import { PostBounty } from 'people/widgetViews/postBounty';
import WantedView from 'people/widgetViews/WantedView';
import React, { useEffect, useState } from 'react';
import { Route, Switch, useHistory, useRouteMatch, useParams } from 'react-router-dom';
import { useStores } from 'store';
import { PersonWanted } from 'store/main';
import styled from 'styled-components';
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

export const Wanted = observer(() => {
  const { ui, main } = useStores();
  const { person, canEdit } = usePerson(ui.selectedPerson);
  const { path, url } = useRouteMatch();
  const history = useHistory();
  const { personPubkey } = useParams<{ personPubkey: string }>();
  const [createdBounties, setCreatedBounties] = useState<PersonWanted[]>([]);

  async function getUserTickets() {
    const userBounties = await main.getPersonCreatedWanteds({}, personPubkey);
    setCreatedBounties(userBounties);
    await main.getPersonAssignedWanteds({}, personPubkey);
  }

  useEffect(() => {
    getUserTickets();
  }, []);

  if (!createdBounties?.length) {
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
      {createdBounties.map((w: any, i: any) => {
        if (w.body.owner_id === person?.owner_pubkey) {
          return (
            <Panel
              key={w.created}
              isMobile={false}
              onClick={() =>
                history.push({
                  pathname: `${url}/${w.body.id}/${i}`
                })
              }
            >
              <WantedView {...w.body} person={person} />
            </Panel>
          );
        }
      })}
    </Container>
  );
});
