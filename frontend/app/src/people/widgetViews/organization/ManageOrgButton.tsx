import React, { useState, useEffect, useCallback } from 'react';
import { useStores } from 'store';
import { OrganizationUser } from 'store/main';
import styled from 'styled-components';

const Button = styled.button`
  padding: 0.5rem 1rem;
  width: 7rem;
  height: 2.5rem;
  border-radius: 0.375rem;
  border: 1px solid var(--Input-Outline-1, #d0d5d8);
  background: #fff;
  box-shadow: 0px 1px 2px 0px rgba(0, 0, 0, 0.06);
  color: #5f6368;
  font-family: 'Barlow';
  font-size: 0.875rem;
  font-style: normal;
  font-weight: 500;
  line-height: 0rem;
  letter-spacing: 0.00875rem;
`;

const ManageButton = (props: { user_pubkey: string; org: any; action: () => void }) => {
  const [organizationUser, setOrganizationUser] = useState<OrganizationUser | undefined>();
  const { main, ui } = useStores();

  const { user_pubkey, org, action } = props;

  const isOrganizationAdmin = org?.owner_pubkey === ui.meInfo?.owner_pubkey;
  const pubkey = organizationUser?.owner_pubkey;
  const isUser = pubkey !== '' && ui.meInfo?.owner_pubkey === pubkey;

  const hasAccess = isOrganizationAdmin || isUser;

  const getUserRoles = useCallback(async () => {
    try {
      const user = await main.getOrganizationUser(org.uuid);
      setOrganizationUser(user);
    } catch (e) {
      console.error('User roles error', e);
    }
  }, [org.uuid, main, user_pubkey]);

  useEffect(() => {
    getUserRoles();
  }, [getUserRoles]);

  return <>{hasAccess && <Button onClick={action}>Manage</Button>}</>;
};

export default ManageButton;
