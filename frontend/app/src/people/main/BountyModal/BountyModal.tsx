import { Modal } from 'components/common';
import { usePerson } from 'hooks';
import { widgetConfigs } from 'people/utils/Constants';
import React from 'react';
import { useHistory, useParams } from 'react-router-dom';
import { useStores } from 'store';
import { BountyModalProps } from 'people/interfaces';
import FocusedView from '../FocusView';

const config = widgetConfigs.wanted;
export const BountyModal = ({ basePath }: BountyModalProps) => {
  const history = useHistory();
  const { wantedId } = useParams<{ wantedId: string }>();

  const { ui } = useStores();
  const { canEdit, person } = usePerson(ui.selectedPerson);

  const wantedLength = person?.extras ? person?.extras.wanted?.length : 0;

  const changeWanted = (step: any) => {
    if (!wantedLength) return;
    const currentStep = Number(wantedId);
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
    history.push({
      pathname: basePath
    });
  };

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
        maxWidth: '90%',
        zIndex: 20
      }}
      nextArrow={() => changeWanted(1)}
      prevArrow={() => changeWanted(-1)}
      overlayClick={() => {
        onGoBack();
      }}
      bigClose={() => {
        onGoBack();
      }}
    >
      <FocusedView
        person={person}
        canEdit={canEdit}
        selectedIndex={Number(wantedId)}
        config={config}
        goBack={() => {
          onGoBack();
        }}
      />
    </Modal>
  );
};
