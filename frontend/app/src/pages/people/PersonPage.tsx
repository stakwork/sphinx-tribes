import { Modal } from 'components/common';
import { useIsMobile, usePerson } from 'hooks';
import { observer } from 'mobx-react-lite';
import { UserInfo } from 'people/userInfo';
import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { useStores } from 'store';
import styled from 'styled-components';
import { PeopleList } from './peopleList';
import { TabsPages } from './tabs';

// eslint-disable-next-line @typescript-eslint/no-inferrable-types
const getHtml = (owner_pubkey: string = '', img: string = '') => `
<sphinx-widget pubkey=${owner_pubkey}
  amount="500"
  title="Support Me"
  subtitle="Because I'm awesome"
  buttonlabel="Donate"
  defaultinterval="weekly"
  imgurl="${img || 'https://i.scdn.co/image/28747994a80c78bc2824c2561d101db405926a37'}"
  ></sphinx-widget>`;

const Content = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  align-items: center;
  color: #000000;
  padding-bottom: 60px;
  background: #f0f1f3;
  .desktop {
    position: relative;
    display: flex;
    width: 100%;
    height: 100%;
  }
`;
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
  padding-bottom: ${(p: any) => (p.isMobile ? 80 : 0)};
`;
export const PersonPage = observer(() => {
  const { main, ui } = useStores();
  const isMobile = useIsMobile();
  const { personPubkey } = useParams<{ personPubkey: string }>();
  const [showSupport, setShowSupport] = useState(false);

  useEffect(() => {
    (async () => {
      const p = await main.getPersonByPubkey(personPubkey);
      ui.setSelectedPerson(p?.id);
      ui.setSelectingPerson(p?.id);
    })();
  }, [main, personPubkey, ui]);

  const personId = ui.selectedPerson;
  const { person, canEdit } = usePerson(personId);

  return (
    <Content>
      {!isMobile && (
        <div className="desktop">
          {!canEdit && <PeopleList />}
          <UserInfo setShowSupport={setShowSupport} />
          <TabsPages />
        </div>
      )}
      {isMobile && (
        <div
          style={{
            display: 'flex',
            flexDirection: 'column',
            width: '100%',
            overflow: 'auto',
            height: '100%'
          }}
        >
          <Panel isMobile={isMobile}>
            <UserInfo setShowSupport={setShowSupport} />
            <TabsPages />
          </Panel>
        </div>
      )}
      <Modal
        visible={showSupport}
        close={() => setShowSupport(false)}
        style={{
          height: '100%'
        }}
        envStyle={{
          marginTop: isMobile || canEdit ? 64 : 123,
          borderRadius: 0
        }}
      >
        <>
          {person?.owner_pubkey && person?.img && (
            <div
              dangerouslySetInnerHTML={{
                __html: getHtml(person.owner_pubkey, person?.img)
              }}
            />
          )}
        </>
      </Modal>
    </Content>
  );
});
