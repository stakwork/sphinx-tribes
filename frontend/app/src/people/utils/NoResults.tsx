import React from 'react';
import styled from 'styled-components';
import { observer } from 'mobx-react-lite';
import { useStores } from '../../store';
import PageLoadSpinner from './PageLoadSpinner';

const H = styled.div`
  font-size: 16px;
  font-style: normal;
  font-weight: 500;
  line-height: 37px;
  letter-spacing: 0.1em;
  text-align: center;

  font-family: Roboto;
  font-style: normal;
  line-height: 26px;
  display: flex;
  align-items: center;
  text-align: center;

  /* Primary Text 1 */

  color: #292c33;
  letter-spacing: 0px;
  color: rgb(60, 63, 65);
`;
function NoResults() {
  const { ui } = useStores();
  const { searchText } = ui || {};

  if (searchText) {
    return (
      <div
        style={{
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'center',
          alignItems: 'center',
          width: '100%',
          marginTop: 20
        }}
      >
        <H>No results</H>
      </div>
    );
  } else {
    return <PageLoadSpinner show={true} />;
  }
}

export default observer(NoResults);
