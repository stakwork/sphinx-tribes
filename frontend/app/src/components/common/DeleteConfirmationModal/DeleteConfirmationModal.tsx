import React, { PropsWithChildren } from 'react';

import { Stack } from '@mui/system';
import { BaseModal } from '../BaseModal';
import IconButton from '../IconButton2';
import { useCreateModal } from '../useCreateModal';
import { icon } from './deleteIcon';

export type DeleteConfirmationModalProps = PropsWithChildren<{
  onClose: () => void;
  onDelete: () => void;
  onCancel?: () => void;
}>;
export const DeleteConfirmationModal = ({
  onClose,
  children,
  onCancel,
  onDelete
}: DeleteConfirmationModalProps) => {
  const closeHandler = () => {
    onClose();
    onCancel?.();
  };

  const deleteHandler = () => {
    onDelete();
    onClose();
  };
  return (
    <BaseModal backdrop="white" open onClose={closeHandler}>
      <Stack minWidth={350} p={4} alignItems="center" spacing={3}>
        {icon}
        {children}
        <Stack width="100%" direction="row" justifyContent="space-between">
          <IconButton width={120} height={44} color="white" text="Cancel" onClick={closeHandler} />
          <IconButton
            width={120}
            height={44}
            color="error"
            hovercolor={'#E96464'}
            activecolor={'#E15A5A'}
            shadowcolor={'#ED747480'}
            text="Delete"
            onClick={deleteHandler}
          />
        </Stack>
      </Stack>
    </BaseModal>
  );
};

export const useDeleteConfirmationModal = () => {
  const openModal = useCreateModal();

  const openDeleteConfirmation = (props: Omit<DeleteConfirmationModalProps, 'onClose'>) => {
    openModal({
      Component: DeleteConfirmationModal,
      props
    });
  };
  return { openDeleteConfirmation };
};
