
import { Modal } from 'components/common';
import React from 'react'
import { useUserEdit } from './useEditUser';
import FocusedView from '../focusView';
import { formConfig } from './config';

export const EditUserDesctopView = () => {
  const {canEdit, closeHandler, person, showModal} = useUserEdit();
  return (
    <Modal
      visible={showModal}
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
        zIndex: 20, 
      }}
      overlayClick={closeHandler}
      bigClose={closeHandler}
    >
      <FocusedView
        person={person}
        canEdit={canEdit}
        selectedIndex={0}
        config={formConfig}
        onSuccess={closeHandler}
        goBack={closeHandler}
      />
    </Modal>
  )
}
