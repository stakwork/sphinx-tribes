import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import ConnectCard from 'people/utils/ConnectCard';
import { TicketModalPage } from './TicketModalPage';
import Tickets from './Tickets';
import OrgTickets from './OrgTickets';

export const TicketsPage = () => {
  const [connectPerson, setConnectPerson] = useState<any>(null);
  const { uuid } = useParams<{ uuid: string }>();
  return (
    <>
      {uuid ? <OrgTickets /> : <Tickets />}
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
