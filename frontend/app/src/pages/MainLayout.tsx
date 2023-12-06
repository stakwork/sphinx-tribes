import React, { FC, PropsWithChildren } from 'react';
import { useLocation } from 'react-router-dom';
import { colors } from '../config/colors';

export type MainLayoutProps = PropsWithChildren<{
  header: React.ReactElement;
}>;

export const MainLayout: FC<MainLayoutProps> = ({
  header,
  children
}: {
  header: JSX.Element;
  children?: React.ReactNode;
}) => {
  const c = colors['light'];
  const location = useLocation(); 

  return (
    <>
      {location.pathname === '/superadmin' ? (
        <div className="app" style={{ background: c.background }}>
          {children}
        </div>
      ) : (
        <div className="app" style={{ background: c.background }}>
          {header}
          {children}
        </div>
      )}
    </>
  );
};
