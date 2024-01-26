import { observer } from 'mobx-react-lite';
import StartUpModal from 'people/utils/StartUpModal';
import React from 'react';
import { useStores } from 'store';

export const GetWorkStartupModal = observer(() => {
  const { modals } = useStores();
  return (
    <>
      {modals.startupModal && (
        <StartUpModal
          closeModal={() => modals.setStartupModal(false)}
          dataObject={'getWork'}
          buttonColor={'primary'}
        />
      )}
    </>
  );
});
