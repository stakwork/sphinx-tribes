import React, { ChangeEvent, useState } from 'react';
import { useIsMobile } from 'hooks/uiHooks';
import { EuiLoadingSpinner } from '@elastic/eui';
import { useStores } from 'store';
import { Person } from 'store/main';
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
  const { main, ui } = useStores();
  const [selectedPubkey, setSelectedPubkey] = useState<string>();
  const [searchTerm, setSearchName] = useState<string>('');
  const [people, setPeople] = useState<Person[]>(
    (main.people && main.people.filter((f: any) => !f.hide)) || []
  );
  const currentUserPubkey = ui.meInfo?.owner_pubkey;

  const config = nonWidgetConfigs['organizationusers'];

  function checkIsActive(pubkey: string) {
    return !!(selectedPubkey && pubkey !== selectedPubkey);
  }

  const handleSearchUser = async (e: ChangeEvent<HTMLInputElement>) => {
    const name = e.target.value;
    setSearchName(name);
    const persons = await main.getPeopleByNameAliasPubkey(name);
    setPeople(persons.filter((person: Person) => person.owner_pubkey !== currentUserPubkey));
  };

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
        maxWidth: '22rem'
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
          <SearchUserInput
            value={searchTerm}
            onChange={handleSearchUser}
            placeholder="Type to search ..."
          />
        </AddUserHeaderContainer>
        <UsersListContainer>
          {people.length > 0 ? (
            people.map((person: any, index: number) => (
              <UserContianer key={index}>
                <UserInfo inactive={checkIsActive(person.owner_pubkey)}>
                  <UserImg src={person.img} alt="user" />
                  <Username>{person.owner_alias}</Username>
                </UserInfo>
                <SmallBtn
                  selected={person.owner_pubkey === selectedPubkey}
                  onClick={() => setSelectedPubkey(person.owner_pubkey)}
                >
                  Select
                </SmallBtn>
              </UserContianer>
            ))
          ) : (
            <p>No user with such alias</p>
          )}
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
