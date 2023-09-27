/* eslint-disable no-use-before-define */
import React, { ComponentProps, MutableRefObject, createContext, useContext } from 'react';
import { createPortal } from 'react-dom';

export function getEmptyContext<T>(): React.Context<T> {
  return createContext<T>({} as never);
}

interface ModifiedModalProps {
  onClose?: () => void;
}

export type ShowProp<T> = Omit<T, 'onClose'> & ModifiedModalProps;

interface RequiredModalProps {
  onClose: () => void;
}

export type ModalComponentProps<T> = ShowProp<T> & RequiredModalProps;

export function CreateModalPortal<T>(
  Component: React.FC<ModalComponentProps<T>>,
  props: ShowProp<T>
) {
  const modal = document.getElementById('modal-root');
  if (!modal) throw new Error('"#modal-root" element not found');
  return createPortal(
    <Component
      {...props}
      onClose={() => {
        props.onClose?.();
      }}
    />,
    modal
  );
}

type GetProps<T extends React.FC<any>> = ComponentProps<T>;

export type CreateModalProps<T extends React.FC<any> = React.FC<any>> = {
  Component: T;
  props: ShowProp<GetProps<T>>;
};

export const CreateModalContext =
  getEmptyContext<
    MutableRefObject<
      <T extends React.FC<any> = React.FC<any>>(modalData: CreateModalProps<T>) => void
    >
  >();

export const useCreateModal = () => useContext(CreateModalContext).current;
