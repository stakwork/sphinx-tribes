import React, { useEffect, useState, useRef } from "react";
import { useStores } from "../../store";
import { useObserver } from "mobx-react-lite";
import Form from "../../form";
import ConfirmMe from "../confirmMe";
import type { MeInfo } from '../../store/ui'
import api from '../../api'
import styled, { css } from "styled-components";
import { getHostIncludingDockerHosts } from "../../host";
import { Button, IconButton, Modal } from "../../sphinxUI";
import moment from 'moment'
import SummaryViewer from '../widgetViews/summaryViewer'
import { useIsMobile } from "../../hooks";
import FocusedView from "./focusView";
import { aboutSchema } from "../../form/schema";
import FadeLeft from "../../animated/fadeLeft";

// this is where we see others posts (etc) and edit our own
export default function FirstTimeScreen() {
    const { ui } = useStores();

    const formHeader = <div style={{ marginTop: 60 }}>
        <Title>
            <B>Hi {ui.meInfo?.owner_alias},</B>
            <div>thank you for joining.</div>
        </Title>
        <SubTitle>
            Please, enter few basic info about yourself and create a public profile.
        </SubTitle>
    </div>

    return <Modal
        visible={true}
        style={{ height: '100%' }}
        envStyle={{ height: '100%', borderRadius: 0, width: '100%', maxWidth: 375 }}>
        <div style={{ height: '100%', padding: 20, paddingTop: 0 }}>
            <FocusedView
                formHeader={formHeader}
                buttonsOnBottom={true}
                person={ui.meInfo}
                canEdit={true}
                selectedIndex={-1}
                config={{
                    label: 'About',
                    name: 'about',
                    single: true,
                    skipEditLayer: true,
                    submitText: 'Submit',
                    schema: aboutSchema,
                }}
                onSuccess={() => {
                    console.log('success')
                }}

            />
        </div>
    </Modal>
}

const B = styled.div`
font-weight: bold;
`

const Title = styled.div`
font-size: 24px;
font-style: normal;

line-height: 30px;
letter-spacing: 0em;
text-align: center;
margin-bottom:20px;

`

const SubTitle = styled.div`
font-family: Roboto;
font-size: 15px;
font-style: normal;
font-weight: 400;
line-height: 20px;
letter-spacing: 0em;
text-align: center;

`
