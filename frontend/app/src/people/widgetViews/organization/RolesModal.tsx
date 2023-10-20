import React from 'react';
import { Wrap } from 'components/form/style';
import { useIsMobile } from 'hooks/uiHooks';
import { nonWidgetConfigs } from 'people/utils/Constants';
import styled from 'styled-components';
import avatarIcon from '../../../public/static/profile_avatar.svg';
import { Button, Modal } from '../../../components/common';
import { colors } from '../../../config/colors';
import { Check, CheckLabel, CheckLi, CheckUl, ModalTitle } from './style';
import { UserRolesModalProps } from './interface';
import { UserImage } from './style';

const color = colors['light'];

const UserRolesHeader = styled.div`
  display: flex;
  flex-direction: row;
`;

const UserRolesName = styled.p`
  color: #8E969C;
  margin: 5px 0px;
`;

const UserRolesWrap = styled(Wrap)`
  width: 100%;
`;

const RolesModal = (props: UserRolesModalProps) => {
  const isMobile = useIsMobile();
  const { isOpen, close, bountyRolesData, roleChange, submitRoles, user } = props;

  const config = nonWidgetConfigs['organizationusers'];

  return (
    <Modal
      visible={isOpen}
      style={{
        height: '100%',
        flexDirection: 'column',
      }}
      envStyle={{
        marginTop: isMobile ? 64 : 0,
        background: color.pureWhite,
        zIndex: 20,
        ...(config?.modalStyle ?? {}),
        maxHeight: '100%',
        borderRadius: '10px',
        padding: '20px 60px 10px 60px'
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
      <UserRolesWrap newDesign={true}>
        <UserRolesHeader>
          <div>
            <UserRolesName>{user?.unique_name}</UserRolesName>
            <ModalTitle>User Roles</ModalTitle>
          </div>
          <UserImage style={{height: '80px', width:'auto', marginLeft: 'auto'}} src={user?.img || avatarIcon} />
        </UserRolesHeader>
        <CheckUl>
          {bountyRolesData.map((role: any, i: number) => {
            const capitalizeWords:string = role.name.charAt(0).toUpperCase() + role.name.slice(1).toLocaleLowerCase();
            
            return (<CheckLi key={i}>
              <Check
                checked={role.status}
                onChange={roleChange}
                type="checkbox"
                name={role.name}
                value={role.name}
              />
              <CheckLabel>{capitalizeWords}</CheckLabel>
            </CheckLi>
          )})}
        </CheckUl>
        <Button
          onClick={() => submitRoles()}
          style={{ width: '150px', height: '50px', borderRadius: '5px', alignSelf: 'center' }}
          color={'primary'}
          text={'Update roles'}
        />
      </UserRolesWrap>
    </Modal>
  );
};

export default RolesModal;
