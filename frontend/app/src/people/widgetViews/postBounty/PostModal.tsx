import React, { FC, useState, useCallback } from 'react';
import { useHistory } from 'react-router-dom';
import { observer } from 'mobx-react-lite';
import { colors } from '../../../config/colors';
import { useIsMobile } from '../../../hooks';
import { Modal } from '../../../components/common';
import { useStores } from '../../../store';
import FocusedView from '../../main/FocusView';
import { Widget } from '../../main/types';
import { widgetConfigs } from '../../utils/Constants';

const color = colors['light'];
export interface PostModalProps {
  isOpen: boolean;
  widget: Widget;
  onClose: () => void;
  onSucces?: () => void;
  onGoBack?: () => void;
}
export const PostModal: FC<PostModalProps> = observer(
  ({ isOpen, onClose, widget, onGoBack, onSucces }: any) => {
    const { main, ui } = useStores();
    const isMobile = useIsMobile();
    const [focusIndex, setFocusIndex] = useState(-1);
    const history = useHistory();

    const person: any = (main.people ?? []).find((f: any) => f.id === ui.selectedPerson);
    const { id } = person || {};
    const canEdit = id === ui.meInfo?.id;
    const config = widgetConfigs[widget];

    const getBountyData = useCallback(async () => {
      const response = await main.getPeopleBounties();
      return response[0].body.id;
    }, [main]);

    const ReCallBounties = async () => {
      /*
      TODO : after getting the better way to reload the bounty, this code will be removed.
      */
      const number = await getBountyData();
      history.push(`/bounty/${number}`);
    };

    const closeHandler = () => {
      onClose();
      onGoBack && onGoBack();
      setFocusIndex(-1);
    };
    const successHandler = () => {
      onClose();
      setFocusIndex(-1);
      onSucces && onSucces();
    };

    if (isMobile) {
      return (
        <>
          {isOpen && (
            <Modal visible={isOpen} fill={true}>
              <FocusedView
                person={person}
                canEdit={!canEdit}
                selectedIndex={focusIndex}
                config={config}
                onSuccess={successHandler}
                goBack={closeHandler}
              />
            </Modal>
          )}
        </>
      );
    }
    return (
      <>
        {' '}
        {isOpen && (
          <Modal
            visible={isOpen}
            style={{
              height: '100%'
            }}
            envStyle={{
              marginTop: isMobile ? 64 : 0,
              background: color.pureWhite,
              zIndex: 20,
              ...(config?.modalStyle ?? {}),
              maxHeight: '100%',
              borderRadius: '10px'
            }}
            overlayClick={undefined}
            bigCloseImage={closeHandler}
            bigCloseImageStyle={{
              top: '-18px',
              right: '-18px',
              background: '#000',
              borderRadius: '50%'
            }}
          >
            <FocusedView
              ReCallBounties={ReCallBounties}
              newDesign={true}
              person={person}
              canEdit={!canEdit}
              selectedIndex={focusIndex}
              config={config}
              onSuccess={successHandler}
              goBack={closeHandler}
            />
          </Modal>
        )}
      </>
    );
  }
);
