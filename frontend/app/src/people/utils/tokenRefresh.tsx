import React, { useEffect, useState } from 'react';
import { Button, Modal } from '../../sphinxUI';
import { useStores } from '../../store';

let timeout;

export default function TokenRefresh() {
  const { main, ui } = useStores();
  const [show, setShow] = useState(false);

  useEffect(() => {
    timeout = setTimeout(async () => {
      console.log('run token refresh!');
      if (ui.meInfo) {
        const res = await main.refreshJwt();
        if (res && res.jwt) {
          console.log('token refreshed!');
          ui.setMeInfo({ ...ui.meInfo, jwt: res.jwt });
        } else {
          console.log('kick!');
          ui.setMeInfo(null);
          ui.setSelectedPerson(0);
          ui.setSelectingPerson(0);
          setShow(true);
          // run this to reset state
          main.getPeople();
        }
      }
    }, 6000);

    return function cleanup() {
      clearTimeout(timeout);
    };
  }, []);

  return (
    <>
      <Modal visible={show}>
        <div style={{ display: 'flex', flexDirection: 'column', width: 250 }}>
          <div style={{ marginBottom: 20, textAlign: 'center' }}>
            Your session expired. Please log in again.
          </div>
          <Button text={'OK'} color={'widget'} onClick={() => setShow(false)} />
        </div>
      </Modal>
    </>
  );
}
