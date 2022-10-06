import React, { useState } from 'react';
import OfferView from '../widgetViews/offerView';
import WantedView from '../widgetViews/wantedView';
import PostView from '../widgetViews/postView';
import styled from 'styled-components';
import { useIsMobile } from '../../hooks';
import { useStores } from '../../store';
import { useObserver } from 'mobx-react-lite';
import { widgetConfigs } from '../utils/constants';
import { Spacer } from '../main/body';
import NoResults from '../utils/noResults';
import api from '../../api';
import {
  EuiButton,
  EuiButtonEmpty,
  EuiModal,
  EuiModalBody,
  EuiModalFooter,
  EuiModalHeader,
  EuiModalHeaderTitle,
  EuiOverlayMask,
  EuiText
} from '@elastic/eui';
import { uiStore } from '../../store/ui';

export default function WidgetSwitchViewer(props) {
  const { main, ui } = useStores();
  const isMobile = useIsMobile();
  const [showDeleteModal, setShowDeleteModal] = useState<boolean>(false);
  const [deletePayload, setDeletePayload] = useState<object>({});
  const closeModal = () => setShowDeleteModal(false);
  const showModal = () => setShowDeleteModal(true);

  const panelStyles = isMobile
    ? {
        minHeight: 132
      }
    : {
        maxWidth: 291,
        minWidth: 291,
        marginRight: 20,
        marginBottom: 20,
        minHeight: 472
      };

  return useObserver(() => {
    const { peoplePosts, peopleWanteds, peopleOffers } = main;

    let { selectedWidget, onPanelClick } = props;

    if (!selectedWidget) {
      return <div style={{ height: 200 }} />;
    }

    const listSource = {
      post: peoplePosts,
      wanted: peopleWanteds,
      offer: peopleOffers
    };

    const activeList = listSource[selectedWidget];

    let searchKeys: any = widgetConfigs[selectedWidget]?.schema?.map((s) => s.name) || [];
    let foundDynamicSchema = widgetConfigs[selectedWidget]?.schema?.find((f) => f.dynamicSchemas);
    // if dynamic schema, get all those fields
    if (foundDynamicSchema) {
      let dynamicFields: any = [];
      foundDynamicSchema.dynamicSchemas?.forEach((ds) => {
        ds.forEach((f) => {
          if (!dynamicFields.includes(f.name)) dynamicFields.push(f.name);
        });
      });
      searchKeys = dynamicFields;
    }

    const deleteTicket = async (payload: any) => {
      const info = uiStore.meInfo as any;
      const URL = info.url.startsWith('http') ? info.url : `https://${info.url}`;

      await fetch(URL + `/delete_ticket`, {
        method: 'POST',
        body: JSON.stringify(payload),
        headers: {
          'x-jwt': info.jwt,
          'Content-Type': 'application/json'
        }
      });
    };

    const confirmDelete = async () => {
      if (!!deletePayload) deleteTicket(deletePayload);
      closeModal();
    };

    let modal;

    if (showDeleteModal) {
      modal = (
        <EuiOverlayMask>
          <EuiModal
            onClose={closeModal}
            initialFocus="[name=popswitch]"
            style={{
              background: '#fff',
              padding: '50px 50px 30px 50px'
            }}>
            <EuiText>Are you sure you want to delete this Ticket?</EuiText>
            <ModalButtonContainer>
              <EuiButtonEmpty
                onClick={closeModal}
                style={{
                  color: '#000'
                }}>
                Cancel
              </EuiButtonEmpty>
              <EuiButton
                onClick={confirmDelete}
                style={{
                  background: '#fff',
                  textDecoration: 'none',
                  color: '#303030',
                  border: '1px solid #909090'
                }}>
                Delete
              </EuiButton>
            </ModalButtonContainer>
          </EuiModal>
        </EuiOverlayMask>
      );
    }

    const listItems =
      activeList && activeList.length ? (
        activeList.map((item, i) => {
          const { person, body } = item;

          const conditionalStyles = body?.paid
            ? {
                border: isMobile ? '2px 0 0 0 solid #dde1e5' : '1px solid #dde1e5',
                boxShadow: 'none'
              }
            : {};

          // if this person has entries for this widget
          return (
            <Panel
              isMobile={isMobile}
              key={person?.owner_pubkey + i + body?.created}
              onClick={() => {
                if (onPanelClick) onPanelClick(person, body);
              }}
              style={{
                ...panelStyles,
                ...conditionalStyles,
                cursor: 'pointer',
                padding: 0,
                overflow: 'hidden'
              }}>
              {selectedWidget === 'post' ? (
                <PostView
                  showName
                  key={i + person.owner_pubkey + 'pview'}
                  person={person}
                  {...body}
                />
              ) : selectedWidget === 'offer' ? (
                <OfferView
                  showName
                  key={i + person.owner_pubkey + 'oview'}
                  person={person}
                  {...body}
                />
              ) : selectedWidget === 'wanted' ? (
                <WantedView
                  showName
                  key={i + person.owner_pubkey + 'wview'}
                  person={person}
                  showModal={showModal}
                  setDeletePayload={setDeletePayload}
                  {...body}
                />
              ) : null}
            </Panel>
          );
        })
      ) : (
        <NoResults />
      );

    return (
      <>
        {listItems}
        <Spacer key={'spacer'} />
        {modal}
      </>
    );
  });
}

interface PanelProps {
  isMobile?: boolean;
}

const Panel = styled.div<PanelProps>`
  position: relative;
  background: #ffffff;
  color: #000000;
  padding: 20px;
  box-shadow: ${(p) => (p.isMobile ? 'none' : '0px 0px 6px rgb(0 0 0 / 7%)')};
  border-bottom: ${(p) => (p.isMobile ? '2px solid #EBEDEF' : 'none')};
`;

const ModalButtonContainer = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  margin-top: 30px;
  padding: 2px;
`;
