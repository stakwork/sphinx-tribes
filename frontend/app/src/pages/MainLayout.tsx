import React, { FC, PropsWithChildren } from 'react';
import { colors } from '../config/colors';

export type MainLayotProps = PropsWithChildren<{
  header: React.ReactElement;
}>;

export const MainLayout: FC<MainLayotProps> = ({ header, children }: { header: JSX.Element, children?: React.ReactNode }) => {
  const c = colors['light'];
  return (
    <div className="app" style={{ background: c.background }}>
      {header}
      {children}
    </div>
  );
};
