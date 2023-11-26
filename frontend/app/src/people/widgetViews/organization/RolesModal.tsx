import React, { useState, useEffect, useCallback } from 'react';
import { Wrap } from 'components/form/style';
import { useIsMobile } from 'hooks/uiHooks';
import { nonWidgetConfigs } from 'people/utils/Constants';
import styled from 'styled-components';
import { userHasRole, Roles, s_RolesCategories } from 'helpers';
import { useStores } from 'store';
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
  border-style: solid;
  border-radius: 50%;
  border-width: 4px;
  border-color: white;
`;

const HLine = styled.div`
  background-color: #ebedef;
  height: 1px;
  width: 100%;
  margin: 5px 0px 20px;
`;

const RolesModal = (props: UserRolesModalProps) => {
  const { main } = useStores();
  const isMobile = useIsMobile();
  const { isOpen, close, user, uuid, addToast } = props;

  const roleData = main.bountyRoles.map((role: any) => ({
    name: role.name,
    status: false
  }));

  const config = nonWidgetConfigs['organizationusers'];
  const [bountyRolesData, setBountyRolesData] = useState(roleData);

  const rolesCheck = (bountyRoles: any[], roleName: string): boolean => {
    const role = bountyRoles.find((role: any) => role.name === roleName);
    return role.status;
  };

  const newRoles = s_RolesCategories.map((cat: any) => {
    let hasRoles = false;
    cat.roles.forEach((element: Roles) => {
      hasRoles = hasRoles || rolesCheck(bountyRolesData, element);
      cat.status = hasRoles;
    });
    return cat;
  });

  const [rolesCategories, setRolesCategories] = useState(newRoles);

  const pubkey = user?.owner_pubkey || '';

  const formatBountyRolesData = useCallback((userRoles: any[], bountyRolesData: any[]) => {
    // set all values to false, so every user data will be fresh
    const rolesData = bountyRolesData;

    userRoles.forEach((userRole: any) => {
      const index = rolesData.findIndex((role: any) => role.name === userRole.role);
      if (index !== -1) {
        rolesData[index]['status'] = true;
      }
    });

    setBountyRolesData(rolesData);
  }, []);

  const getRoles = useCallback(async () => {
    if (user) {
      const userRoles = await main.getUserRoles(uuid || '', pubkey);
      formatBountyRolesData(userRoles, bountyRolesData);

      const newRoles = s_RolesCategories.map((cat: any) => {
        let hasRoles = false;
        cat.roles.forEach((element: Roles) => {
          hasRoles = hasRoles || userHasRole(bountyRolesData, userRoles, element);
        });

        cat.status = hasRoles;
        return cat;
      });

      setRolesCategories(newRoles);
    }
  }, [main, pubkey, user, uuid, formatBountyRolesData]);

  useEffect(() => {
    getRoles();
  }, [getRoles]);

  const roleChange = (e: Roles, s: any): any => {
    const rolesData = bountyRolesData.map((role: any) => {
      if (role.name === e) {
        role.status = s.target.checked;
      }
      return role;
    });
    setBountyRolesData(rolesData);
  };

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

  const submitRoles = async () => {
    const roleData = bountyRolesData
      .filter((r: any) => r.status)
      .map((role: any) => ({
        owner_pubkey: user?.owner_pubkey,
        org_uuid: uuid,
        role: role.name
      }));

    if (uuid && user?.owner_pubkey) {
      const res = await main.addUserRoles(roleData, uuid, user.owner_pubkey);
      if (res.status === 200) {
        await main.getUserRoles(uuid, user.owner_pubkey);
      } else {
        addToast('Error: could not add user roles', 'danger');
      }
      close();
    }
  };

  return (
    <>
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
    </>
  );
};

export default RolesModal;
