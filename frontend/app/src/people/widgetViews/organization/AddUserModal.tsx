import React, { useState } from 'react';
import { useIsMobile } from 'hooks/uiHooks';
import styled from 'styled-components';
import { EuiLoadingSpinner } from '@elastic/eui';
import { useStores } from 'store';
import { nonWidgetConfigs } from 'people/utils/Constants';
import { Modal } from '../../../components/common';
import { colors } from '../../../config/colors';
import { AddUserModalProps } from './interface';

const color = colors['light'];

interface SmallBtnProps {
  selected: boolean;
}

interface UserProps {
  inactive: boolean;
}

const AddUserContainer = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
`;

const AddUserHeaderContainer = styled.div`
  display: flex;
  padding: 1.875rem;
  flex-direction: column;
`;

const AddUserHeader = styled.h2`
  color: #3c3f41;
  font-family: 'Barlow';
  font-size: 1.625rem;
  font-style: normal;
  font-weight: 800;
  line-height: normal;
  margin-bottom: 1.25rem;
`;

const SearchUserInput = styled.input`
  padding: 0.9375rem 0.875rem;
  border-radius: 0.375rem;
  border: 1px solid #dde1e5;
  background: #fff;
  width: 100%;
  color: #292c33;
  font-family: 'Barlow';
  font-size: 0.8125rem;
  font-style: normal;
  font-weight: 400;

  ::placeholder {
    color: #8e969c;
  }
`;

const UsersListContainer = styled.div`
  display: flex;
  flex-direction: column;
  padding: 1rem 1.875rem;
  background-color: #f2f3f5;
  height: 16rem;
  overflow-y: auto;
`;

const UserContianer = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1rem;
`;

const UserInfo = styled.div<UserProps>`
  display: flex;
  align-items: center;
  opacity: ${(p: any) => (p.inactive ? 0.3 : 1)};
`;

const UserImg = styled.img`
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  margin-right: 0.63rem;
  object-fit: cover;
`;

const Username = styled.p`
  color: #3c3f41;
  font-family: 'Barlow';
  font-size: 0.8125rem;
  font-style: normal;
  font-weight: 500;
  line-height: 1rem;
  margin-bottom: 0;
`;

const SmallBtn = styled.button<SmallBtnProps>`
  width: 5.375rem;
  height: 2rem;
  padding: 0.625rem;
  border-radius: 0.375rem;
  background: ${(p: any) => (p.selected ? '#618AFF' : '#dde1e5')};
  color: ${(p: any) => (p.selected ? '#FFF' : '#5f6368')};
  font-family: 'Barlow';
  font-size: 0.8125rem;
  font-style: normal;
  font-weight: 600;
  line-height: 0rem; /* 0% */
  letter-spacing: 0.00813rem;
  border: none;
`;

const FooterContainer = styled.div`
  display: flex;
  padding: 1.125rem 1.875rem;
  flex-direction: column;
  justify-content: center;
  align-items: center;
`;

const AddUserBtn = styled.button`
  height: 3rem;
  padding: 0.5rem 1rem;
  width: 100%;
  border-radius: 0.375rem;
  font-family: 'Barlow';
  font-size: 0.9375rem;
  font-style: normal;
  font-weight: 500;
  line-height: 0rem;
  letter-spacing: 0.00938rem;
  background: #618aff;
  box-shadow: 0px 2px 10px 0px rgba(97, 138, 255, 0.5);
  border: none;
  color: #fff;
  &:disabled {
    border: 1px solid rgba(0, 0, 0, 0.07);
    background: rgba(0, 0, 0, 0.04);
    color: rgba(142, 150, 156, 0.85);
    cursor: not-allowed;
    box-shadow: none;
  }
`;
const AddUserModal = (props: AddUserModalProps) => {
  const isMobile = useIsMobile();
  const { isOpen, close, onSubmit, loading } = props;
  const { main } = useStores();
  const people: any = (main.people && main.people.filter((f: any) => !f.hide)) || [];
  const [selectedPubkey, setSelectedPubkey] = useState<string>();

  const config = nonWidgetConfigs['organizationusers'];

  function handleSelectUser(pubkey: string) {
    setSelectedPubkey(pubkey);
  }

  function checkIsActive(pubkey: string) {
    if (selectedPubkey && pubkey !== selectedPubkey) {
      return true;
    } else {
      return false;
    }
  }

  return (
    <Modal
      visible={isOpen}
      style={{
        height: '100%',
        flexDirection: 'column'
      }}
      envStyle={{
        marginTop: isMobile ? 64 : 0,
        background: color.pureWhite,
        zIndex: 20,
        ...(config?.modalStyle ?? {}),
        maxHeight: '100%',
        borderRadius: '10px',
        minWidth: '22rem'
      }}
      overlayClick={close}
      bigCloseImage={close}
      bigCloseImageStyle={{
        top: '-18px',
        right: '-18px',
        background: '#000',
        borderRadius: '50%'
      }}
    >
      <AddUserContainer>
        <AddUserHeaderContainer>
          <AddUserHeader>Add New User</AddUserHeader>
          <SearchUserInput placeholder="Type to search ..." />
        </AddUserHeaderContainer>
        <UsersListContainer>
          {people.map((person: any, index: number) => (
            <UserContianer key={index}>
              <UserInfo inactive={checkIsActive(person.owner_pubkey)}>
                <UserImg src={person.img} alt="user" />
                <Username>{person.owner_alias}</Username>
              </UserInfo>
              <SmallBtn
                selected={person.owner_pubkey === selectedPubkey}
                onClick={() => handleSelectUser(person.owner_pubkey)}
              >
                Add
              </SmallBtn>
            </UserContianer>
          ))}
        </UsersListContainer>
        <FooterContainer>
          <AddUserBtn
            onClick={() => onSubmit({ owner_pubkey: selectedPubkey })}
            disabled={!selectedPubkey}
          >
            {loading ? <EuiLoadingSpinner size="s" /> : 'Add User'}
          </AddUserBtn>
        </FooterContainer>
      </AddUserContainer>
    </Modal>
  );
};

export default AddUserModal;
