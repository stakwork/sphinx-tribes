import React, { useEffect, useState, useRef } from "react";
import { useStores } from "../../store";
import { useObserver } from "mobx-react-lite";
import styled, { css } from "styled-components";
import { Button, IconButton } from "../../sphinxUI";
import moment from 'moment'
import PostSummary from './summaries/postSummary'
import WantedSummary from './summaries/wantedSummary'
import OfferSummary from './summaries/offerSummary'

// this is where we see others posts (etc) and edit our own
export default function SummaryViewer(props: any) {
    const { item, config, person } = props
    const { ui, main } = useStores();

    function wrapIt(child) {
        return <Wrap>
            {child}
        </Wrap>
    }

    switch (config.name) {
        case 'post':
            return wrapIt(<PostSummary {...item} />)
        case 'offer':
            return wrapIt(<OfferSummary {...item} />)
        case 'wanted':
            return wrapIt(<WantedSummary {...item} />)
        default:
            return wrapIt(<div>none</div>)
    }
}

const Wrap = styled.div`
height: calc(100% - 60px);
overflow: auto;
display: flex;
flex-direction:column;
width:100%;
min-width:100%;
padding:20px;
`;