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

    console.log('config.name', config.name)
    switch (config.name) {
        case 'post':
            return <PostSummary {...item} />
        case 'offer':
            return <OfferSummary {...item} />
        case 'wanted':
            return <WantedSummary {...item} />
        default:
            return <div>none</div>
    }
}
