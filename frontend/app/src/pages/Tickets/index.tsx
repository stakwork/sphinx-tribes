import React, { useState } from 'react';
import { TicketModalPage } from './TicketModalPage';
import Tickets from './Tickets';
import ConnectCard from 'people/utils/connectCard';

export const TicketsPage = () => {
  const [connectPerson, setConnectPerson] = useState<any>(null);
  return (
    <>
      <Tickets />
      <TicketModalPage setConnectPerson={setConnectPerson} />
      {connectPerson && (
        <ConnectCard
          dismiss={() => setConnectPerson(null)}
          modalStyle={{
            top: '-64px',
            height: 'calc(100% + 64px)'
          }}
          person={connectPerson}
          visible={true}
        />
      )}
    </>
  );
};
