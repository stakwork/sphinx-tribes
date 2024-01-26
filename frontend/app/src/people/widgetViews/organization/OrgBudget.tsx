import React, { useState, useEffect, useCallback } from 'react';
import { DollarConverter, satToUsd, userHasRole } from 'helpers';
import { useStores } from 'store';
import styled from 'styled-components';

const OrganizationTextWrap = styled.div`
  margin-left: 1.5rem;
  display: flex;
  flex-direction: column;
  justify-content: center;

  @media only screen and (max-width: 470px) {
    margin-left: 0px;
    margin-top: 15px;
    justify-content: center;
    align-items: center;
  }
`;

const OrganizationText = styled.p<{ hasAccess: boolean }>`
  color: var(--Text-2, var(--Hover-Icon-Color, #3c3f41));
  font-family: 'Barlow';
  font-size: 1.25rem;
  font-style: normal;
  font-weight: 600;
  line-height: 1.1875rem; /* 95% */
  margin-bottom: ${(p: any) => (p.hasAccess ? '14px' : '0px')};
  text-transform: capitalize;
  @media only screen and (max-width: 700px) {
    font-size: 0.85rem;
  }
  @media only screen and (max-width: 500px) {
    font-size: 0.79rem;
  }
  @media only screen and (max-width: 470px) {
    font-size: 0.85rem;
    text-align: center;
  }
`;

const OrganizationBudgetText = styled.small`
  color: #5f6368;
  font-family: 'Barlow';
  font-size: 0.9375rem;
  font-style: normal;
  font-weight: 500;
  line-height: normal;
  letter-spacing: 0.00938rem;
  text-transform: uppercase;
  margin-top: auto;
  font-size: 0.9rem;
  @media only screen and (max-width: 700px) {
    font-size: 0.8rem;
  }
  @media only screen and (max-width: 500px) {
    font-size: 0.75rem;
  }
`;

const SatsGap = styled.span`
  margin: 0px 0.625rem;
  @media only screen and (max-width: 700px) {
    margin: 0px 0.3125rem;
  }
`;

const CurrencyUnit = styled.span`
  color: #8e969c;
`;

const OrganizationBudget = (props: { user_pubkey: string; org: any }) => {
  const [userRoles, setUserRoles] = useState<any[]>([]);

  const { main, ui } = useStores();
  const { user_pubkey, org } = props;

  const isOrganizationAdmin = org?.owner_pubkey === ui.meInfo?.owner_pubkey;

  const hasAccess =
    isOrganizationAdmin ||
    userHasRole(main.bountyRoles, userRoles, 'ADD USER') ||
    userHasRole(main.bountyRoles, userRoles, 'VIEW REPORT');

  const getUserRoles = useCallback(async () => {
    const userRoles = await main.getUserRoles(org.uuid, user_pubkey);
    setUserRoles(userRoles);
  }, [org.uuid, user_pubkey, main]);

  useEffect(() => {
    getUserRoles();
  }, [getUserRoles]);

  return (
    <OrganizationTextWrap>
      <OrganizationText hasAccess={hasAccess}>{org.name}</OrganizationText>
      {hasAccess && (
        <OrganizationBudgetText>
          {DollarConverter(org.budget ?? 0)}
          <CurrencyUnit>
            {' SAT'}
            <SatsGap>/</SatsGap>
            {satToUsd(org.budget ?? 0)} USD
          </CurrencyUnit>
        </OrganizationBudgetText>
      )}
    </OrganizationTextWrap>
  );
};

export default OrganizationBudget;
