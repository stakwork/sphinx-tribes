import React, { useRef, useState, useEffect } from "react";
import { QRCode } from "react-qr-svg";
import styled from "styled-components";
import { getHost } from "../host";
import qrCode from "../utils/invoice-qr-code.svg";
import { useObserver } from 'mobx-react-lite'
import { useStores } from '../store'
import {
    EuiModal,
    EuiModalBody,
    EuiModalHeader,
    EuiModalHeaderTitle,
    EuiOverlayMask,
} from "@elastic/eui";

import AboutView from "./widgetViews/aboutView";
import BlogView from "./widgetViews/blogView";
import OfferView from "./widgetViews/offerView";
import TwitterView from "./widgetViews/twitterView";
import SupportMeView from "./widgetViews/supportMeView";
import WantedView from "./widgetViews/wantedView";
import PostView from "./widgetViews/postView";

import FadeLeft from "../animated/fadeLeft";
import { Button, IconButton, Modal } from "../sphinxUI";
import MaterialIcon from "@material/react-material-icon";
import FocusedView from './mobile/focusView'
import { aboutSchema, postSchema, wantedSchema, meSchema, offerSchema } from "../form/schema";
import { useIsMobile } from "../hooks";

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
    const { meInfo } = ui || {}

    const person = (main.people && main.people.length && main.people.find(f => f.id === personId))

    const {
        id,
        img,
        tags,
        description,
        owner_alias,
        unique_name,
        price_to_meet,
        extras,
        owner_pubkey
    } = person || {}


    const canEdit = id === meInfo?.id


    const [selectedWidget, setSelectedWidget] = useState(canEdit ? 'post' : 'about');
    const [focusIndex, setFocusIndex] = useState(-1);
    const [newSelectedWidget, setNewSelectedWidget] = useState(canEdit ? 'post' : 'about');
    const [animating, setAnimating] = useState(false);
    const [showQR, setShowQR] = useState(false);
    const [showFocusView, setShowFocusView] = useState(false);
    const qrString = makeQR(owner_pubkey || '');

    const isMobile = useIsMobile()

    function switchWidgets(name) {
        // setting newSelectedWidget will dismount the FadeLeft, 
        // and on dismount, endAnimation runs
        // if (!animating && selectedWidget !== name) {
        setNewSelectedWidget(name)
        setSelectedWidget(name)
        setAnimating(true)
        // }
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

    function logout() {
        ui.setEditMe(false)
        ui.setChallenge('')
        ui.setMeInfo(null)
        goBack()
    }

    if (loading) return <div>Loading...</div>

    let widgetSchemas: any = meSchema.find(f => f.name === 'extras')
    if (widgetSchemas && widgetSchemas.extras) {
        widgetSchemas = widgetSchemas && widgetSchemas.extras
    }

    const qrWidth = 209

    let fullSelectedWidget: any = (extras && selectedWidget) ? extras[selectedWidget] : null

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
            if (filteredExtras && thisSchema && thisSchema.single) {
                delete filteredExtras[name]
            }
        })

        emptyArrayKeys.forEach(e => {
            if (filteredExtras && e) delete filteredExtras[e]
        })
    }


    const tabs = {
        about: {
            label: 'About',
            name: 'about',
            single: true,
            skipEditLayer: true,
            submitText: 'Save',
            schema: aboutSchema,
            action: {
                text: 'Edit Profile',
                icon: 'edit'
            }
        },
        post: {
            label: 'Posts',
            name: 'post',
            submitText: 'Post',
            schema: postSchema,
            action: {
                text: 'Create a Post',
                icon: 'add',
                info: "What's on your mind?",
                infoIcon: 'chat_bubble_outline'
            }
        },
        offer: {
            label: 'Offer',
            name: 'offer',
            submitText: 'Post',
            schema: offerSchema,
            action: {
                text: 'Sell Something',
                icon: 'local_offer'
            }
        },
        wanted: {
            label: 'Wanted',
            name: 'wanted',
            submitText: 'Save',
            schema: wantedSchema,
            action: {
                text: 'Add to Wanted',
                icon: 'favorite_outline'
            }
        },
    }

    function renderWidgets() {
        if (!selectedWidget) {
            return <div style={{ height: 200 }} />
        }

        const widgetSchema: any = widgetSchemas && widgetSchemas.find(f => f.name === selectedWidget) || {}
        const single = widgetSchema.single
        let fields = widgetSchema.fields && [...widgetSchema.fields]
        // remove show from display
        fields = fields && fields.filter(f => f.name !== 'show')

        function wrapIt(child) {
            if (single) {
                return <Panel>
                    {child}
                </Panel>
            }

            if (!fullSelectedWidget) return <div />

            const elementArray: any = []

            fullSelectedWidget && fullSelectedWidget.map((s, i) => {

                elementArray.push(<Panel key={i}
                    onClick={() => {
                        setShowFocusView(true)
                        setFocusIndex(i)
                    }}
                    style={{ width: '100%' }}>
                    {React.cloneElement(child, { ...s })}
                </Panel>)
            })

            // </Panel>
            return elementArray
        }

        switch (selectedWidget) {
            case 'about':
                return <Panel>
                    <AboutView {...person} />
                </Panel>
            case 'post':
                return wrapIt(<PostView {...fullSelectedWidget} />)
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
                return wrapIt(<></>)

        }
    }

    function renderEditButton() {
        if (!canEdit || !selectedWidget) return <div />

        let { action } = tabs[selectedWidget] || {}
        action = action || {}
        return <div style={{ padding: 10, margin: '8px 0 5px' }}>
            {!fullSelectedWidget && action.info &&
                <ActionInfo>
                    <MaterialIcon icon={action.infoIcon} style={{ fontSize: 80 }} />
                    <>{action.info}</>
                </ActionInfo>
            }
            <Button
                text={action.text}
                color='widget'
                leadingIcon={action.icon}
                width='100%'
                height={48}
                onClick={() => {
                    setShowFocusView(true)
                }}
            />
        </div>
    }



    return (
        <Content>
            <div style={{
                display: 'flex', flexDirection: 'column',
                width: '100%', overflow: 'auto', height: '100%'
            }}>
                <Panel style={{ paddingBottom: 0, paddingTop: 80 }}>
                    <div style={{
                        position: 'absolute',
                        top: 20, left: 0,
                        display: 'flex',
                        justifyContent: 'space-between', width: '100%',
                    }}>
                        <IconButton
                            onClick={goBack}
                            icon='arrow_back'
                        />
                        {canEdit ?
                            <IconButton
                                onClick={logout}
                                icon='logout'
                            /> : <div />
                        }
                    </div>

                    {/* profile photo */}
                    <Head>
                        <Img src={img || '/static/sphinx.png'} />
                        <RowWrap>
                            <Name>{owner_alias}</Name>
                        </RowWrap>

                        {/* only see buttons on other people's profile */}
                        {canEdit ? <div style={{ height: 40 }} /> :
                            <RowWrap style={{ marginBottom: 30, marginTop: 25 }}>
                                {isMobile ?
                                    <a href={qrString}>
                                        <Button
                                            text='Connect'
                                            onClick={add}
                                            color='primary'
                                            height={42}
                                            width={120}
                                        />
                                    </a>
                                    :
                                    <Button
                                        text='Connect'
                                        onClick={() => setShowQR(true)}
                                        color='primary'
                                        height={42}
                                        width={120}
                                    />
                                }
                                <div style={{ width: 15 }} />
                                <Button
                                    text='Support'
                                    color='link'
                                    height={42}
                                    width={120} />
                            </RowWrap>
                        }
                    </Head>

                    <Tabs>
                        {tabs && Object.keys(tabs).map((name, i) => {
                            const t = tabs[name]
                            const label = t.label
                            const selected = name === newSelectedWidget

                            return <Tab key={i}
                                selected={selected}
                                onClick={() => {
                                    switchWidgets(name)
                                }}>
                                {label}
                            </Tab>
                        })}

                    </Tabs>

                </Panel>

                <Sleeve>
                    {renderEditButton()}
                    {renderWidgets()}
                </Sleeve>

                {
                    showQR &&
                    <EuiOverlayMask onClick={() => setShowQR(false)}>
                        <EuiModal onClose={() => setShowQR(false)}
                            initialFocus="[name=popswitch]">
                            <EuiModalHeader>
                                <EuiModalHeaderTitle>{`Connect with ${owner_alias}`}</EuiModalHeaderTitle>
                            </EuiModalHeader>
                            <EuiModalBody style={{ padding: 0, textAlign: 'center' }}>
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
                                <div style={{ marginTop: 10, color: '#fff' }}>Scan with your Sphinx Mobile App</div>
                            </EuiModalBody>
                        </EuiModal>
                    </EuiOverlayMask >
                }
            </div>

            <Modal
                fill
                visible={showFocusView}>
                <FocusedView
                    person={person}
                    canEdit={canEdit}
                    selectedIndex={focusIndex}
                    config={tabs[selectedWidget] && tabs[selectedWidget]}
                    onSuccess={() => {
                        console.log('success')
                        setFocusIndex(-1)
                    }}
                    goBack={() => {
                        setShowFocusView(false)
                        setFocusIndex(-1)
                    }}
                />
            </Modal>


            {/* <Bottom>
                <a href={qrString}>
                    <EuiButton
                        onClick={add}
                        fill={true}
                        style={{
                            backgroundColor: "#6089ff",
                            borderColor: "#6089ff",
                            fontWeight: 600,
                            fontSize: 12,
                            width: '40%',
                            minWidth: 140,
                            maxWidth: 140,
                            color: '#fff'
                        }}
                        aria-label="join"
                    >
                        JOIN
                    </EuiButton>
                </a>
                <div style={{ width: 20 }} />
                <EuiButton
                    onClick={toggleQR}
                    fill={true}
                    style={{
                        background: "#fff",
                        borderColor: "#5F6368",
                        color: '#5F6368',
                        fontWeight: 600,
                        fontSize: 12,
                        width: '40%',
                        minWidth: 140,
                        maxWidth: 140
                    }}
                    // iconType={qrCode}
                    aria-label="qr-code"
                >
                    QR
                </EuiButton>
            </Bottom> */}
        </Content >

    );
}
interface ContentProps {
    selected: boolean;
}

