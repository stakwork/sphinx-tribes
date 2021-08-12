import React, { useRef, useState } from "react";
import { QRCode } from "react-qr-svg";
import styled from "styled-components";
import { getHost } from "../host";
import qrCode from "../utils/invoice-qr-code.svg";
import { EuiCheckableCard, EuiButton, EuiButtonIcon, EuiToolTip } from "@elastic/eui";
import { useObserver } from 'mobx-react-lite'
import { useStores } from '../store'
import {
    EuiModal,
    EuiModalBody,
    EuiModalHeader,
    EuiModalHeaderTitle,
    EuiOverlayMask,
} from "@elastic/eui";

import { meSchema } from '../form/schema'

import BlogView from "./widgetViews/blogView";
import OfferView from "./widgetViews/offerView";
import TwitterView from "./widgetViews/twitterView";
import SupportMeView from "./widgetViews/supportMeView";
import WantedView from "./widgetViews/wantedView";
import FadeLeft from "../animated/fadeLeft";
import { useEffect } from "react";

const host = getHost();
function makeQR(pubkey: string) {
    return `sphinx.chat://?action=person&host=${host}&pubkey=${pubkey}`;
}

export default function PersonView(props: any) {

    const {
        personId,
        loading,
        goBack
    } = props

    const { main, ui } = useStores()

    const person = (main.people && main.people.length && main.people.find(f => f.id === personId))

    const {
        id,
        img,
        tags,
        description,
        owner_alias,
        unique_name,
        price_to_meet,
        extras
    } = person || {}

    const owner_pubkey = ''

    const [selectedWidget, setSelectedWidget] = useState('');
    const [newSelectedWidget, setNewSelectedWidget] = useState('');
    const [animating, setAnimating] = useState(false);
    const [showQR, setShowQR] = useState(false);
    const qrString = makeQR(owner_pubkey);

    useEffect(() => {
        if (extras && (Object.keys(extras).length > 0)) {
            const name = Object.keys(extras)[0]
            setSelectedWidget(name)
            setNewSelectedWidget(name)
        }
    }, [extras])

    function switchWidgets(name) {
        // setting newSelectedWidget will dismount the FadeLeft, 
        // and on dismount, endAnimation runs
        if (!animating && selectedWidget !== name) {
            setNewSelectedWidget(name)
            setAnimating(true)
        }
    }

    function endAnimation() {
        setSelectedWidget(newSelectedWidget)
        setAnimating(false)
    }

    let tagsString = "";
    tags && tags.forEach((t: string, i: number) => {
        if (i !== 0) tagsString += ",";
        tagsString += t;
    });

    function add(e) {
        e.stopPropagation();
    }
    function toggleQR(e) {
        e.stopPropagation();
        setShowQR((current) => !current);
    }

    if (loading) return <div>Loading...</div>

    let widgetSchemas: any = meSchema.find(f => f.name === 'extras')
    if (widgetSchemas && widgetSchemas.extras) {
        widgetSchemas = widgetSchemas && widgetSchemas.extras
    }

    const qrWidth = 209

    let fullSelectedWidget = (extras && selectedWidget) ? extras[selectedWidget] : {}

    // we do this because sometimes the widgets are empty arrays
    let filteredExtras = extras && { ...extras }
    if (filteredExtras) {
        let emptyArrayKeys = ['']

        Object.keys(filteredExtras).forEach(name => {
            const p = extras && extras[name]
            if (Array.isArray(p) && !p.length) {
                emptyArrayKeys.push(name)
            }
            const thisSchema = widgetSchemas && widgetSchemas.find(e => e.name === name)
            if (filteredExtras && thisSchema.single) {
                delete filteredExtras[name]
            }
        })

        emptyArrayKeys.forEach(e => {
            if (filteredExtras && e) delete filteredExtras[e]
        })
    }

    function renderSelectedWidget() {
        if (!selectedWidget || !fullSelectedWidget) {
            return <div style={{ height: 200 }} />
        }

        const widgetSchema: any = widgetSchemas && widgetSchemas.find(f => f.name === selectedWidget) || {}
        const single = widgetSchema.single
        let fields = [...widgetSchema.fields]
        // remove show from display
        fields = fields.filter(f => f.name !== 'show')

        function wrapIt(child) {
            if (single) {
                return null
            }

            return <FadeLeft
                direction={'up'}
                style={{ width: '100%' }}
                isMounted={newSelectedWidget === widgetSchema.name}
                speed={100}
                drift={20}
                dismountCallback={endAnimation}>
                <SelectedWidgetWrap>
                    {(fullSelectedWidget.length > 0) && fullSelectedWidget.map((s, i) => {
                        return <Card key={i} style={{ width: '100%' }}>
                            {React.cloneElement(child, { ...s })}
                        </Card>
                    })}
                </SelectedWidgetWrap>
            </FadeLeft>
        }

        switch (widgetSchema.name) {
            case 'twitter':
                return wrapIt(<TwitterView {...fullSelectedWidget} />)
            case 'supportme':
                return wrapIt(<SupportMeView {...fullSelectedWidget} />)
            case 'offer':
                return wrapIt(<OfferView {...fullSelectedWidget} />)
            case 'wanted':
                return wrapIt(<WantedView {...fullSelectedWidget} />)
            case 'blog':
                return wrapIt(<BlogView {...fullSelectedWidget} />)
            default:
                return <></>

        }
    }



    return (
        <Content>
            <div style={{ display: 'flex', justifyContent: 'space-between', width: '100%', padding: '0 20px', marginBottom: 40 }}>
                <EuiButton onClick={goBack}
                    iconType="sortLeft" aria-label="goback"
                />


                <div style={{ display: 'flex', width: 'fit-content', margin: 0 }}>

                    {extras && extras.supportme &&
                        <EuiToolTip position="bottom"
                            content={`Donate`}>
                            <EuiButtonIcon
                                style={{
                                    border: "1px solid #6B7A8D",
                                    color: "white",
                                    padding: 10,
                                    marginRight: 10
                                }}
                                iconType={'cheer'}
                                aria-label="donate"
                            />
                        </EuiToolTip>
                    }

                    <EuiToolTip position="bottom"
                        content={`Price to Meet: ${price_to_meet} sats`}>
                        <div style={{ display: 'flex', width: 'fit-content', margin: 0 }}>
                            <EuiButtonIcon
                                onClick={toggleQR}
                                style={{
                                    border: "1px solid #6B7A8D",
                                    color: "white",
                                    padding: 10,
                                    marginRight: 10
                                }}
                                iconType={qrCode}
                                aria-label="qr-code"
                            />

                            <a href={qrString}>
                                <EuiButton
                                    onClick={add}
                                    fill={true}
                                    style={{
                                        backgroundColor: "#6089ff",
                                        borderColor: "#6089ff",
                                        color: "white",
                                        fontWeight: 600,
                                        fontSize: 12,
                                        width: 80,
                                        maxWidth: 80,
                                        minWidth: 80
                                    }}
                                    aria-label="add"
                                >
                                    ADD
                                </EuiButton>
                            </a>
                        </div>
                    </EuiToolTip>
                </div>

            </div>
            {/* profile photo */}
            <Head>
                <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center' }}>
                    <Img src={img || '/static/sphinx.png'} />
                    <RowWrap>
                        <Name>{owner_alias}</Name>
                    </RowWrap>

                    {extras && extras.twitter &&
                        <RowWrap style={{ alignItems: 'center', margin: 0 }}>
                            <Icon source={'/static/twitter.png'} style={{ width: 14, height: 14, margin: '0 3px 0 0' }} />
                            <div style={{ fontSize: 14, color: '#ffffffd3' }}>{extras.twitter.handle}</div>
                        </RowWrap>
                    }

                    <RowWrap>
                        <Row style={{
                            padding: 10, maxWidth: 400, maxHeight: 400, margin: 10,
                            overflow: 'auto', background: '#ffffff21', borderRadius: 5
                        }}>
                            <Description>{description}</Description>
                        </Row>
                    </RowWrap>
                </div>
            </Head>

            <RowWrap>
                <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', width: '100%', margin: '0 20px', borderBottom: '1px solid #ffffff21' }}>
                    <TabRow>
                        {filteredExtras && Object.keys(filteredExtras).map((name, i) => {
                            const widgetSchema: any = widgetSchemas && widgetSchemas.find(f => f.name === name) || {}
                            const label = widgetSchema.label
                            const icon = widgetSchema.icon
                            const selected = name === newSelectedWidget

                            return < WidgetEnv key={i} selected={selected} onClick={() => switchWidgets(name)}>
                                <Widget key={i}>
                                    {label}
                                </Widget>
                            </WidgetEnv>
                        })}
                    </TabRow>

                    {/* <EuiButtonIcon
                        iconType={'arrowRight'}
                        aria-label="next"
                        style={{ width: 50, height: 50, margin: 0 }}
                    /> */}
                    <div />
                </div>
            </RowWrap>

            <RowWrap>
                <Row>
                    {renderSelectedWidget()}
                </Row>
            </RowWrap>

            {
                showQR &&
                <EuiOverlayMask onClick={() => setShowQR(false)}>
                    <EuiModal onClose={() => setShowQR(false)}
                        initialFocus="[name=popswitch]">
                        <EuiModalHeader>
                            <EuiModalHeaderTitle>{`Add ${owner_alias}`}</EuiModalHeaderTitle>
                        </EuiModalHeader>
                        <EuiModalBody style={{ padding: 0, color: '#fff' }}>
                            <RowWrap style={{ marginTop: -20, marginBottom: 10 }}>
                                {`Price to Meet: ${price_to_meet} sats`}
                            </RowWrap>
                            <QRWrapWrap>
                                <QRWrap className="qr-wrap float-r">
                                    <QRCode
                                        bgColor={"#FFFFFF"}
                                        fgColor="#000000"
                                        level="Q"
                                        style={{ width: qrWidth }}
                                        value={qrString}
                                    />
                                </QRWrap>
                            </QRWrapWrap>
                        </EuiModalBody>
                    </EuiModal>
                </EuiOverlayMask >
            }
        </Content >

    );
}
interface ContentProps {
    selected: boolean;
}
const Content = styled.div`
            display: flex;
            flex-direction:column;
            flex:1;
            width:100%;
            max-width:800px;
            align-items:center;
            color:#fff;
            padding-bottom:200px;
            `;
