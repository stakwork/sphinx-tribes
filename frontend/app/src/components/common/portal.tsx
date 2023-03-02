import { PropsWithChildren } from 'react';
import { createPortal } from 'react-dom';

export const Portal = ({ children }: PropsWithChildren<{}>) => {
  const portal = createPortal(children, document.body);
  return portal;
};