const Panel = styled.div`
            position:relative;
            background:#ffffff;
            color:#000000;
            margin-bottom:10px;
            padding:20px;
            box-shadow:0px 0px 3px rgb(0 0 0 / 29%);
            `;
const Content = styled.div`
            display: flex;
            flex-direction:column;

            width:100%;
            height: 100%;
            align-items:center;
            color:#000000;
            background:#f0f1f3;
            `;
const QRWrapWrap = styled.div`
            display: flex;
            justify-content: center;
            `;
const QRWrap = styled.div`
            background: white;
            padding: 5px;
            `;

const ActionInfo = styled.div`
            font-style: normal;
            font-weight: normal;
            font-size: 22px;
            line-height: 26px;
            display: flex;
            align-items: center;
            text-align: center;
            padding: 20px;
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            color:#B0B7BC;
            margin-bottom:10px;
            `;


/* Placeholder Text */

const Tabs = styled.div`
            display:flex;
            width:100%;
            overflow-x:auto;
            ::-webkit-scrollbar {
                display: none;
            }
            `;

interface TagProps {
    selected: boolean;
}
const Tab = styled.div<TagProps>`
                display:flex;
                padding:10px;
                margin-right:25px;
                color:${p => p.selected ? '#292C33' : '#8E969C'};
                border-bottom: ${p => p.selected && '4px solid #618AFF'};
                cursor:hover;
                font-weight: 500;
                font-size: 16px;
                line-height: 19px;
                cursor:pointer;
                `;



