import { Modal } from 'components/common';
import { colors } from 'config';
import { useIsMobile } from 'hooks';
import { observer } from 'mobx-react-lite';
import FocusedView from 'people/main/FocusView';
import { widgetConfigs } from 'people/utils/Constants';
import React, { useEffect, useMemo, useState } from 'react';
import { useHistory, useLocation } from 'react-router-dom';
import { useStores } from 'store';

const color = colors['light'];
const focusedDesktopModalStyles = widgetConfigs.wanted.modalStyle;

const findPerson = (search: any) => (item: any) => {
  const { person, body } = item;
  return search.owner_id === person.owner_pubkey && search.created === `${body.created}`;
};

type Props = {
  setConnectPerson: (p: any) => void;
};
export const TicketModalPage = observer(({ setConnectPerson }: Props) => {
  const location = useLocation();
  const { main, modals, ui } = useStores();

  const history = useHistory();
  const [connectPersonBody, setConnectPersonBody] = useState<any>();
  const [activeListIndex, setActiveListIndex] = useState<number>(0);
  const [publicFocusIndex, setPublicFocusIndex] = useState(0);
  const isMobile = useIsMobile();

  const search = useMemo(() => {
    const s = new URLSearchParams(location.search);
    return {
      owner_id: s.get('owner_id'),
      created: s.get('created')
    };
  }, [location.search]);

  const publicFocusPerson = useMemo(
    () => main.people.find(({ owner_pubkey }: any) => owner_pubkey === search.owner_id),
    [main.people, search.owner_id]
  );

  useEffect(() => {
    const activeIndex = (main.peopleWanteds ?? []).findIndex(findPerson(search));
    const connectPerson = (main.peopleWanteds ?? [])[activeIndex];

    setPublicFocusIndex(activeIndex);
    setActiveListIndex(activeIndex);

    setConnectPersonBody(connectPerson?.person);
  }, [main.peopleWanteds, publicFocusPerson, search]);

  const goBack = () => {
    history.push('/tickets');
  };

  const prevArrHandler = () => {
    if (activeListIndex === 0) return;

    const { person, body } = main.peopleWanteds[activeListIndex-1];
    if (person && body) {
      history.replace({
        pathname: history?.location?.pathname,
        search: `?owner_id=${person?.owner_pubkey}&created=${body?.created}`,
        state: {
          owner_id: person?.owner_pubkey,
          created: body?.created
        }
      });
    }
  };
  const nextArrHandler = () => {
    console.log("edakun")
    if (activeListIndex + 1 > main.peopleWanteds?.length) return;

    const { person, body } = main.peopleWanteds[activeListIndex + 1];
    if (person && body) {
      history.replace({
        pathname: history?.location?.pathname,
        search: `?owner_id=${person?.owner_pubkey}&created=${body?.created}`,
        state: {
          owner_id: person?.owner_pubkey,
          created: body?.created
        }
      });
    }
  };

  if (isMobile) {
    return (
      <Modal visible={search.created && search.owner_id} fill={true}>
        <FocusedView
          person={publicFocusPerson}
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
      visible={search.created && search.owner_id}
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
        person={publicFocusPerson}
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
