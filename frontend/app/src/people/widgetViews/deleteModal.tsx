import { EuiButton, EuiButtonEmpty, EuiModal, EuiOverlayMask, EuiText } from '@elastic/eui';
import React from 'react';
import styled from 'styled-components';

const DeleteTicketModal = ({ closeModal, confirmDelete }) => {
  return (
    <>
      <EuiOverlayMask>
        <EuiModal
          onClose={closeModal}
          style={{
            background: '#F2F3F5',
            padding: '50px 50px 30px 50px'
          }}
        >
          <EuiText>Are you sure you want to delete this Ticket?</EuiText>
          <ModalButtonContainer>
            <EuiButtonEmpty
              onClick={closeModal}
              style={{
                color: '#000'
              }}
            >
              Cancel
            </EuiButtonEmpty>
            <EuiButton
              onClick={confirmDelete}
              style={{
                background: '#fff',
                textDecoration: 'none',
                color: '#303030',
                border: '1px solid #909090'
              }}
            >
              Delete
            </EuiButton>
          </ModalButtonContainer>
        </EuiModal>
      </EuiOverlayMask>
    </>
  );
};

export default DeleteTicketModal;

const ModalButtonContainer = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  margin-top: 30px;
  padding: 2px;
`;