const QRWrapWrap = styled.div`
            display: flex;
            justify-content: center;
            `;
const QRWrap = styled.div`
            background: white;
            padding: 5px;
            `;
const Widget = styled.div`

            `;

const Head = styled.div`
            display:flex;
            flex-direction:column;
            justify-content:center;
            align-items:center;
            width:100%;
            `;

const Card = styled.div`
            min-width: 300px;
            max-width: 700px;
            padding: 20px;
            border: 1px solid #ffffff21;
            background:#ffffff07;
            border-radius: 5px;
            overflow:hidden;
            margin-bottom:20px;
            `;

const SupportMe = styled.div`
            min-width: 300px;
            max-width: 700px;
            padding: 20px;
            border: 1px solid #ffffff21;
            background:#ffffff07;
            border-radius: 5px;
            overflow:hidden;
            margin-bottom:20px;
            `;


const SelectedWidgetWrap = styled.div`
            display:flex;
            width:100%;
            justify-content:space-around;
            flex-wrap:wrap;
            `;
interface WidgetEnvProps {
    selected: boolean;
}
const WidgetEnv = styled.div<WidgetEnvProps>`
                display:flex;
                flex-direction:column;
                align-items:center;
                justify-content:center;
                padding:10px;
                min-width:80px;
                border-radius:5px;
                cursor:pointer;
                background:${p => p.selected && '#ffffff31'};
                &:hover{
                    background: ${p => !p.selected && '#ffffff21'};
            }
                `;
