import { Modal } from 'components/common';
import { colors } from 'config';
import { useIsMobile } from 'hooks';
import { observer } from 'mobx-react-lite';
import FocusedView from 'people/main/FocusView';
import { widgetConfigs } from 'people/utils/Constants';
import React, { useEffect, useState, useMemo } from 'react';
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
  const { bountyId } = useParams<{ uuid: string; bountyId: string }>();

  const isMobile = useIsMobile();

  const search = useMemo(() => {
    const s = new URLSearchParams(location.search);
    return {
      owner_id: s.get('owner_id'),
      created: s.get('created')
    };
  }, [location.search]);

  useEffect(() => {
    const activeIndex = bountyId
      ? main.peopleBounties.findIndex((bounty: PersonBounty) => bounty.body.id === Number(bountyId))
      : (main.peopleBounties ?? []).findIndex(findPerson(search));
    const connectPerson = (main.peopleBounties ?? [])[activeIndex];

    setPublicFocusIndex(activeIndex);
    setActiveListIndex(activeIndex);

    setConnectPersonBody(connectPerson?.person);
  }, [main.peopleBounties, bountyId, search]);

  const goBack = () => {
    history.push('/bounties');
  };

  const prevArrHandler = () => {
    if (activeListIndex === 0) return;

    const { person, body } = main.peopleBounties[activeListIndex - 1];
    if (person && body) {
      if (bountyId) history.replace(`/bounty/${body.id}`);
      else {
        history.replace({
          pathname: history?.location?.pathname,
          search: `?owner_id=${person?.owner_pubkey}&created=${body?.created}`,
          state: {
            owner_id: person?.owner_pubkey,
            created: body?.created
          }
        });
      }
    }
  };
  const nextArrHandler = () => {
    if (activeListIndex + 1 > main.peopleBounties?.length) return;

    const { person, body } = main.peopleBounties[activeListIndex + 1];
    if (person && body) {
      if (bountyId) history.replace(`/bounty/${body.id}`);
      else {
        history.replace({
          pathname: history?.location?.pathname,
          search: `?owner_id=${person?.owner_pubkey}&created=${body?.created}`,
          state: {
            owner_id: person?.owner_pubkey,
            created: body?.created
          }
        });
      }
    }
  };

  if (isMobile) {
    return (
      <Modal visible={activeListIndex !== -1 && (bountyId || search.created)} fill={true}>
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
      visible={activeListIndex !== -1 && (bountyId || search.created)}
      envStyle={{
        background: color.pureWhite,
        ...focusedDesktopModalStyles,
        maxHeight: '100vh',
        zIndex: 20,
        overflow: 'hidden'
      }}
      style={{
        background: 'rgba( 0 0 0 /75% )',
        overflow: 'hidden'
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
