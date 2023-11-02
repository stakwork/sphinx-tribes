import { Modal } from 'components/common';
import React from 'react';
import { observer } from 'mobx-react-lite';
import AboutFocusView from '../AboutFocusView';
import { useUserEdit } from './useEditUser';
import { formConfig } from './config';

export const EditUserDesktopView = observer(() => {
  const { canEdit, closeHandler, person, modals } = useUserEdit();

  return (
    <Modal
      visible={modals.userEditModal}
      style={{
        height: '100%'
      }}
      envStyle={{
        borderRadius: 0,
        background: '#fff',
        height: '100%',
        width: '60%',
        minWidth: 500,
        maxWidth: 602,
        zIndex: 20
      }}
      overlayClick={closeHandler}
      bigClose={closeHandler}
    >
      <AboutFocusView
        person={person}
        canEdit={canEdit}
        selectedIndex={0}
        config={formConfig}
        onSuccess={closeHandler}
        goBack={closeHandler}
      />
    </Modal>
  );
});
