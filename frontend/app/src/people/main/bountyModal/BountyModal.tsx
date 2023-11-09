import { Modal } from 'components/common';
import { useIsMobile, usePerson } from 'hooks';
import { widgetConfigs } from 'people/utils/Constants';
import React, { useEffect, useState } from 'react';
import { useHistory, useParams } from 'react-router-dom';
import { useStores } from 'store';
import { BountyModalProps } from 'people/interfaces';
import { PersonBounty } from 'store/main';
import FocusedView from '../FocusView';

const config = widgetConfigs.wanted;
export const BountyModal = ({ basePath }: BountyModalProps) => {
  const history = useHistory();
  const { wantedId, wantedIndex, personPubkey } = useParams<{
    wantedId: string;
    wantedIndex: string;
    personPubkey: string;
  }>();

  const { ui, main } = useStores();
  const { person } = usePerson(ui.selectedPerson);
  const [bounty, setBounty] = useState<PersonBounty[]>([]);

  const onGoBack = async () => {
    await main.getPersonCreatedBounties({}, personPubkey);
    await main.getPersonAssignedBounties({}, personPubkey);

    ui.setBountyPerson(0);
    history.push({
      pathname: basePath
    });
  };

  useEffect(() => {
    async function getBounty() {
      /** check for the bounty length, else the request
       * will be made continously which will lead to an
       * infinite loop and crash the app
       */
      if (wantedId && !bounty.length) {
        const bounty = await main.getBountyById(Number(wantedId));
        setBounty(bounty);
      }
    }

    getBounty();
  }, [bounty, main, wantedId]);
  const isMobile = useIsMobile();

  if (isMobile) {
    return (
      <Modal visible={true} fill={true}>
        <FocusedView
          person={person}
          personBody={person}
          canEdit={false}
          selectedIndex={Number(wantedIndex)}
          config={config}
          goBack={onGoBack}
        />
      </Modal>
    );
  }

  return (
    <Modal
      visible={true}
      style={{
        background: 'rgba( 0 0 0 /75% )'
      }}
      envStyle={{
        maxHeight: '100vh',
        marginTop: 0,
        borderRadius: 0,
        background: '#fff',
        width: 'auto',
        minWidth: 500,
        maxWidth: '80%',
        zIndex: 20
      }}
      overlayClick={onGoBack}
      bigCloseImage={onGoBack}
      bigCloseImageStyle={{
        top: '18px',
        right: '-50px',
        borderRadius: '50%'
      }}
    >
      <FocusedView
        person={person}
        personBody={person}
        canEdit={false}
        selectedIndex={Number(wantedIndex)}
        config={config}
        bounty={bounty}
        goBack={() => {
          onGoBack();
        }}
        fromBountyPage={true}
      />
    </Modal>
  );
};
