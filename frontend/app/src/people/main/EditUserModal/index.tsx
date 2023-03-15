import React from 'react'
import { useIsMobile } from 'hooks'
import { EditUserDesctopView } from './EditUserDesctopView'
import { EditUserMobileView } from './EditUserMobileView'
import { useModalsVisibility } from 'store/modals'

export const EditUserModal = () => {
  const isMobile = useIsMobile();
  const modals = useModalsVisibility();
  console.log('EditUserModal')
  console.log(modals.userEditModal)
  if (!modals.userEditModal) {
    return null
  }

  return <>{isMobile ? <EditUserMobileView /> : <EditUserDesctopView />}</>
}