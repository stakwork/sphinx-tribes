import { Modal } from 'components/common';
import { colors } from 'config';
import { useIsMobile } from 'hooks';
import { observer } from 'mobx-react-lite';
import FocusedView from 'people/main/FocusView';
import { widgetConfigs } from 'people/utils/Constants';
import React, { useEffect, useState } from 'react';
import { useHistory, useParams } from 'react-router-dom';
import { useStores } from 'store';
import { PersonBounty } from 'store/main';

const color = colors['light'];
const focusedDesktopModalStyles = widgetConfigs.wanted.modalStyle;

type Props = {
  setConnectPerson: (p: any) => void;
};

export const TicketModalPage = observer(({ setConnectPerson }: Props) => {
  const { main, modals, ui } = useStores();

  const history = useHistory();
  const [connectPersonBody, setConnectPersonBody] = useState<any>();
  const [activeListIndex, setActiveListIndex] = useState<number>(0);
  const [publicFocusIndex, setPublicFocusIndex] = useState(0);
  const { bountyId } = useParams<{ uuid: string; bountyId: string }>();

  const isMobile = useIsMobile();

  useEffect(() => {
    const activeIndex = main.peopleBounties.findIndex(
      (bounty: PersonBounty) => bounty.body.id === Number(bountyId)
    );
    const connectPerson = (main.peopleBounties ?? [])[activeIndex];

    setPublicFocusIndex(activeIndex);
    setActiveListIndex(activeIndex);

    setConnectPersonBody(connectPerson?.person);
  }, [main.peopleBounties, bountyId]);

  const goBack = () => {
    history.push('/bounties');
  };

  const prevArrHandler = () => {
    if (activeListIndex === 0) return;

    const { person, body } = main.peopleBounties[activeListIndex - 1];
    if (person && body) {
      history.replace(`/bounty/${body.id}`);
    }
  };
  const nextArrHandler = () => {
    if (activeListIndex + 1 > main.peopleBounties?.length) return;

    const { person, body } = main.peopleBounties[activeListIndex + 1];
    if (person && body) {
      history.replace(`/bounty/${body.id}`);
    }
  };

  if (isMobile) {
    return (
      <Modal visible={bountyId} fill={true}>
        <FocusedView
          person={connectPersonBody}
          personBody={connectPersonBody}
          canEdit={false}
          selectedIndex={publicFocusIndex}
          config={widgetConfigs.wanted}
          goBack={goBack}
        />
      </Modal>
    );
  }

  return (
    <Modal
      visible={bountyId && activeListIndex !== -1}
      envStyle={{
        background: color.pureWhite,
        ...focusedDesktopModalStyles,
        maxHeight: '100vh',
        zIndex: 20
      }}
      style={{
        background: 'rgba( 0 0 0 /75% )'
      }}
      overlayClick={goBack}
      bigCloseImage={goBack}
      bigCloseImageStyle={{
        top: '18px',
        right: '-50px',
        borderRadius: '50%'
      }}
      prevArrowNew={prevArrHandler}
      nextArrowNew={nextArrHandler}
    >
      <FocusedView
        person={connectPersonBody}
        personBody={connectPersonBody}
        canEdit={false}
        selectedIndex={publicFocusIndex}
        config={widgetConfigs.wanted}
        goBack={goBack}
        fromBountyPage={true}
        extraModalFunction={() => {
          goBack();
          if (ui.meInfo) {
            setConnectPerson(connectPersonBody);
          } else {
            modals.setStartupModal(true);
          }
        }}
      />
    </Modal>
  );
});
