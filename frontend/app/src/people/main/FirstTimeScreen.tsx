import React from 'react';
import styled from 'styled-components';
import { useStores } from '../../store';
import { Modal } from '../../components/common';
import { firstScreenSchema } from '../../components/form/schema';
import AboutFocusedView from './AboutFocusView';

const B = styled.div`
  font-weight: bold;
`;

const Title = styled.div`
  font-size: 24px;
  font-style: normal;

  line-height: 30px;
  letter-spacing: 0em;
  text-align: center;
  margin-bottom: 20px;
`;

const SubTitle = styled.div`
  font-family: Roboto;
  font-size: 15px;
  font-style: normal;
  font-weight: 400;
  line-height: 20px;
  letter-spacing: 0em;
  text-align: center;
`;
// this is where we see others posts (etc) and edit our own
export default function FirstTimeScreen() {
  const { ui, main } = useStores();

  const formHeader = (
    <div style={{ marginTop: 60 }}>
      <Title>
        <B>Hi {ui.meInfo?.owner_alias},</B>
        <div>thank you for joining.</div>
      </Title>
      <SubTitle>Please enter some basic info about yourself and create a public profile.</SubTitle>
    </div>
  );

  return (
    <Modal
      visible={true}
      envStyle={{
        height: 'fit-content',
        borderRadius: 8,
        overflow: 'hidden',
        width: '100%',
        maxWidth: 600
      }}
    >
      <div style={{ height: '100%', padding: 20, paddingTop: 0, width: '100%' }}>
        <AboutFocusedView
          formHeader={formHeader}
          isFirstTimeScreen={true}
          buttonsOnBottom={true}
          person={ui.meInfo}
          canEdit={true}
          manualGoBackOnly={true}
          goBack={() => {
            ui.setMeInfo(null);
            main.getPeople();
          }}
          selectedIndex={-1}
          config={{
            label: 'About',
            name: 'about',
            single: true,
            skipEditLayer: true,
            submitText: 'Submit',
            schema: firstScreenSchema
          }}
          onSuccess={() => {
            console.log('success');
          }}
        />
      </div>
    </Modal>
  );
}
