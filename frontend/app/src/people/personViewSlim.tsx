import React, { useState, useEffect, useRef } from "react";
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

import { Button, IconButton, Modal, SearchTextInput } from "../sphinxUI";
import MaterialIcon from "@material/react-material-icon";
import FocusedView from './main/focusView'
import { meSchema } from "../form/schema";
import { useIsMobile, usePageScroll } from "../hooks";
import Person from "./person";
import NoneSpace from "./utils/noneSpace";
import ConnectCard from "./utils/connectCard";
import { widgetConfigs } from "./utils/constants";
import { extractGithubIssue } from "../helpers";
import { useHistory, useLocation } from "react-router";
import { EuiLoadingSpinner } from '@elastic/eui';
import { queryLimit } from '../store/main'
import NoResults from "./utils/noResults";
import PageLoadSpinner from "./utils/pageLoadSpinner";

import Badges from './utils/badges'


const host = getHost();
function makeQR(pubkey: string) {
    return `sphinx.chat://?action=person&host=${host}&pubkey=${pubkey}`;
}

let deeplinkTimeout

export default function PersonView(props: any) {

    const {
        personId,
        loading,
        selectPerson,
        goBack,
    } = props

    const { main, ui } = useStores()
    const { meInfo, peoplePageNumber } = ui || {}

    const [loadingPerson, setLoadingPerson]: any = useState(false)
    const [loadedPerson, setLoadedPerson]: any = useState(null)

    const history = useHistory()
    const location = useLocation()
    const pathname = history?.location?.pathname

    // FOR PEOPLE VIEW
    let person: any = (main.people && main.people.length && main.people.find(f => f.id === personId))

    // migrating to loading person on person view load
    if (loadedPerson) {
        person = loadedPerson
    }

    // if i select myself, fill person with meInfo
    if (personId === ui.meInfo?.id) {
        person = {
            ...ui.meInfo
        }
    }

    let people: any = (main.people && main.people.filter(f => !f.hide)) || []

    const {
        id,
        img,
        tags,
        owner_alias,
        unique_name,
        price_to_meet,
        extras,
        owner_pubkey
    } = person || {}

    let { description } = person || {}

    // backend is adding 'description' to empty descriptions, short term fix
    if (description === 'description') description = ''

    const canEdit = id === meInfo?.id
    const isMobile = useIsMobile()

    const initialWidget = (!isMobile || canEdit) ? 'badges' : 'about'

    const [selectedWidget, setSelectedWidget] = useState(initialWidget);
    const [newSelectedWidget, setNewSelectedWidget] = useState(initialWidget);
    const [focusIndex, setFocusIndex] = useState(-1);
    const [showSupport, setShowSupport] = useState(false);

    const [showQR, setShowQR] = useState(false);
    const [showFocusView, setShowFocusView] = useState(false);
    const qrString = makeQR(owner_pubkey || '');

    async function loadMorePeople(direction) {
        let newPage = peoplePageNumber + direction
        if (newPage < 1) newPage = 1

        await main.getPeople({ page: newPage })
    }

    // if no people, load people on mount
    useEffect(() => {
        if (!people.length) main.getPeople({ page: 1, resetPage: true });
    }, [])

    // deeplink load person
    useEffect(() => {
        if (loadedPerson) {
            doDeeplink()
        } else {
            deeplinkTimeout = setTimeout(() => {
                doDeeplink()
            }, 200)
        }

        return function cleanup() {
            clearTimeout(deeplinkTimeout)
        }
    }, [personId])

    async function doDeeplink() {
        if (pathname) {
            let splitPathname = pathname?.split('/')
            let personPubkey: string = splitPathname[2]
            if (personPubkey) {
                setLoadingPerson(true)
                let p = await main.getPersonByPubkey(personPubkey)
                setLoadedPerson(p)
                setLoadingPerson(false)

                const search = location?.search

                // deeplink for widgets
                let widgetName: any = new URLSearchParams(search).get("widget")
                let widgetTimestamp: any = new URLSearchParams(search).get("timestamp")

                if (widgetName) {
                    setNewSelectedWidget(widgetName)
                    setSelectedWidget(widgetName)
                    if (widgetTimestamp) {
                        const thisExtra = p?.extras && p?.extras[widgetName]
                        const thisItemIndex = thisExtra && thisExtra.length && thisExtra.findIndex(f => f.created === parseInt(widgetTimestamp))
                        if (thisItemIndex > -1) {
                            // select it!
                            setFocusIndex(thisItemIndex)
                            setShowFocusView(true)
                        }
                    }
                }
            }
        }
    }

    function updatePath(name) {
        history.push(location.pathname + `?widget=${name}`)
    }

    function updatePathIndex(timestamp) {
        history.push(location.pathname + `?widget=${selectedWidget}&timestamp=${timestamp}`)
    }

    function switchWidgets(name) {
        setNewSelectedWidget(name)
        setSelectedWidget(name)
        updatePath(name)
        setShowFocusView(false)
        setFocusIndex(-1)
    }

    function selectPersonWithinFocusView(id, unique_name, pubkey) {
        setShowFocusView(false)
        setFocusIndex(-1)
        selectPerson(id, unique_name, pubkey)
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
        main.getPeople({ resetPage: true })
        goBack()
    }

    const { loadingTop, loadingBottom, handleScroll } = usePageScroll(() => loadMorePeople(1), () => loadMorePeople(-1))

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
        if (selectedWidget === 'badges') {
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

        if (selectedWidget === 'badges') {
            return <Badges person={person} />
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

                if (!canEdit &&
                    'show' in s &&
                    s.show === false) {
                    // skip hidden items
                    return
                }

                elementArray.push(<Panel key={i}
                    onClick={() => {
                        setShowFocusView(true)
                        setFocusIndex(i)
                        if (s.created) updatePathIndex(s.created)
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

            const noneKey = canEdit ? 'me' : 'otherUser'
            const panels: any = elementArray.length ? elementArray : (<NoneSpace
                action={() => setShowFocusView(true)}
                {...tabs[selectedWidget]?.noneSpace[noneKey]}
            />)

            return <>
                <PageLoadSpinner show={loadingPerson} />
                {panels}
            </>
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
            let g = person.extras[tabs[selectedWidget]?.name]
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
            let g = person?.extras[tabs[selectedWidget]?.name]
            let previndex = focusIndex - 1
            if (g[previndex]) setFocusIndex(previndex)
            else setFocusIndex(g.length - 1)
        }
    }

    function renderEditButton(style: any) {
        if (!canEdit || !selectedWidget) return <div />

        if (selectedWidget === 'badges') return <div />

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

    const defaultPic = '/static/person_placeholder.png'
    const mediumPic = img && img + '?medium=true'

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
                    <Img src={mediumPic || defaultPic} />
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
                        let count = (extras && extras[name] && extras[name].length > 0) ? extras[name].length : null

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

    const loaderTop = <PageLoadSpinner show={loadingTop} style={{ paddingTop: 10 }} />
    const loaderBottom = <PageLoadSpinner noAnimate show={loadingBottom} style={{ position: 'absolute', bottom: 0, left: 0 }} />

    function renderDesktopView() {
        const focusedDesktopModalStyles = newSelectedWidget ? {
            ...tabs[newSelectedWidget]?.modalStyle
        } : {}

        return <div style={{
            display: 'flex',
            width: '100%', height: '100%'
        }}>

            {!canEdit &&
                <PeopleList >

                    <DBack >
                        <Button
                            color='clear'
                            leadingIcon='arrow_back'
                            text='Back'
                            onClick={goBack}
                        />

                        <SearchTextInput
                            small
                            name='search'
                            type='search'
                            placeholder='Search'
                            value={ui.searchText}
                            style={{ width: 120, height: 40, border: '1px solid #DDE1E5', background: '#fff' }}
                            onChange={e => {
                                console.log('handleChange', e)
                                ui.setSearchText(e)
                            }}
                        />
                    </DBack>



                    <div style={{ width: '100%', overflowY: 'auto', height: '100%' }} onScroll={handleScroll}>
                        {loaderTop}
                        {people?.length ? people.map(t => <Person {...t} key={t.id}
                            selected={personId === t.id}
                            hideActions={true}
                            small={true}
                            select={selectPersonWithinFocusView}
                        />) : <NoResults />}

                        {/* make sure you can always scroll ever with too few people */}
                        {people?.length < queryLimit &&
                            <div style={{ height: 400 }} />
                        }
                    </div>

                    {loaderBottom}


                </PeopleList>
            }

            <div style={{
                width: 364,
                minWidth: 364, overflowY: 'auto',
                // position: 'relative',
                background: '#ffffff',
                color: '#000000',
                padding: 40,
                zIndex: 5,
                // height: '100%',
                marginTop: canEdit ? 64 : 0,
                height: canEdit ? 'calc(100% - 64px)' : '100%',
                borderLeft: '1px solid #ebedef',
                borderRight: '1px solid #ebedef',
                boxShadow: '1px 2px 6px -2px rgba(0, 0, 0, 0.07)'
            }}>

                {canEdit && <div style={{
                    position: 'absolute',
                    top: 0, left: 0,
                    display: 'flex',
                    background: '#ffffff',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    width: 364,
                    minWidth: 364,
                    boxShadow: '0px 2px 6px rgba(0, 0, 0, 0.07)',
                    borderBottom: 'solid 1px #ebedef',
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
                    <div style={{ height: 35 }} />

                    <Img src={mediumPic || defaultPic} >
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
                            style={{ fontSize: 13, background: '#f2f3f5' }}
                            leadingIcon={'edit'}
                            iconSize={15}
                        />
                        <Button
                            text='Sign out'
                            onClick={logout}
                            height={42}
                            style={{ fontSize: 13, color: '#3c3f41', }}
                            iconStyle={{ color: '#8e969c' }}
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
                minWidth: 250, zIndex: canEdit ? 6 : 4,
            }}>
                <Tabs style={{
                    background: '#fff', padding: '0 20px',
                    borderBottom: 'solid 1px #ebedef',
                    boxShadow: canEdit ? '0px 2px 0px rgba(0, 0, 0, 0.07)' : '0px 2px 6px rgba(0, 0, 0, 0.07)'

                }}>
                    {tabs && Object.keys(tabs).map((name, i) => {
                        if (name === 'about') return <div key={i} />
                        const t = tabs[name]
                        const label = t.label
                        const selected = name === newSelectedWidget
                        const hasExtras = (extras && extras[name] && (extras[name].length > 0))
                        let count: any = hasExtras ? extras[name].length : null
                        // count only open ones
                        if (hasExtras && name === 'wanted') {
                            count = 0
                            extras[name].forEach((w => {
                                const { repo, issue } = w
                                const { status } = extractGithubIssue(person, repo, issue)
                                if (status === '' || status === 'open') count++
                            }))
                        }
                        return <Tab key={i}
                            style={{ height: 64, alignItems: 'center' }}
                            selected={selected}
                            onClick={() => {
                                switchWidgets(name)
                            }}>
                            {label}
                            {(count > 0) && <Counter>
                                {count}
                            </Counter>}
                        </Tab>
                    })}

                </Tabs>


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
                modalStyle={{ top: -63, height: 'calc(100% + 64px)' }}
                person={person} visible={showQR} />


            <Modal
                visible={showFocusView}
                style={{
                    top: -64,
                    height: 'calc(100% + 64px)'
                }}
                envStyle={{
                    marginTop: (isMobile || canEdit) ? 64 : 123, borderRadius: 0, background: '#fff',
                    height: (isMobile || canEdit) ? 'calc(100% - 64px)' : '100%', width: '60%',
                    minWidth: 500, maxWidth: 602, zIndex: 20,//minHeight: 300, 
                    ...focusedDesktopModalStyles
                }}
                nextArrow={nextIndex}
                prevArrow={prevIndex}
                overlayClick={() => {
                    setShowFocusView(false)
                    setFocusIndex(-1)
                    if (selectedWidget === 'about') switchWidgets('badges')
                }}
                bigClose={() => {
                    setShowFocusView(false)
                    setFocusIndex(-1)
                    if (selectedWidget === 'about') switchWidgets('badges')
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
                        if (selectedWidget === 'about') switchWidgets('badges')
                    }}
                    goBack={() => {
                        setShowFocusView(false)
                        setFocusIndex(-1)
                        if (selectedWidget === 'about') switchWidgets('badges')
                    }}
                />
            </Modal>
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
                                imgurl="${mediumPic || 'https://i.scdn.co/image/28747994a80c78bc2824c2561d101db405926a37'}"
                            ></sphinx-widget>` }} />
            </Modal>
        </Content >

    );

}
interface ContentProps {
    selected: boolean;
}

const PeopleList = styled.div`
            position:relative;
            display:flex;
            flex-direction:column;
            background:#ffffff;
            width: 265px;
            `;

const DBack = styled.div`
            min-height:64px;
            height:64px;
            display:flex;
            padding-right:10px;
            align-items:center;
            justify-content:space-between;
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

const Name = styled.div`
                    font-style: normal;
                    font-weight: 600;
                    font-size: 24px;
                    line-height: 28px;
                    /* or 73% */

                    text-align: center;

                    /* Text 2 */

                    color: #3C3F41;
                    `;


const Sleeve = styled.div`

                    `;

const Loader = styled.div`
            position:absolute;
            width:100%;
            display:flex;
            justify-content:center;
            padding:10px;
            left:0px;
            z-index:20;
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
