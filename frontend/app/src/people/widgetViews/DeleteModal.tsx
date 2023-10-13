import { EuiButton, EuiModal, EuiOverlayMask, EuiText } from '@elastic/eui';
import { DeleteTicketModalProps } from 'people/interfaces';
import React from 'react';
import styled from 'styled-components';
import avatarIcon from '../../public/static/profile_avatar.svg';

const ModalButtonContainer = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  margin-top: 30px;
  padding: 2px;
`;

const BoldText = styled.span`
  font-weight: bold;
`;

const WarningText = styled(EuiText)`
  text-align: center;
`;

const UserImage = styled.img`
  width: 80px;
  height: 80px;
  border-radius: 50%;
  margin-left: 26%;
  z-index: 200;
  position: absolute;
  margin-top: -24%;
  border: 4px solid #fff;
`;

const CancelButton = styled(EuiButton)`
  color: #000;
  background: #fff;
  border: 1px solid #d0d5d8;
  text-decoration: none !important;
  box-shadow: 0px 1px 2px 0px rgba(0, 0, 0, 0.06);
`;

const DeleteButton = styled(EuiButton)`
  background: #ed7474;
  color: #ffffff;
  border: none;
  text-decoration: none !important;
  box-shadow: 0px 2px 10px 0px rgba(237, 116, 116, 0.5);
`;

const DeleteTicketModal = ({
  closeModal,
  confirmDelete,
  text,
  imgUrl,
  userDelete
}: DeleteTicketModalProps) => (
  <EuiOverlayMask>
    <EuiModal
      onClose={closeModal}
      style={{
        background: '#F2F3F5',
        padding: '50px 50px 30px 50px'
      }}
    >
      {userDelete && <UserImage src={imgUrl || avatarIcon} />}
      <WarningText>
        Are you sure you want to <br />
        <BoldText>Delete this {text ? text : 'Ticket'}?</BoldText>
      </WarningText>
      <ModalButtonContainer>
        <CancelButton onClick={closeModal}>Cancel</CancelButton>
        <DeleteButton onClick={confirmDelete}>Delete</DeleteButton>
      </ModalButtonContainer>
    </EuiModal>
  </EuiOverlayMask>
);

export default DeleteTicketModal;
