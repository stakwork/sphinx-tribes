
import { Modal } from 'components/common';
import React from 'react';
import FocusedView from '../focusView';
import { formConfig } from './config';
import { useUserEdit } from './useEditUser';
import { observer } from 'mobx-react-lite';

export const EditUserMobileView = observer(() => {
  const {canEdit, closeHandler, person, modals} = useUserEdit()

      return (
  <Modal fill visible={modals.userEditModal}>
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
})