const Bottom = styled.div`
                height:80px;
                width:100%;
                display:flex;
                justify-content:center;
                align-items:center;
                background: #FFFFFF;
                box-shadow: 0px -2px 4px rgba(0, 0, 0, 0.1);
                `;
const Head = styled.div`
                display:flex;
                flex-direction:column;
                justify-content:center;
                align-items:center;
                width:100%;
                `;

const B = styled.span`
                color:#000;
                font-weight:bold;
                margin-right:5px;
                `;

const Card = styled.div`
                margin-bottom:10px;
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
                    font-style: normal;
                    font-weight: 500;
                    font-size: 30px;
                    line-height: 19px;
                    /* or 73% */

                    text-align: center;

                    /* Text 2 */

                    color: #3C3F41;
                    `;


const Sleeve = styled.div`

                    `;

const Description = styled.div`
                    font-family: Roboto;
                    font-style: normal;
                    font-weight: normal;
                    font-size: 13px;
                    line-height: 19px;
                    /* or 146% */


                    /* Secondary Text 4 */

                    color: #8E969C;
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
                        margin-bottom:20px;
                        width:150px;
                        height:150px;
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
                        font-size: 14px;
                        margin: 10px;
                        `;
interface IconProps {
    source: string;
}

const Icon = styled.div<IconProps>`
                            background-image: ${p => `url(${p.source})`};
                            width:150px;
                            height:150px;
                            margin-top:10px;
                            background-position: center; /* Center the image */
                            background-repeat: no-repeat; /* Do not repeat the image */
                            background-size: contain; /* Resize the background image to cover the entire container */
                            border-radius:5px;
                            overflow:hidden;
                            `;