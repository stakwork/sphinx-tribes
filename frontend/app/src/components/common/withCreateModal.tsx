/* eslint-disable react/display-name */
import React, { FC, MutableRefObject, useCallback, useEffect, useRef, useState } from 'react';

import { CreateModalContext, CreateModalPortal, CreateModalProps } from './useCreateModal';

const ModalInjector = (props: {
  createCb: MutableRefObject<(props: CreateModalProps) => void>;
}) => {
  const [modals, setModals] = useState<CreateModalProps[]>([]);

  const closeDialog = useCallback(() => {
    setModals((modals: CreateModalProps[]) => {
      const latestDialog = modals.pop();
      if (!latestDialog) {
        return modals;
      }
      if (latestDialog.props.onClose) {
        latestDialog.props.onClose();
      }
      return [...modals];
    });
  }, []);

  const createModal = useCallback((option: CreateModalProps) => {
    setModals((m: CreateModalProps[]) => [...m, option]);
  }, []);

  useEffect(
    function applyInjectCb() {
      props.createCb.current = createModal;
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [createModal]
  );
  return (
    <>
      {modals.map((modal: CreateModalProps) =>
        CreateModalPortal(modal.Component, { ...modal.props, onClose: closeDialog })
      )}
    </>
  );
};

// eslint-disable-next-line @typescript-eslint/ban-types
export function withCreateModal<T extends Object>(Component: FC<T>) {
  return function (props: T) {
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    const injectModalCb = useRef<(props: CreateModalProps) => void>(() => {});

    return (
      <CreateModalContext.Provider value={injectModalCb}>
        <ModalInjector createCb={injectModalCb} />
        <Component {...props} />
      </CreateModalContext.Provider>
    );
  };
}
