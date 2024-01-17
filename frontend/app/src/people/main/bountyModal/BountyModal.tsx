import { Modal } from 'components/common';
import { useIsMobile, usePerson } from 'hooks';
import { widgetConfigs } from 'people/utils/Constants';
import React, { useEffect, useState, useCallback } from 'react';
import { useHistory, useParams } from 'react-router-dom';
import { useStores } from 'store';
import { BountyModalProps } from 'people/interfaces';
import { PersonBounty } from 'store/main';
import FocusedView from '../FocusView';

const config = widgetConfigs.wanted;
export const BountyModal = ({ basePath, fromPage, bountyOwner }: BountyModalProps) => {
  const history = useHistory();
  const { wantedId, wantedIndex, personPubkey } = useParams<{
    wantedId: string;
    wantedIndex: string;
    personPubkey: string;
  }>();

  const { ui, main } = useStores();
  const { person } = usePerson(ui.selectedPerson);
  const [bounty, setBounty] = useState<PersonBounty[]>([]);
  const [afterEdit, setAfterEdit] = useState(false);
  const [loading, setLoading] = useState(true);
  const [ToDisplay, setToDisplay] = useState<any>();

  const personToDisplay = fromPage === 'usertickets' ? bountyOwner : person;

  const onGoBack = async () => {
    await main.getPersonCreatedBounties({}, personPubkey);
    await main.getPersonAssignedBounties({}, personPubkey);

    ui.setBountyPerson(0);
    history.push({
      pathname: basePath
    });
  };

  const getBounty = useCallback(
    async (afterEdit?: boolean) => {
      /** check for the bounty length, else the request
       * will be made continously which will lead to an
       * infinite loop and crash the app
       */
      if ((wantedId && !bounty.length) || afterEdit) {
        try {
          const bountyData = await main.getBountyById(Number(wantedId));
          console.log(bountyData, 'bd');
          setBounty(bountyData);
          if (personToDisplay === undefined) {
            setToDisplay(bountyData[0].person);
          } else {
            setToDisplay(personToDisplay);
          }
        } catch (error) {
          console.error('Error fetching bounty:', error);
        } finally {
          setLoading(false);
        }
      }
    },
    [bounty, main, wantedId, personToDisplay]
  );

  useEffect(() => {
    getBounty();
  }, [getBounty]);

  useEffect(() => {
    if (afterEdit) {
      getBounty(afterEdit);
      setAfterEdit(false);
    }
  }, [afterEdit, getBounty]);

  const isMobile = useIsMobile();

  if (loading) {
    return null;
  }

  if (isMobile) {
    return (
      <Modal visible={true} fill={true}>
        <FocusedView
          person={person}
          personBody={person}
          canEdit={false}
          selectedIndex={Number(wantedIndex)}
          config={config}
          goBack={onGoBack}
          setAfterEdit={setAfterEdit}
          bounty={bounty}
          fromBountyPage={true}
        />
      </Modal>
    );
  }

  return (
    <Modal
      visible={true}
      style={{
        background: 'rgba( 0 0 0 /75% )'
      }}
      envStyle={{
        maxHeight: '100vh',
        marginTop: 0,
        borderRadius: 0,
        background: '#fff',
        width: 'auto',
        minWidth: 500,
        maxWidth: '80%',
        zIndex: 20
      }}
      overlayClick={onGoBack}
      bigCloseImage={onGoBack}
      bigCloseImageStyle={{
        top: '18px',
        right: '-50px',
        borderRadius: '50%'
      }}
    >
      <FocusedView
        person={ToDisplay}
        personBody={person}
        canEdit={false}
        selectedIndex={Number(wantedIndex)}
        config={config}
        bounty={bounty}
        goBack={() => {
          onGoBack();
        }}
        fromBountyPage={true}
        setAfterEdit={setAfterEdit}
      />
    </Modal>
  );
};