const Name = styled.div`
                color: white;
                font-weight: 500;
                margin-top:15px;
                `;
const Description = styled.div`
                color: white;
                font-weight: 340;
                `;
const Left = styled.div`
                height: 100%;
                max-width: 100%;
                display: flex;
                flex-direction: column;
                flex: 1;
                `;
const Row = styled.div`
                display: flex;
                align-items: center;
                width:100%;
                margin: 20px 0;
                justify-content: space-evenly;
                `;

const TabRow = styled.div`
                display: flex;
                flex-wrap:flex;
                align-items: center;
                width:100%;
                user-select:none;
                // margin: 10px 0;
                margin-top:10px;
                `;
const RowWrap = styled.div`
                display:flex;
                justify-content:center;

                width:100%`;
const Title = styled.h3`
                margin-right: 12px;
                font-size: 22px;
                white-space: nowrap;
                overflow: hidden;
                text-overflow: ellipsis;
                max-width: 100%;
                min-height: 24px;
                `;
interface DescriptionProps {
    oneLine: boolean;
}

interface ImageProps {
    readonly src: string;
}
const Img = styled.div<ImageProps>`
                    background-image: url("${(p) => p.src}");
                    background-position: center;
                    background-size: cover;
                    height: 90px;
                    width: 90px;
                    border-radius: 50%;
                    position: relative;
                    `;
const Tokens = styled.div`
                    display: flex;
                    align-items: center;
                    `;
const TagsWrap = styled.div`
                    display: flex;
                    flex-direction: row;
                    justify-content: flex-start;
                    align-items: center;
                    margin-top: 10px;
                    `;
const Tag = styled.h5`
                    margin-right: 10px;
                    `;
const Intro = styled.div`
                    color: white;
                    font-size: 14px;
                    margin: 10px;
                    `;
interface IconProps {
    source: string;
}

const Icon = styled.div<IconProps>`
                        background-image: ${p => `url(${p.source})`};
                        width:40px;
                        height:40px;
                        margin-top:10px;
                        background-position: center; /* Center the image */
                        background-repeat: no-repeat; /* Do not repeat the image */
                        background-size: contain; /* Resize the background image to cover the entire container */
                        border-radius:5px;
                        overflow:hidden;
                        `;