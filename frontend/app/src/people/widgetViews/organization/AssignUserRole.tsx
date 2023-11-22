import React, { useState, useEffect, ChangeEvent } from 'react';
import { useIsMobile } from 'hooks/uiHooks';
import styled from 'styled-components';
import { EuiLoadingSpinner } from '@elastic/eui';
import { useStores } from 'store';
import { RolesCategory, s_RolesCategories } from 'helpers/helpers-extended';
import { nonWidgetConfigs } from 'people/utils/Constants';
import { Modal } from '../../../components/common';
import { colors } from '../../../config/colors';
import { AssignUserModalProps } from './interface';

interface BackendRoles {
  name: string;
  status: boolean;
}

const color = colors['light'];

const AssignUserContainer = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
  position: relative;
`;

const UserInfoContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: absolute;
  top: -2.5rem;
  left: 0;
  right: 0;
`;

const UserImage = styled.img`
  width: 5rem;
  height: 5rem;
  border-radius: 50%;
  background: #dde1e5;
  border: 4px solid #fff;
  object-fit: cover;
`;

const Username = styled.p`
  color: #3c3f41;
  text-align: center;
  font-family: 'Barlow';
  font-size: 1.25rem;
  font-style: normal;
  font-weight: 500;
  line-height: 1.625rem;
  margin-top: 0.69rem;
  margin-bottom: 0;
  text-transform: capitalize;
`;

const UserRolesContainer = styled.div`
  padding: 3.25rem 3rem 3rem 3rem;
  margin-top: 3.25rem;
`;

const UserRolesTitle = styled.h2`
  color: #3c3f41;
  font-family: 'Barlow';
  font-size: 1.625rem;
  font-style: normal;
  font-weight: 800;
  line-height: 1.625rem;
  margin-bottom: 2.81rem;
`;

const RolesContainer = styled.div`
  display: flex;
  flex-direction: column;
`;

const RoleContainer = styled.div`
  display: flex;
  align-items: center;
  margin-bottom: 1rem;
`;

const Checkbox = styled.input`
  margin-right: 1rem;
  width: 1rem;
  height: 1rem;
`;

const Label = styled.label`
  margin-bottom: 0;
  color: #1e1f25;
  font-family: 'Barlow';
  font-size: 0.9375rem;
  font-style: normal;
  font-weight: 500;
  line-height: 1.125rem;
`;

const AssingUserBtn = styled.button`
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
  margin-top: 3rem;
  color: #fff;
  &:disabled {
    border: 1px solid rgba(0, 0, 0, 0.07);
    background: rgba(0, 0, 0, 0.04);
    color: rgba(142, 150, 156, 0.85);
    cursor: not-allowed;
    box-shadow: none;
  }
`;

const AssignUserRoles = (props: AssignUserModalProps) => {
  const { isOpen, close, loading, user, setLoading, onSubmit, addToast } = props;
  const { main } = useStores();

  const config = nonWidgetConfigs['organizationusers'];

  const isMobile = useIsMobile();

  const roleData = main.bountyRoles.map((role: any) => ({
    name: role.name,
    status: false
  }));

  const [rolesData, setRolesData] = useState(roleData);
  const [displayedRoles, setDisplayedRoles] = useState(s_RolesCategories);

  useEffect(() => {
    // Set default data roles for first assign user
    const defaultRole = {
      'Manage bounties': true,
      'Fund organization': true,
      'Withdraw from organization': true,
      'View transaction history': true
    };

    const tempDataRole: { [id: string]: boolean } = {};

    const newDisplayedRoles = displayedRoles.map((role: RolesCategory) => {
      if (defaultRole[role.name]) {
        role.status = true;
        role.roles.forEach((dataRole: string) => (tempDataRole[dataRole] = true));
      }
      return role;
    });

    setRolesData((prev: BackendRoles[]) =>
      prev.map((role: BackendRoles) => {
        if (tempDataRole[role.name]) {
          role.status = true;
        }
        return role;
      })
    );
    setDisplayedRoles(newDisplayedRoles);

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleOnchange = (e: ChangeEvent<HTMLInputElement>, role: RolesCategory) => {
    // Handle display role
    const newDisplayRoles = displayedRoles.map((displayRole: RolesCategory) => {
      if (displayRole.name === role.name) {
        displayRole.status = e.target.checked;
      }
      return displayRole;
    });
    setDisplayedRoles(newDisplayRoles);

    // update backend role
    const newBackendObj = {};
    role.roles.forEach(
      (value: string) => (newBackendObj[`${value}`] = { value: e.target.checked })
    );

    setRolesData((prev: BackendRoles[]) =>
      prev.map((role: BackendRoles) => {
        if (newBackendObj[role.name]) {
          role.status = newBackendObj[role.name].value;
        }
        return role;
      })
    );
  };

  async function handleSubmitRole() {
    setLoading(true);
    try {
      await onSubmit(rolesData);
      setLoading(false);
      addToast('Roles Assigned Successfully', 'success');
      close();
    } catch (error) {
      setLoading(false);
      addToast('An internal error occured, please contact support', 'danger');
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
      <AssignUserContainer>
        <UserInfoContainer>
          <UserImage src={user?.img} alt="profile" />
          <Username>{user?.owner_alias}</Username>
        </UserInfoContainer>
        <UserRolesContainer>
          <UserRolesTitle>User Roles</UserRolesTitle>
          <RolesContainer>
            {displayedRoles.map((role: RolesCategory) => (
              <RoleContainer key={role.name}>
                <Checkbox
                  type={'checkbox'}
                  id={role.name}
                  checked={role.status}
                  onChange={(e: ChangeEvent<HTMLInputElement>) => handleOnchange(e, role)}
                />
                <Label htmlFor={role.name}>{role.name}</Label>
              </RoleContainer>
            ))}
          </RolesContainer>
          <AssingUserBtn onClick={handleSubmitRole}>
            {loading ? <EuiLoadingSpinner size="s" /> : 'Assign'}
          </AssingUserBtn>
        </UserRolesContainer>
      </AssignUserContainer>
    </Modal>
  );
};

export default AssignUserRoles;
