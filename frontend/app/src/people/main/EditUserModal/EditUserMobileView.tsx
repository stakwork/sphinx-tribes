
import { Modal } from 'components/common';
import React from 'react';
import FocusedView from '../focusView';
import { formConfig } from './config';
import { useUserEdit } from './useEditUser';

export const EditUserMobileView = () => {

  const {canEdit, closeHandler, person, showModal} = useUserEdit()
  return (
    <Modal fill visible={showModal}>
    <FocusedView
      person={person}
      canEdit={canEdit}
      selectedIndex={1}
      config={formConfig}
      onSuccess={closeHandler}
      goBack={closeHandler}
    />
    </Modal>
  )
}
