import { Modal } from 'components/common';
import { usePerson } from 'hooks';
import { widgetConfigs } from 'people/utils/Constants';
import React, { useEffect, useState } from 'react';
import { useHistory, useParams } from 'react-router-dom';
import { useStores } from 'store';
import { BountyModalProps } from 'people/interfaces';
import { PersonWanted } from 'store/main';
import FocusedView from '../FocusView';

const config = widgetConfigs.wanted;
export const BountyModal = ({ basePath }: BountyModalProps) => {
  const history = useHistory();
  const { wantedId, wantedIndex } = useParams<{ wantedId: string; wantedIndex: string }>();

  const { ui, main } = useStores();
  const { canEdit, person } = usePerson(ui.selectedPerson);
  const [bounty, setBounty] = useState<PersonWanted[]>([]);

  const wantedLength = person?.extras ? person?.extras.wanted?.length : 0;

  const changeWanted = (step: any) => {
    if (!wantedLength) return;
    const currentStep = Number(wantedIndex);
    const newStep = currentStep + step;

    if (step === 1) {
      if (newStep < wantedLength) {
        history.replace({
          pathname: `${basePath}/${newStep}`
        });
      }
    }
    if (step === -1) {
      if (newStep >= 0) {
        history.replace({
          pathname: `${basePath}/${newStep}`
        });
      }
    }
  };

  const onGoBack = () => {
    ui.setBountyPerson(0);
    history.push({
      pathname: basePath
    });
  };

  useEffect(() => {
    async function getBounty() {
      if (wantedId && !bounty.length) {
        const bounty = await main.getWantedById(Number(wantedId));
        setBounty(bounty);
      }
    }

    getBounty();
  }, [bounty, main, wantedId]);

  return (
    <Modal
      visible={true}
      style={{
        minHeight: '100%',
        height: 'auto'
      }}
      envStyle={{
        marginTop: 0,
        borderRadius: 0,
        background: '#fff',
        height: '100%',
        width: 'auto',
        minWidth: 500,
        maxWidth: '80%',
        zIndex: 20
      }}
      overlayClick={onGoBack}
      bigCloseImage={onGoBack}
    >
      <FocusedView
        person={person}
        canEdit={ui.bountyPerson ? person?.id === ui.bountyPerson : canEdit}
        selectedIndex={Number(wantedIndex)}
        config={config}
        bounty={bounty}
        goBack={() => {
          onGoBack();
        }}
      />
    </Modal>
  );
};
