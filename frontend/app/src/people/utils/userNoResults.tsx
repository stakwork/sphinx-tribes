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

    const tabs = widgetConfigs;

    if (props.loading) {
        return <PageLoadSpinner show={true} />;
    } else {
        return (
            <NoneSpace
                small
                style={{
                    minWidth: '60vw',
                    minHeight: '90vh',
                }}
                {...tabs['usertickets']?.noneSpace['noResult']}
            />
        )
    }
}