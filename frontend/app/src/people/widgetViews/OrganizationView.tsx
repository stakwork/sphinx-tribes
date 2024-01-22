import React, { useState, useEffect, useCallback } from 'react';
import styled from 'styled-components';
import PageLoadSpinner from 'people/utils/PageLoadSpinner';
import NoResults from 'people/utils/OrgNoResults';
import { useStores } from 'store';
import { Organization } from 'store/main';
import { Button } from 'components/common';
import { useIsMobile } from 'hooks/uiHooks';
import { Modal } from '../../components/common';
import avatarIcon from '../../public/static/profile_avatar.svg';
import { colors } from '../../config/colors';
import { widgetConfigs } from '../utils/Constants';
import { Person } from '../../store/main';
import OrganizationDetails from './OrganizationDetails';
import ManageButton from './organization/ManageOrgButton';
import OrganizationBudget from './organization/OrgBudget';
import AddOrganization from './organization/AddOrganization';

const color = colors['light'];

const Container = styled.div`
  display: flex;
  flex-flow: column wrap;
  min-width: 100%;
  min-height: 100%;
  flex: 1 1 100%;
  margin: -20px -30px;

  .organizations {
    padding: 1.25rem 2.5rem;
    @media only screen and (max-width: 800px) {
      padding: 1.25rem;
    }
  }
`;

const OrganizationWrap = styled.a`
  display: flex;
  flex-direction: row;
  width: 100%;
  background: white;
  padding: 1.5rem;
  border-radius: 0.375rem;
  box-shadow: 0px 1px 2px 0px rgba(0, 0, 0, 0.15);
  @media only screen and (max-width: 800px) {
    padding: 1rem 0px;
  }
  @media only screen and (max-width: 700px) {
    padding: 0.75rem 0px;
    margin-bottom: 10px;
  }
  @media only screen and (max-width: 500px) {
    padding: 0px;
  }

  &:hover {
    text-decoration: none !important;
  }
`;

const ButtonIconLeft = styled.button`
  display: flex;
  align-items: center;
  justify-content: center;
  height: 2.5rem;
  column-gap: 0.75rem;
  padding: 0.5rem 0.75rem;
  color: #5f6368;
  font-family: 'Barlow';
  font-size: 0.875rem;
  font-style: normal;
  font-weight: 500;
  line-height: 0rem;
  letter-spacing: 0.00875rem;
  border-radius: 0.375rem;
  border: 1px solid #d0d5d8;
  background: #fff;
  box-shadow: 0px 1px 2px 0px rgba(0, 0, 0, 0.06);

  :disabled {
    cursor: not-allowed;
  }
`;

const IconImg = styled.img`
  width: 1.25rem;
  height: 1.25rem;
`;

const OrganizationData = styled.div`
  display: flex;
  align-items: center;
  flex-direction: row;
  width: 100%;
  @media only screen and (max-width: 470px) {
    flex-direction: column;
    justify-content: center;
    border: 1px solid #ccc;
    border-radius: 10px;
    padding: 15px 0px;
  }
`;

const OrganizationImg = styled.img`
  width: 4rem;
  height: 4rem;
  border-radius: 50%;
  object-fit: cover;
  @media only screen and (max-width: 700px) {
    width: 55px;
    height: 55px;
  }
  @media only screen and (max-width: 500px) {
    width: 3rem;
    height: 3rem;
  }
  @media only screen and (max-width: 470px) {
    width: 3.75rem;
    height: 3.75rem;
  }
`;

const OrganizationContainer = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
  gap: 1rem;
`;

const OrgHeadWrap = styled.div`
  display: flex;
  align-items: center;
  margin-top: 5px;
  margin-bottom: 20px;
`;

const OrgText = styled.div`
  color: #3c3f41;
  font-family: 'Barlow';
  font-size: 1.5rem;
  font-style: normal;
  font-weight: 600;
  line-height: 1.1875rem;
  @media only screen and (max-width: 700px) {
    font-size: 1.1rem;
  }
  @media only screen and (max-width: 700px) {
    font-size: 0.95rem;
  }
