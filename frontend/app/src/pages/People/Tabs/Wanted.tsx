import { usePerson } from 'hooks';
import { observer } from 'mobx-react-lite';
import { BountyModal } from 'people/main/bountyModal';
import { widgetConfigs } from 'people/utils/Constants';
import NoneSpace from 'people/utils/NoneSpace';
import { PostBounty } from 'people/widgetViews/postBounty';
import WantedView from 'people/widgetViews/WantedView';
import React from 'react';
import { Route, Switch, useHistory, useRouteMatch } from 'react-router-dom';
import { useStores } from 'store';
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
  const { ui } = useStores();
  const { person, canEdit } = usePerson(ui.selectedPerson);
  const { path, url } = useRouteMatch();
  const history = useHistory();

  const fullSelectedWidgets = person?.extras?.wanted;

  if (!fullSelectedWidgets?.length) {
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
        <Route path={`${path}/:wantedId`}>
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
      {fullSelectedWidgets.map((w: any, i: number) => (
        <Panel
          key={w.created}
          isMobile={false}
          onClick={() =>
            history.push({
              pathname: `${url}/${i}`
            })
          }
        >
          <WantedView titleString={w.title} {...w} person={person} />
        </Panel>
      ))}
    </Container>
  );
});
