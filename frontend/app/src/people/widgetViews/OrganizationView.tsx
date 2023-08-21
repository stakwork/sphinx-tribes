import React, { useState, useEffect } from 'react';
import styled from 'styled-components';
import PageLoadSpinner from 'people/utils/PageLoadSpinner';
import NoResults from 'people/utils/OrgNoResults';
import { useStores } from 'store';
import { Organization } from 'store/main';
import { IconButton } from 'components/common';
import { useIsMobile } from 'hooks/uiHooks';
import avatarIcon from '../../public/static/profile_avatar.svg';

const Container = styled.div`
  display: flex;
  flex-flow: column wrap;
  gap: 1rem;
  min-width: 77vw;
  flex: 1 1 100%;
`;

const OrganizationText = styled.p`
  font-size: 1rem;
  text-transform: capitalize;
  font-weight: bold;
  margin-top: 15px;
`;

const OrganizationImg = styled.img`
  width: 60px;
  height: 60px;
`;

const OrganizationWrap = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  width: calc(19% - 40px);
  margin-left: 0.5%;
  margin-right: 0.5%;
  margin; 10px;
  background: white;
  padding: 20px;
  border-radius: 2px;
`;

const Organizations = () => {
    const [loading, setIsLoading] = useState<boolean>(false);
    const { main } = useStores();
    const isMobile = useIsMobile();

    async function getUserOrganizations() {
        setIsLoading(true);
        await main.getUserOrganizations();
        setIsLoading(false);
    }

    useEffect(() => {
        getUserOrganizations();
    }, []);


    const renderOrganizations = () => {
        if (main.organizations.length) {
            return main.organizations.map((org: Organization, i: number) => (

                <OrganizationWrap key={i}>
                    <OrganizationImg src={org.img || avatarIcon} />
                    <OrganizationText>{org.name}</OrganizationText>
                </OrganizationWrap>
            ))
        } else {
            return <NoResults loading={loading} />
        }
    }

    return (
        <div>
            <Container>
                <PageLoadSpinner show={loading} />
                <IconButton
                    width={150}
                    height={isMobile ? 36 : 48}
                    text="Add Organization"
                    onClick={() => false}
                    style={{
                        marginLeft: '10px'
                    }}
                />
                {renderOrganizations()}
            </Container>
        </div>
    );
};

export default Organizations;