`;
const OrganizationActionWrap = styled.div`
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 15px;
  @media only screen and (max-width: 470px) {
    margin-left: 0;
    margin-top: 20px;
  }
`;

const Organizations = (props: { person: Person }) => {
  const [loading, setIsLoading] = useState<boolean>(false);
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const [detailsOpen, setDetailsOpen] = useState<boolean>(false);
  const [organization, setOrganization] = useState<Organization>();
  const { main, ui } = useStores();
  const isMobile = useIsMobile();
  const config = widgetConfigs['organizations'];
  const isMyProfile = ui?.meInfo?.pubkey === props?.person?.owner_pubkey;

  const user_pubkey = ui.meInfo?.owner_pubkey;

  const getUserOrganizations = useCallback(async () => {
    setIsLoading(true);
    if (ui.selectedPerson) {
      await main.getUserOrganizations(ui.selectedPerson);
    }
    setIsLoading(false);
  }, [main, ui.selectedPerson]);

  useEffect(() => {
    getUserOrganizations();
  }, [getUserOrganizations]);

  const closeHandler = () => {
    setIsOpen(false);
  };

  const closeDetails = () => {
    setDetailsOpen(false);
  };

  // renders org as list item
  const orgUi = (org: any, key: number) => {
    const btnDisabled = (!org.bounty_count && org.bount_count !== 0) || !org.uuid;
    return (
      <OrganizationWrap key={key}>
        <OrganizationData>
          <OrganizationImg src={org.img || avatarIcon} />
          <OrganizationBudget org={org} user_pubkey={user_pubkey ?? ''} />
          <OrganizationActionWrap>
            {user_pubkey && (
              <ManageButton
                org={org}
                user_pubkey={user_pubkey ?? ''}
                action={() => {
                  setOrganization(org);
                  setDetailsOpen(true);
                }}
              />
            )}
            <ButtonIconLeft
              disabled={btnDisabled}
              onClick={() => window.open(`/org/bounties/${org.uuid}`, '_target')}
            >
              View Bounties
              <IconImg src="/static/open_in_new_grey.svg" alt="open_in_new_tab" />
            </ButtonIconLeft>
          </OrganizationActionWrap>
        </OrganizationData>
      </OrganizationWrap>
    );
  };

  // renders list of orgs with header
  const renderOrganizations = () => {
    if (main.organizations.length) {
      return (
        <div className="organizations">
          <OrgHeadWrap>
            <OrgText>Organizations</OrgText>
            {isMyProfile && (
              <Button
                leadingIcon={'add'}
                height={isMobile ? 40 : 45}
                text="Add Organization"
                onClick={() => setIsOpen(true)}
                style={{ marginLeft: 'auto', borderRadius: 10 }}
              />
            )}
          </OrgHeadWrap>
          <OrganizationContainer>
            {main.organizations.map((org: Organization, i: number) => orgUi(org, i))}
          </OrganizationContainer>
        </div>
      );
    } else {
      return (
        <Container>
          <NoResults showAction={isMyProfile} action={() => setIsOpen(true)} />
        </Container>
      );
    }
  };

  return (
    <Container>
      <PageLoadSpinner show={loading} />
      {detailsOpen && (
        <OrganizationDetails
          close={closeDetails}
          org={organization}
          resetOrg={(newOrg: Organization) => setOrganization(newOrg)}
          getOrganizations={getUserOrganizations}
        />
      )}
      {!detailsOpen && (
        <>
          {renderOrganizations()}
          {isOpen && (
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
                minWidth: isMobile ? '100%' : '34.4375rem',
                minHeight: isMobile ? '100%' : '22.1875rem'
              }}
              overlayClick={closeHandler}
              bigCloseImage={closeHandler}
              bigCloseImageStyle={{
                top: '-18px',
                right: '-18px',
                background: '#000',
                borderRadius: '50%'
              }}
            >
              <AddOrganization
                closeHandler={closeHandler}
                getUserOrganizations={getUserOrganizations}
                owner_pubkey={ui.meInfo?.owner_pubkey}
              />
            </Modal>
          )}
        </>
      )}
    </Container>
  );
};

export default Organizations;
