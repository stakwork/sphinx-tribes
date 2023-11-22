import React, { useState } from 'react';
import { useIsMobile } from 'hooks/uiHooks';
import { EuiLoadingSpinner } from '@elastic/eui';
import { useStores } from 'store';
import { nonWidgetConfigs } from 'people/utils/Constants';
import { Modal } from '../../../components/common';
import { colors } from '../../../config/colors';
import { AddUserModalProps } from './interface';
import {
  AddUserBtn,
  AddUserContainer,
  AddUserHeader,
  AddUserHeaderContainer,
  FooterContainer,
  SearchUserInput,
  SmallBtn,
  UserContianer,
  UserImg,
  UserInfo,
  Username,
  UsersListContainer
} from './style';

const color = colors['light'];

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
