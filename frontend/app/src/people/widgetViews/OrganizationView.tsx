import React, { useState, useEffect, useCallback } from 'react';
import styled from 'styled-components';
import PageLoadSpinner from 'people/utils/PageLoadSpinner';
import NoResults from 'people/utils/OrgNoResults';
import { useStores } from 'store';
import { Organization } from 'store/main';
import { EuiGlobalToastList } from '@elastic/eui';
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
`;

const OrganizationWrap = styled.div`
  display: flex;
  flex-direction: row;
  width: 100%;
  background: white;
  padding: 25px 30px;
  border-radius: 6px;
  @media only screen and (max-width: 800px) {
    padding: 15px 0px;
  }
  @media only screen and (max-width: 700px) {
    padding: 12px 0px;
    margin-bottom: 10px;
  }
  @media only screen and (max-width: 500px) {
    padding: 0px;
  }
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
  width: 65px;
  height: 65px;
  @media only screen and (max-width: 700px) {
    width: 55px;
    height: 55px;
  }
  @media only screen and (max-width: 500px) {
    width: 48px;
    height: 48px;
  }
  @media only screen and (max-width: 470px) {
    width: 60px;
    height: 60px;
  }
`;

const OrganizationContainer = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
  gap: 15px;
`;

const OrgHeadWrap = styled.div`
  display: flex;
  align-items: center;
  margin-top: 5px;
  margin-bottom: 20px;
`;

const OrgText = styled.div`
  font-size: 1.4rem;
  font-weight: bold;
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
  const [toasts, setToasts]: any = useState([]);
  const [user, setUser] = useState<Person>();
  const { main, ui } = useStores();
  const isMobile = useIsMobile();
  const config = widgetConfigs['organizations'];
  const isMyProfile = ui?.meInfo?.pubkey === props?.person?.owner_pubkey;

  // function addToast(title: string) {
  //   setToasts([
  //     {
  //       id: '1',
  //       title,
  //       color: 'danger'
  //     }
  //   ]);
  // }

  function removeToast() {
    setToasts([]);
  }

  const getUserOrganizations = useCallback(async () => {
    setIsLoading(true);
    if (ui.selectedPerson) {
      await main.getUserOrganizations(ui.selectedPerson);
      const user = await main.getPersonById(ui.selectedPerson);
      setUser(user);
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
          <OrganizationBudget org={org} user_pubkey={user?.owner_pubkey ?? ''} />
          <OrganizationActionWrap>
            {ui.meInfo?.owner_pubkey && (
              <ManageButton
                org={org}
                user_pubkey={user?.owner_pubkey ?? ''}
                action={() => {
                  setOrganization(org);
                  setDetailsOpen(true);
                }}
              />
            )}
            <Button
              disabled={btnDisabled}
              color={!btnDisabled ? 'white' : 'grey'}
              text="View Bounties"
              endingIcon="open_in_new"
              onClick={() => window.open(`/org/bounties/${org.uuid}`, '_target')}
              style={{
                height: 40,
                color: '#000000',
                borderRadius: 10
              }}
            />
          </OrganizationActionWrap>
        </OrganizationData>
      </OrganizationWrap>
    );
  };

  // renders list of orgs with header
  const renderOrganizations = () => {
    if (main.organizations.length) {
      return (
        <>
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
        </>
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
      <EuiGlobalToastList toasts={toasts} dismissToast={removeToast} toastLifeTimeMs={5000} />
    </Container>
  );
};

export default Organizations;
