import React, { useState, useEffect, ChangeEvent } from 'react';
import { useIsMobile } from 'hooks/uiHooks';
import { EuiLoadingSpinner } from '@elastic/eui';
import { useStores } from 'store';
import { RolesCategory, s_RolesCategories } from 'helpers/helpers-extended';
import { nonWidgetConfigs } from 'people/utils/Constants';
import { handleDisplayRole } from 'helpers';
import { Modal } from '../../../components/common';
import { colors } from '../../../config/colors';
import { AssignUserModalProps } from './interface';
import {
  AssignRoleUserImage,
  AssignRoleUsername,
  AssignUserContainer,
  AssingUserBtn,
  Checkbox,
  Label,
  RoleContainer,
  RolesContainer,
  UserInfoContainer,
  UserRolesContainer,
  UserRolesTitle
} from './style';

interface BackendRoles {
  name: string;
  status: boolean;
}

const color = colors['light'];

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
    const { newDisplayedRoles, tempDataRole } = handleDisplayRole(displayedRoles);

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
          <AssignRoleUserImage
            src={user?.img || main.getUserAvatarPlaceholder(user?.owner_pubkey || '')}
            alt="profile"
          />
          <AssignRoleUsername>{user?.owner_alias}</AssignRoleUsername>
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
