import React, { useState, useEffect } from 'react';
import { Wrap } from 'components/form/style';
import { useIsMobile } from 'hooks/uiHooks';
import { nonWidgetConfigs } from 'people/utils/Constants';
import styled from 'styled-components';
import { userHasRole, Roles } from 'helpers';
import avatarIcon from '../../../public/static/profile_avatar.svg';
import { Button, Modal } from '../../../components/common';
import { colors } from '../../../config/colors';
import { Check, CheckLabel, CheckLi, CheckUl, ModalTitle } from './style';
import { UserRolesModalProps } from './interface';
import { UserImage } from './style';
const color = colors['light'];

const UserRolesName = styled.p`
  color: #8e969c;
  margin: auto;
  margin-bottom: 20px;
  margin-top: 10px;
`;

const UserRolesWrap = styled(Wrap)`
  width: 100%;
`;

const UserRolesImage = styled(UserImage)`
  height: 80px;
  width: 80px;
  margin-left: auto;
  position: fixed;
  left: 50%;
  transform: translate(-40px, -90px);
  borderstyle: solid;
  borderradius: 50%;
  borderwidth: 4px;
  bordercolor: white;
`;

const HLine = styled.div`
  background-color: #ebedef;
  height: 1px;
  width: 100%;
  margin: 5px 0px 20px;
`;

const s_RolesCategories = [
  {
    name: 'Manage organization',
    roles: ['EDIT ORGANIZATION'],
    status: false
  },
  {
    name: 'Manage bounties',
    roles: ['ADD BOUNTY', 'UPDATE BOUNTY', 'DELETE BOUNTY', 'PAY BOUNTY', 'ADD ROLES'],
    status: false
  },
  {
    name: 'Fund organization',
    roles: ['ADD BUDGET'],
    status: false
  },
  {
    name: 'Withdraw from organization',
    roles: ['WITHDRAW BUDGET'],
    status: false
  },
  {
    name: 'View transaction history',
    roles: ['VIEW REPORT'],
    status: false
  },
  {
    name: 'Update members',
    roles: ['ADD USER', 'UPDATE USER', 'DELETE USER'],
    status: false
  }
];

const RolesModal = (props: UserRolesModalProps) => {
  const isMobile = useIsMobile();
  const { isOpen, close, bountyRolesData, roleChange, submitRoles, user, userRoles } = props;

  const config = nonWidgetConfigs['organizationusers'];

  const [rolesCategories, setRolesCategories] = useState(s_RolesCategories);

  const getRoles = () => {
    const newRoles = s_RolesCategories.map((cat: any) => {
      let hasRoles = false;
      cat.roles.forEach((element: Roles) => {
        hasRoles = hasRoles || userHasRole(bountyRolesData, userRoles, element);
      });

      cat.status = hasRoles;
      return cat;
    });
    setRolesCategories(newRoles);
  };

  useEffect(() => {
    getRoles();
  }, [userRoles]);

  const rolesChange = (role: any, s: any) => {
    // set the backend roles status using the map 'rolesCategories'
    role.roles.forEach((role: Roles) => roleChange(role, s));
    // set the checkbox status
    const newRoles = rolesCategories.map((r: any) => {
      if (r === role) {
        r.status = s.target.checked;
      }
      return r;
    });
    setRolesCategories(newRoles);
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
        padding: '20px 60px 20px 60px'
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
        <UserRolesImage src={user?.img || avatarIcon} />
        <UserRolesName>{user?.unique_name}</UserRolesName>
        <HLine />
        <ModalTitle style={{ fontWeight: '800', fontSize: '26px' }}>User Roles</ModalTitle>
        <CheckUl>
          {rolesCategories.map((role: any, i: number) => (
            <CheckLi key={i}>
              <Check
                checked={role.status}
                onChange={(s: any) => rolesChange(role, s)}
                type="checkbox"
                name={role.name}
                value={role.name}
              />
              <CheckLabel>{role.name}</CheckLabel>
            </CheckLi>
          ))}
        </CheckUl>
        <Button
          onClick={() => submitRoles()}
          style={{ width: '100%', height: '50px', borderRadius: '6px', alignSelf: 'center' }}
          color={'primary'}
          text={'Update roles'}
        />
      </UserRolesWrap>
    </Modal>
  );
};

export default RolesModal;
