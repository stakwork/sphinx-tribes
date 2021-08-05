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

    const person = main.people && main.people[0]

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
    const [showQR, setShowQR] = useState(false);
    const qrString = makeQR(owner_pubkey);

    let tagsString = "";
    tags.forEach((t: string, i: number) => {
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

    let fullSelectedWidget = (person && person.extras && selectedWidget) ? person.extras[selectedWidget] : {}

    // we do this because sometimes the widgets are empty arrays
    let filteredExtras = person && { ...person.extras }
    if (filteredExtras) {
        let emptyArrayKeys = ['']

        Object.keys(filteredExtras).forEach(name => {
            const p = person.extras[name]
            if (Array.isArray(p) && !p.length) {
                emptyArrayKeys.push(name)
            }
        })

        emptyArrayKeys.forEach(e => {
            if (e) delete filteredExtras[e]
        })
    }

    function renderSelectedWidget() {
        if (!selectedWidget || !fullSelectedWidget) {
            return <div />
        }

        const p = fullSelectedWidget
        const widgetSchema: any = widgetSchemas && widgetSchemas.find(f => f.name === selectedWidget) || {}
        const single = widgetSchema.single
        const label = widgetSchema.label
        const icon = widgetSchema.icon
        let fields = [...widgetSchema.fields]
        // remove show from display
        fields = fields.filter(f => f.name !== 'show')

        if (single) {
            return <WidgetItem>
                {fields && fields.map((f, i) => {
                    return <div style={{ marginBottom: 5 }} key={i}>
                        <div>
                            {f.label}
                        </div>
                        <div>
                            {fullSelectedWidget[f.name]}
                        </div>
                    </div>
                })}
            </WidgetItem>
        }

        return <SelectedWidgetWrap>
            {fullSelectedWidget.map((s, si) => {
                return <WidgetItem key={si}>
                    {fields && fields.map((f, i) => {
                        return <div style={{ marginBottom: 5 }} key={i}>
                            <div>
                                {f.label}
                            </div>
                            <div>
                                {s[f.name]}
                            </div>
                        </div>
                    })}
                </WidgetItem>
            })}

        </SelectedWidgetWrap>

        // do this instead, build each widget its own view style
        switch (widgetSchema.type) {
            case 'twitter':
                return <TwitterView {...props} />
            case 'supportme':
                return <SupportMeView {...props} />
            case 'offers':
                return <OfferView {...props} />
            case 'wanted':
                return <WantedView {...props} />
            case 'blog':
                return <BlogView {...props} />
            case 'hidden':
                return <></>
            default:
                return <></>

        }
    }

    return (
        <Content>
            <div style={{ display: 'flex', justifyContent: 'space-between', width: '100%', padding: '0 20px' }}>
                <EuiButton onClick={goBack}
                    iconType="sortLeft" aria-label="goback"
                />

                <EuiToolTip position="bottom"
                    content={`Price to Meet: ${price_to_meet} sats`}>
                    <Row style={{ width: 'fit-content', margin: 0 }}>

                        <EuiButtonIcon
                            onClick={toggleQR}
                            style={{
                                border: "1px solid #6B7A8D",
                                color: "white",
                                padding: 10,
                                marginRight: 10
                            }}
                            iconType={qrCode}
                        // aria-label="qr-code"
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
                                }}
                            // aria-label="add"
                            >
                                ADD
                            </EuiButton>
                        </a>
                    </Row>
                </EuiToolTip>
            </div>
            {/* profile photo */}
            <Img src={img || '/static/sphinx.png'} />
            <RowWrap>
                <Name>{owner_alias}</Name>
            </RowWrap>
            <RowWrap>
                <Row style={{ padding: 20, maxWidth: 300, background: '#ffffff21', borderRadius: 5 }}>
                    <Description>{description}</Description>
                </Row>
            </RowWrap>

            <RowWrap>
                <Row style={{ flexWrap: 'wrap' }}>
                    {filteredExtras && Object.keys(filteredExtras).map((name, i) => {
                        const p = person.extras[name]
                        const widgetSchema: any = widgetSchemas && widgetSchemas.find(f => f.name === name) || {}
                        const label = widgetSchema.label
                        const icon = widgetSchema.icon
                        const selected = name === selectedWidget
                        console.log('p', p)

                        return < WidgetEnv key={i} selected={selected} onClick={() => setSelectedWidget(name)}>
                            {
                                widgetSchema.single ?
                                    <Widget key={i}>
                                        {label}
                                    </Widget>

                                    : <Widget key={i}>{label}</Widget>
                            }
                            <Icon source={`/static/${icon || 'sphinx'}.png`} />
                        </WidgetEnv>
                    })}
                </Row>
            </RowWrap>

            <RowWrap>
                <Row>
                    {fullSelectedWidget &&
                        <>
                            {renderSelectedWidget()}
                        </>
                    }
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

const WidgetItem = styled.div`
    width: 300px;
    padding: 20px;
    border: 1px solid #fff;
    border-radius: 5px;
    overflow:hidden;
            `;



const SelectedWidgetWrap = styled.div`
            display:flex;
            width:100%;
            justify-content:space-around;
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
            padding-bottom:20px;
            min-width:120px;
            border-radius:5px;
            cursor:pointer;
            background:${p => p.selected && '#ffffff41'};
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