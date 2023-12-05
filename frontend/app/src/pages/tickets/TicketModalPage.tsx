import { Modal } from 'components/common';
import { colors } from 'config';
import { useIsMobile } from 'hooks';
import { observer } from 'mobx-react-lite';
import FocusedView from 'people/main/FocusView';
import { widgetConfigs } from 'people/utils/Constants';
import React, { useEffect, useState, useMemo, useCallback } from 'react';
import { useHistory, useLocation, useParams } from 'react-router-dom';
import { useStores } from 'store';
import { PersonBounty } from 'store/main';

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
  const [removeNextAndPrev, setRemoveNextAndPrev] = useState(false);
  const { bountyId } = useParams<{ uuid: string; bountyId: string }>();
  const [activeBounty, setActiveBounty] = useState<PersonBounty[]>([]);
  const [visible, setVisible] = useState(false);

  const isMobile = useIsMobile();

  const search = useMemo(() => {
    const s = new URLSearchParams(location.search);
    return {
      owner_id: s.get('owner_id'),
      created: s.get('created')
    };
  }, [location.search]);

  const getBounty = useCallback(async () => {
    let bounty;

    if (bountyId) {
      bounty = await main.getBountyById(Number(bountyId));
    } else if (search && search.created) {
      bounty = await main.getBountyByCreated(Number(search.created));
    }

    const activeIndex = bountyId
      ? main.peopleBounties.findIndex((bounty: PersonBounty) => bounty.body.id === Number(bountyId))
      : (main.peopleBounties ?? []).findIndex(findPerson(search));
    const connectPerson = bounty && bounty.length ? bounty[0].person : [];

    setPublicFocusIndex(activeIndex);
    setActiveListIndex(activeIndex);
    setConnectPersonBody(connectPerson);

    const visible = bounty && bounty.length > 0;

    setActiveBounty(bounty);
    setVisible(visible);
  }, [bountyId, main, search]);

  useEffect(() => {
    getBounty();
  }, [getBounty, removeNextAndPrev]);

  const goBack = async () => {
    setVisible(false);
    await main.getPeopleBounties({ page: 1, resetPage: true });
    history.goBack();
  };

  const directionHandler = (person: any, body: any) => {
    if (person && body) {
      if (bountyId) {
        history.replace(`/bounty/${body.id}`);
      }
    }
  };

  const prevArrHandler = () => {
    if (activeListIndex === 0) return;

    const { person, body } = main.peopleBounties[activeListIndex - 1];
    directionHandler(person, body);
  };

  const nextArrHandler = () => {
    if (activeListIndex + 1 > main.peopleBounties?.length) return;

    const { person, body } = main.peopleBounties[activeListIndex + 1];
    directionHandler(person, body);
  };

  if (isMobile) {
    return (
      <>
        {visible && (
          <Modal visible={visible} fill={true}>
            <FocusedView
              person={connectPersonBody}
              personBody={connectPersonBody}
              canEdit={false}
              selectedIndex={publicFocusIndex}
              config={widgetConfigs.wanted}
              bounty={activeBounty}
              goBack={goBack}
            />
          </Modal>
        )}
      </>
    );
  }

  return (
    <>
      {visible && (
        <Modal
          visible={visible}
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
          prevArrowNew={removeNextAndPrev ? undefined : prevArrHandler}
          nextArrowNew={removeNextAndPrev ? undefined : nextArrHandler}
        >
          <FocusedView
            setRemoveNextAndPrev={setRemoveNextAndPrev}
            person={connectPersonBody}
            personBody={connectPersonBody}
            canEdit={false}
            selectedIndex={publicFocusIndex}
            config={widgetConfigs.wanted}
            goBack={goBack}
            bounty={activeBounty}
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
      )}
    </>
  );
});