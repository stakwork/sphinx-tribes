import React from 'react';
import styled from 'styled-components';
import { useStores } from '../../store';
import PageLoadSpinner from './pageLoadSpinner';
import { observer } from 'mobx-react-lite';
import NoneSpace from './noneSpace';
import { widgetConfigs } from '../utils/constants';

export default observer(NoResults);
function NoResults(props) {
    const { ui } = useStores();
    const { searchText } = ui || {};

    const tabs = widgetConfigs;

    if (props.loading) {
        return <PageLoadSpinner show={true} />;
    } else {
        return (
            <NoneSpace
                small
                style={{
                    margin: 'auto',
                    marginTop: '25%'
                }}
                {...tabs['usertickets']?.noneSpace['noResult']}
            />
        )
    }
}

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
