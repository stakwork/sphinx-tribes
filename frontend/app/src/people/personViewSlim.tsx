import React, { useState, useEffect } from "react";
import styled from "styled-components";
import { getHost } from "../host";
import { useStores } from '../store'

import AboutView from "./widgetViews/aboutView";
import BlogView from "./widgetViews/blogView";
import OfferView from "./widgetViews/offerView";
import TwitterView from "./widgetViews/twitterView";
import SupportMeView from "./widgetViews/supportMeView";
import WantedView from "./widgetViews/wantedView";
import PostView from "./widgetViews/postView";

import { Button, IconButton, Modal } from "../sphinxUI";
import MaterialIcon from "@material/react-material-icon";
import FocusedView from './main/focusView'
import { meSchema } from "../form/schema";
import { useIsMobile } from "../hooks";
import Person from "./person";
import NoneSpace from "./utils/noneSpace";
import ConnectCard from "./utils/connectCard";
import { widgetConfigs } from "./utils/constants";

const host = getHost();
function makeQR(pubkey: string) {
    return `sphinx.chat://?action=person&host=${host}&pubkey=${pubkey}`;
}

export default function PersonView(props: any) {

    const {
        personId,
        loading,
        selectPerson,
        goBack,
    } = props

    const { main, ui } = useStores()
    const { meInfo } = ui || {}

    // FOR PEOPLE VIEW
    let person: any = (main.people && main.people.length && main.people.find(f => f.id === personId))

    // if i select myself, fill person with meInfo
    if (personId === ui.meInfo?.id) {
        console.log('before', person.twitter_confirmed)
        console.log('after', ui.meInfo?.twitter_confirmed)
        person = { twitter_confirmed: person.twitter_confirmed, ...ui.meInfo }

    }

    const people = (main.people && main.people.filter(f => !f.hide)) || []

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
    const isMobile = useIsMobile()

    const initialWidget = (!isMobile || canEdit) ? 'post' : 'about'

    const [selectedWidget, setSelectedWidget] = useState(initialWidget);
    const [newSelectedWidget, setNewSelectedWidget] = useState(initialWidget);
    const [focusIndex, setFocusIndex] = useState(-1);
    const [showSupport, setShowSupport] = useState(false);

    const [showQR, setShowQR] = useState(false);
    const [showFocusView, setShowFocusView] = useState(false);
    const qrString = makeQR(owner_pubkey || '');



    function switchWidgets(name) {
        setNewSelectedWidget(name)
        setSelectedWidget(name)
        setShowFocusView(false)
        setFocusIndex(-1)
    }

    function selectPersonWithinFocusView(id, unique_name) {
        setShowFocusView(false)
        setFocusIndex(-1)
        selectPerson(id, unique_name)
    }

    useEffect(() => {
        if (ui.personViewOpenTab) {
            switchWidgets(ui.personViewOpenTab)
            ui.setPersonViewOpenTab('')
        }
    }, [ui.personViewOpenTab])


    let tagsString = "";
    tags && tags.forEach((t: string, i: number) => {
        if (i !== 0) tagsString += ",";
        tagsString += t;
    });

    function add(e) {
        e.stopPropagation();
    }

    function logout() {
        ui.setEditMe(false)
        ui.setMeInfo(null)
        main.getPeople()
        goBack()
    }

    if (loading) return <div>Loading...</div>

    let widgetSchemas: any = meSchema.find(f => f.name === 'extras')
    if (widgetSchemas && widgetSchemas.extras) {
        widgetSchemas = widgetSchemas && widgetSchemas.extras
    }

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

    const tabs = widgetConfigs

    function hasWidgets() {
        let has = false
        if (fullSelectedWidget && fullSelectedWidget.length) {
            has = true
        }
        return has
    }

    function renderWidgets(name: string) {
        if (name) {
            switch (name) {
                case 'about':
                    return <AboutView canEdit={canEdit} {...person} />
                case 'post':
                    return wrapIt(<PostView {...fullSelectedWidget} person={person} />)
                case 'twitter':
                    return wrapIt(<TwitterView {...fullSelectedWidget} person={person} />)
                case 'supportme':
                    return wrapIt(<SupportMeView {...fullSelectedWidget} person={person} />)
                case 'offer':
                    return wrapIt(<OfferView {...fullSelectedWidget} person={person} />)
                case 'wanted':
                    return wrapIt(<WantedView {...fullSelectedWidget} person={person} />)
                case 'blog':
                    return wrapIt(<BlogView {...fullSelectedWidget} person={person} />)
                default:
                    return wrapIt(<></>)
            }
        }
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

            const elementArray: any = []

            const panelStyles = isMobile ? {
                minHeight: 132
            } : {
                maxWidth: 291, minWidth: 291,
                marginRight: 20, marginBottom: 20, minHeight: 472
            }

            fullSelectedWidget && fullSelectedWidget.forEach((s, i) => {

                elementArray.push(<Panel key={i}
                    onClick={() => {
                        setShowFocusView(true)
                        setFocusIndex(i)
                    }}
                    style={{
                        ...panelStyles,
                        cursor: 'pointer',
                        padding: 0, overflow: 'hidden'
                    }}
                >
                    {React.cloneElement(child, { ...s })}
                </Panel>)
            })

            // if empty
            if (!elementArray.length) {
                const noneKey = canEdit ? 'me' : 'otherUser'
                return <NoneSpace
                    action={() => setShowFocusView(true)}
                    {...tabs[selectedWidget].noneSpace[noneKey]}
                />
            }

            return elementArray
        }

        switch (selectedWidget) {
            case 'about':
                return <Panel>
                    <AboutView {...person} />
                </Panel>
            case 'post':
                return wrapIt(<PostView {...fullSelectedWidget} person={person} />)
            case 'twitter':
                return wrapIt(<TwitterView {...fullSelectedWidget} person={person} />)
            case 'supportme':
                return wrapIt(<SupportMeView {...fullSelectedWidget} person={person} />)
            case 'offer':
                return wrapIt(<OfferView {...fullSelectedWidget} person={person} />)
            case 'wanted':
                return wrapIt(<WantedView {...fullSelectedWidget} person={person} />)
            case 'blog':
                return wrapIt(<BlogView {...fullSelectedWidget} person={person} />)
            default:
                return wrapIt(<></>)
        }
    }

    function nextIndex() {
        if (focusIndex < 0) {
            console.log('nope!')
            return
        }
        if (person && person.extras) {
            let g = person.extras[tabs[selectedWidget].name]
            let nextindex = focusIndex + 1
            if (g[nextindex]) setFocusIndex(nextindex)
            else setFocusIndex(0)
        }
    }

    function prevIndex() {
        if (focusIndex < 0) {
            console.log('nope!')
            return
        }
        if (person && person.extras) {
            let g = person?.extras[tabs[selectedWidget].name]
            let previndex = focusIndex - 1
            if (g[previndex]) setFocusIndex(previndex)
            else setFocusIndex(g.length - 1)
        }
    }

    function renderEditButton(style: any) {
        if (!canEdit || !selectedWidget) return <div />
        // don't return button if there are no items in list, the button is returned elsewhere

        if (selectedWidget !== 'about') {
            if (!fullSelectedWidget || (fullSelectedWidget && fullSelectedWidget.length < 1)) return <div />
        }

        let { action } = tabs[selectedWidget] || {}
        action = action || {}
        return <div style={{ padding: isMobile ? 10 : '10 0', margin: isMobile ? '6px 0 5px' : '10px 0', ...style }}>
            {!fullSelectedWidget && action.info &&
                <ActionInfo>
                    <MaterialIcon icon={action.infoIcon} style={{ fontSize: 80 }} />
                    <>{action.info}</>
                </ActionInfo>
            }
            <Button
                text={action.text}
                color={isMobile ? 'widget' : 'desktopWidget'}
                leadingIcon={action.icon}
                // width={isMobile ? '100%' : 291}
                width={'100%'}
                height={48}
                onClick={() => {
                    setShowFocusView(true)
                }}
            />
        </div >
    }


    function renderMobileView() {
        return <div style={{
            display: 'flex', flexDirection: 'column',
            width: '100%', overflow: 'auto', height: '100%'
        }}>
            <Panel style={{ paddingBottom: 0, paddingTop: 80 }}>
                <div style={{
                    position: 'absolute',
                    top: 20, left: 0,
                    display: 'flex',
                    justifyContent: 'space-between', width: '100%',
                    padding: '0 20px'
                }}>
                    <IconButton
                        onClick={goBack}
                        icon='arrow_back'
                    />
                    {canEdit ?
                        <IconButton
                            iconStyle={{ transform: 'rotate(270deg)' }}
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
                            <a href={qrString}>
                                <Button
                                    text='Connect'
                                    onClick={add}
                                    color='primary'
                                    height={42}
                                    width={120}
                                />
                            </a>

                            <div style={{ width: 15 }} />

                            <Button
                                text='Send Tip'
                                color='link'
                                height={42}
                                width={120}
                                onClick={() => setShowSupport(true)} />
                        </RowWrap>
                    }
                </Head>

                <Tabs>
                    {tabs && Object.keys(tabs).map((name, i) => {
                        const t = tabs[name]
                        const label = t.label
                        const selected = name === newSelectedWidget
                        let count = (extras[name] && extras[name].length > 0) ? extras[name].length : null

                        return <Tab key={i}
                            selected={selected}
                            onClick={() => {
                                switchWidgets(name)
                            }}>
                            {label}
                            {count && <Counter>
                                {count}
                            </Counter>}
                        </Tab>
                    })}

                </Tabs>

            </Panel>

            <Sleeve>
                {renderEditButton({})}
                {renderWidgets('')}
                <div style={{ height: 60 }} />
            </Sleeve>


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
        </div>
    }

    function renderDesktopView() {
        const focusedDesktopModalStyles = newSelectedWidget ? {
            ...tabs[newSelectedWidget].modalStyle
        } : {}

        return <div style={{
            display: 'flex',
            width: '100%', height: '100%'
        }}>

            {!canEdit &&
                <PeopleList>
                    <DBack >
                        <Button
                            color='clear'
                            leadingIcon='arrow_back'
                            text='Back'
                            onClick={goBack}
                        />
                    </DBack>

                    <div style={{ width: '100%', overflowY: 'auto' }} >
                        {people.map(t => <Person {...t} key={t.id}
                            selected={personId === t.id}
                            hideActions={true}
                            small={true}
                            select={selectPersonWithinFocusView}
                        />)}
                    </div>

                </PeopleList>
            }

            <div style={{
                width: 364,
                minWidth: 364, overflowY: 'auto',
                position: 'relative',
                background: '#ffffff',
                color: '#000000',
                padding: 40,
                height: '100%',
                borderLeft: '1px solid #F2F3F5',
                borderRight: '1px solid #F2F3F5',
                boxShadow: '1px 0px 6px -2px rgba(0, 0, 0, 0.07)'
            }}>

                {canEdit && <div style={{
                    position: 'absolute',
                    top: 0, left: 0,
                    display: 'flex',
                    justifyContent: 'space-between', width: '100%',
                    alignItems: 'center',
                    boxShadow: '0px 1px 6px rgba(0, 0, 0, 0.07)',
                    paddingRight: 10,
                    height: 64,
                    zIndex: 0
                }}>
                    <Button
                        color='clear'
                        leadingIcon='arrow_back'
                        text='Back'
                        onClick={goBack}
                    />
                    <div />
                </div>
                }

                {/* profile photo */}
                <Head>
                    <div style={{ height: canEdit ? 80 : 35 }} />

                    <Img src={img || '/static/sphinx.png'} >
                        <IconButton
                            iconStyle={{ color: '#5F6368' }}
                            style={{
                                zIndex: 2, width: '40px',
                                height: '40px', padding: 0,
                                background: '#ffffff', border: '1px solid #D0D5D8',
                                boxSizing: 'borderBox', borderRadius: 4
                            }}
                            icon={'qr_code_2'}
                            onClick={() => setShowQR(true)}
                        />
                    </Img>


                    <RowWrap>
                        <Name>{owner_alias}</Name>
                    </RowWrap>

                    {/* only see buttons on other people's profile */}
                    {canEdit ? <RowWrap style={{ marginBottom: 30, marginTop: 25, justifyContent: 'space-around' }}>
                        <Button
                            text='Edit Profile'
                            onClick={() => {
                                switchWidgets('about')
                                setShowFocusView(true)
                            }}
                            color='widget'
                            height={42}
                            style={{ fontSize: 13 }}

                            leadingIcon={'edit'}
                            iconSize={15}
                        />
                        <Button
                            text='Sign out'
                            onClick={logout}
                            height={42}
                            style={{ fontSize: 13 }}
                            iconSize={15}
                            color='white'
                            leadingIcon='logout'
                        />

                    </RowWrap>
                        :
                        <RowWrap style={{ marginBottom: 30, marginTop: 25, justifyContent: 'space-between' }}>

                            <Button
                                text='Connect'
                                onClick={() => setShowQR(true)}
                                color='primary'
                                height={42}
                                width={120}
                            />

                            <Button
                                text='Send Tip'
                                color='link'
                                height={42}
                                width={120}
                                onClick={() => setShowSupport(true)} />
                        </RowWrap>
                    }
                </Head>

                {renderWidgets('about')}

            </div>

            <div style={{
                width: canEdit ? 'calc(100% - 365px)' : 'calc(100% - 628px)',
                minWidth: 250
            }}>
                <Tabs style={{
                    background: '#fff', padding: '0 20px', boxShadow: '0px 1px 6px rgba(0, 0, 0, 0.07)'

                }}>
                    {tabs && Object.keys(tabs).map((name, i) => {
                        if (name === 'about') return <div key={i} />
                        const t = tabs[name]
                        const label = t.label
                        const selected = name === newSelectedWidget
                        let count = (extras && extras[name] && extras[name].length > 0) ? extras[name].length : null
                        return <Tab key={i}
                            style={{ height: 64, alignItems: 'center' }}
                            selected={selected}
                            onClick={() => {
                                switchWidgets(name)
                            }}>
                            {label}
                            {count && <Counter>
                                {count}
                            </Counter>}
                        </Tab>
                    })}

                </Tabs>

                <Modal
                    visible={showFocusView}
                    style={{
                        top: -64,
                        height: 'calc(100% + 64px)'
                    }}
                    envStyle={{
                        marginTop: (isMobile || canEdit) ? 64 : 123, borderRadius: 0, background: '#fff',
                        height: (isMobile || canEdit) ? 'calc(100% - 64px)' : '100%', width: '60%',
                        minWidth: 500, maxWidth: 602, //minHeight: 300,
                        ...focusedDesktopModalStyles
                    }}
                    nextArrow={nextIndex}
                    prevArrow={prevIndex}
                    overlayClick={() => {
                        setShowFocusView(false)
                        setFocusIndex(-1)
                        if (selectedWidget === 'about') switchWidgets('post')
                    }}
                    bigClose={() => {
                        setShowFocusView(false)
                        setFocusIndex(-1)
                        if (selectedWidget === 'about') switchWidgets('post')
                    }}
                >
                    <FocusedView
                        person={person}
                        canEdit={canEdit}
                        selectedIndex={focusIndex}
                        config={tabs[selectedWidget] && tabs[selectedWidget]}
                        onSuccess={() => {
                            console.log('success')
                            setFocusIndex(-1)
                            if (selectedWidget === 'about') switchWidgets('post')
                        }}
                        goBack={() => {
                            setShowFocusView(false)
                            setFocusIndex(-1)
                            if (selectedWidget === 'about') switchWidgets('post')
                        }}
                    />
                </Modal>

                <div style={{
                    padding: 20, height: 'calc(100% - 63px)', background: '#F2F3F5',
                    overflowY: 'auto',
                    position: 'relative',
                }}>
                    {renderEditButton({ marginBottom: 15 })}
                    {/* <div style={{ height: 15 }} /> */}
                    <Sleeve style={{
                        display: 'flex',
                        alignItems: 'flex-start',
                        justifyContent: (fullSelectedWidget && fullSelectedWidget.length > 0) ? 'flex-start' : 'center',
                        flexWrap: 'wrap',
                        height: !hasWidgets() ? 'inherit' : '',
                        paddingTop: !hasWidgets() ? 30 : 0
                    }}>
                        {renderWidgets('')}
                    </Sleeve>
                    <div style={{ height: 60 }} />
                </div>

            </div>

            <ConnectCard
                dismiss={() => setShowQR(false)}
                modalStyle={{ top: -64, height: 'calc(100% + 64px)' }}
                person={person} visible={showQR} />
        </div >
    }


    return (
        <Content>
            {isMobile ? renderMobileView() : renderDesktopView()}

            <Modal
                visible={showSupport}
                close={() => setShowSupport(false)}
                style={{
                    top: -64,
                    height: 'calc(100% + 64px)'
                }}
                envStyle={{
                    marginTop: (isMobile || canEdit) ? 64 : 123, borderRadius: 0
                }}
            >
                <div dangerouslySetInnerHTML={{
                    __html: `<sphinx-widget
                                pubkey=${owner_pubkey}
                                amount="500"
                                title="Support Me"
                                subtitle="Because I'm awesome"
                                buttonlabel="Donate"
                                defaultinterval="weekly"
                                imgurl="${img || 'https://i.scdn.co/image/28747994a80c78bc2824c2561d101db405926a37'}"
                            ></sphinx-widget>` }} />
            </Modal>
        </Content >

    );

}
interface ContentProps {
    selected: boolean;
}

const PeopleList = styled.div`
            display:flex;
            flex-direction:column;
            background:#ffffff;
            width: 265px;
            `;

const DBack = styled.div`
            height:64px;
            display:flex;
            align-items:center;
            background: #FFFFFF;
            box-shadow: 0px 1px 6px rgba(0, 0, 0, 0.07);
            z-index:0;
            `

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

const Counter = styled.div`
font-family: Roboto;
font-style: normal;
font-weight: normal;
font-size: 11px;
line-height: 19px;
margin-bottom:-3px;
/* or 173% */
margin-left:8px;

display: flex;
align-items: center;

/* Placeholder Text */

color: #B0B7BC;
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
            align-items:center;
            // justify-content:center;
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
                font-size: 15px;
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
                    line-height: 28px;
                    /* or 73% */

                    text-align: center;

                    /* Text 2 */

                    color: #3C3F41;
                    `;


const Sleeve = styled.div`

                    `;

const RowWrap = styled.div`
                    display:flex;
                    justify-content:center;

                    width:100%`;


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
                        display:flex;
                        align-items:flex-end;
                        justify-content:flex-end;
                        `;
