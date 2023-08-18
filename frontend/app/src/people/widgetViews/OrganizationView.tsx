import React, { useState } from 'react';
import styled from 'styled-components';
import PageLoadSpinner from 'people/utils/PageLoadSpinner';
import NoResults from 'people/utils/OrgNoResults';

const Container = styled.div`
  display: flex;
  flex-flow: row wrap;
  gap: 1rem;
  min-width: 77vw;
  flex: 1 1 100%;
`;

const Organizations = () => {
    const [loading] = useState<boolean>(false);

    return (
        <div>
            <Container>
                <PageLoadSpinner show={loading} />
                <NoResults loading={loading} />
            </Container>
        </div>
    );
};

export default Organizations;
