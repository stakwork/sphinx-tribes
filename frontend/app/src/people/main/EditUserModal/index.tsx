import React, { useEffect } from 'react'
import { useIsMobile } from 'hooks'
import { EditUserDesctopView } from './EditUserDesctopView'
import { EditUserMobileView } from './EditUserMobileView'
import { useStores } from 'store'
import { observer } from 'mobx-react-lite'

export const EditUserModal = observer(() => {
  const isMobile = useIsMobile();
  const {modals} = useStores();
  useEffect(() => {
    console.log('EditUserModal')
    console.log(modals.userEditModal)

  }, [modals.userEditModal])
  if (!modals.userEditModal) {
    return null
  }

  return <>{isMobile ? <EditUserMobileView /> : <EditUserDesctopView />}</>
})