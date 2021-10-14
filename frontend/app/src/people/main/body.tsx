import React, { useEffect, useState } from 'react'
import styled from 'styled-components'
import { useObserver } from 'mobx-react-lite'
import { useStores } from '../../store'
import {
    EuiGlobalToastList,
    EuiLoadingSpinner,
} from '@elastic/eui';
import Person from '../person'
import PersonViewSlim from '../personViewSlim'
import { useFuse, useScroll, useIsMobile, useScreenWidth } from '../../hooks'
import { colors } from '../../colors'
import FadeLeft from '../../animated/fadeLeft';
import FirstTimeScreen from './firstTimeScreen';
import moment from 'moment';
import NoneSpace from '../utils/noneSpace';
import { Divider, SearchTextInput, Modal } from '../../sphinxUI';
import WidgetSwitchViewer from '../widgetViews/widgetSwitchViewer';
import MaterialIcon from '@material/react-material-icon';
import FocusedView from './focusView';

import { widgetConfigs } from '../utils/constants'

// import { SearchTextInput } from '../../sphinxUI/index'
// avoid hook within callback warning by renaming hooks

const getFuse = useFuse
const getScroll = useScroll

export default function BodyComponent() {
    const { main, ui } = useStores()
    const [loading, setLoading] = useState(false)
    const [showDropdown, setShowDropdown] = useState(false)
    const screenWidth = useScreenWidth()
    const [publicFocusPerson, setPublicFocusPerson]: any = useState(null)
    const [publicFocusIndex, setPublicFocusIndex] = useState(-1)

    const [selectedWidget, setSelectedWidget] = useState('people')

    const c = colors['light']

    const tabs = [
        {
            label: 'People',
            name: 'people',

        },
        widgetConfigs['post'],
        widgetConfigs['offer'],
        widgetConfigs['wanted'],
    ]
    const isMobile = useIsMobile()

    useEffect(() => {
        // clear public focus is selected person
        if (ui.selectedPerson) {
            setPublicFocusPerson(null)
            setPublicFocusIndex(-1)
        }
    }, [ui.selectedPerson])

    function selectPerson(id: number, unique_name: string) {
        console.log('selectPerson', id, unique_name)
        ui.setSelectedPerson(id)
        ui.setSelectingPerson(id)
        // setPublicFocusPerson(null)
        // setPublicFocusIndex(-1)
        if (unique_name && window.history.pushState) {
            window.history.pushState({}, 'Sphinx Tribes', '/p/' + unique_name);
        }
    }

    async function loadPeople() {
        setLoading(true)
        let un = ''
        if (window.location.pathname.startsWith('/p/')) {
            un = window.location.pathname.substr(3)
        }
        const ps = await main.getPeople(un)
        if (un) {
            const initial = ps[0]
            if (initial && initial.unique_name === un) ui.setSelectedPerson(initial.id)
        }
        setLoading(false)
    }

    useEffect(() => {
        loadPeople()
    }, [])

    useEffect(() => {
        if (ui.meInfo) {
            main.getTribesByOwner(ui.meInfo.owner_pubkey || '')
        }
    }, [ui.meInfo])



    function publicPanelClick(person, widget, i) {
        setPublicFocusPerson(person)
        setPublicFocusIndex(i)
    }

    return useObserver(() => {
        const peeps = getFuse(main.people, ["owner_alias"])
        const { handleScroll, n, loadingMore } = getScroll()
        let people = peeps.slice(0, n)

        people = (people && people.filter(f => !f.hide)) || []

        function renderPeople() {
            // clone, sort, reverse, return
            const peopleClone = [...people]
            const p = peopleClone && peopleClone.sort((a: any, b: any) => {
                return moment(a.updated).valueOf() - moment(b.updated).valueOf()
            }).reverse().map(t => <Person {...t} key={t.id}
                small={isMobile}
                squeeze={screenWidth < 1420}
                selected={ui.selectedPerson === t.id}
                select={selectPerson}
            />)
            return p
        }

        if (loading) {
            return <Body style={{ justifyContent: 'center', alignItems: 'center' }}>
                <EuiLoadingSpinner size="xl" />
            </Body>
        }

        const showFirstTime = ui.meInfo && ui.meInfo.id === 0

        if (showFirstTime) {
            return <FirstTimeScreen />
        }

        const widgetLabel = selectedWidget && tabs.find(f => f.name === selectedWidget)

        const toastsEl = <EuiGlobalToastList
            toasts={ui.toasts}
            dismissToast={() => ui.setToasts([])}
            toastLifeTimeMs={3000}
        />

        if (isMobile) {
            return <Body>

                {!ui.meInfo &&
                    <div style={{ marginTop: 60 }}>
                        <NoneSpace
                            buttonText={'Get Started'}
                            buttonIcon={'arrow_forward'}
                            action={() => ui.setShowSignIn(true)}
                            img={'explore.png'}
                            text={'Discover people on Sphinx'}
                            style={{ height: 320, background: '#fff' }}
                        />
                        <Divider />
                    </div>
                }

                <div style={{
                    width: '100%', display: 'flex',
                    justifyContent: 'space-between', alignItems: 'flex-start', padding: 20,
                    height: 62, marginBottom: 20
                }}>
                    <Label style={{ fontSize: 20 }}>
                        Explore
                        <Link
                            onClick={() => setShowDropdown(true)}>
                            <div>{widgetLabel && widgetLabel.label}</div>
                            <MaterialIcon icon={'expand_more'} style={{ fontSize: 18, marginLeft: 5 }} />

                            {showDropdown && <div style={{ position: 'absolute', top: 0, left: 0, zIndex: 10, background: '#fff' }}>
                                {tabs && tabs.map((t, i) => {
                                    const label = t.label
                                    const selected = selectedWidget === t.name

                                    return <Tab key={i}
                                        style={{ borderRadius: 0, margin: 0 }}
                                        selected={selected}
                                        onClick={(e) => {
                                            e.stopPropagation()
                                            setShowDropdown(false)
                                            setSelectedWidget(t.name)
                                        }}>
                                        {label}
                                    </Tab>
                                })}
                            </div>}
                        </Link>
                    </Label>



                    <SearchTextInput
                        small
                        name='search'
                        type='search'
                        placeholder='Search'
                        value={ui.searchText}
                        style={{ width: 164, height: 40, border: '1px solid #DDE1E5', background: '#fff' }}
                        onChange={e => {
                            console.log('handleChange', e)
                            ui.setSearchText(e)
                        }}

                    />

                </div>

                <div style={{ width: '100%' }} >
                    {selectedWidget === 'people' ?
                        renderPeople()
                        : <WidgetSwitchViewer
                            onPanelClick={(person, widget, i) => {
                                publicPanelClick(person, widget, i)
                                // selectPerson(person.id, person.unique_name)
                            }}
                            selectedWidget={selectedWidget} />
                    }
                </div>

                <FadeLeft
                    withOverlay
                    drift={40}
                    overlayClick={() => ui.setSelectingPerson(0)}
                    style={{ position: 'absolute', top: 0, right: 0, zIndex: 10000, width: '100%' }}
                    isMounted={ui.selectingPerson ? true : false}
                    dismountCallback={() => ui.setSelectedPerson(0)}
                >
                    <PersonViewSlim goBack={() => ui.setSelectingPerson(0)}
                        personId={ui.selectedPerson}
                        selectPerson={selectPerson}
                        loading={loading} />
                </FadeLeft>

                <Modal
                    visible={publicFocusPerson ? true : false}
                    fill={true}
                >
                    <FocusedView
                        person={publicFocusPerson}
                        canEdit={false}
                        selectedIndex={publicFocusIndex}
                        config={widgetConfigs[selectedWidget] && widgetConfigs[selectedWidget]}
                        onSuccess={() => {
                            console.log('success')
                            setPublicFocusPerson(null)
                            setPublicFocusIndex(-1)
                        }}
                        goBack={() => {
                            setPublicFocusPerson(null)
                            setPublicFocusIndex(-1)
                        }}
                    />
                </Modal>

                {toastsEl}
            </Body >
        }

        const focusedDesktopModalStyles = (selectedWidget && widgetConfigs[selectedWidget]) ? {
            ...widgetConfigs[selectedWidget].modalStyle
        } : {}

        // desktop mode
        return <Body style={{
            background: '#f0f1f3',
            height: 'calc(100% - 65px)'
        }}>

            {!ui.meInfo &&
                <div>
                    <NoneSpace
                        banner
                        buttonText={'Get Started'}
                        buttonIcon={'arrow_forward'}
                        action={() => ui.setShowSignIn(true)}
                        img={'explore.png'}
                        text={'Discover people on Sphinx'}
                        style={{ height: 320 }}
                    />
                    <Divider />
                </div>
            }

            <div style={{
                width: '100%', display: 'flex',
                justifyContent: 'space-between', alignItems: 'flex-start', padding: 20,
                height: 62
            }}>
                <Label>
                    Explore
                </Label>

                <Tabs>
                    {tabs && tabs.map((t, i) => {
                        const label = t.label
                        const selected = selectedWidget === t.name

                        return <Tab key={i}
                            selected={selected}
                            onClick={() => {
                                setSelectedWidget(t.name)
                            }}>
                            {label}
                        </Tab>
                    })}

                </Tabs>

                <SearchTextInput
                    name='search'
                    type='search'
                    placeholder='Search'
                    value={ui.searchText}
                    style={{ width: 204, height: 40, background: '#DDE1E5' }}
                    onChange={e => {
                        console.log('handleChange', e)
                        ui.setSearchText(e)
                    }}

                />
            </div>
            <>
                <div style={{
                    width: '100%', display: 'flex', flexWrap: 'wrap', height: '100%',
                    justifyContent: 'flex-start', alignItems: 'flex-start', padding: 20
                }}>
                    {selectedWidget === 'people' ?
                        renderPeople()
                        : <WidgetSwitchViewer
                            onPanelClick={(person, widget, i) => {
                                publicPanelClick(person, widget, i)
                                // selectPerson(person.id, person.unique_name)
                            }}
                            selectedWidget={selectedWidget} />
                    }
                </div>
                <div style={{ height: 100 }} />
            </>


            {/* selected view */}
            <FadeLeft
                withOverlay={isMobile}
                drift={40}
                overlayClick={() => ui.setSelectingPerson(0)}
                style={{ position: 'absolute', top: isMobile ? 0 : 65, right: 0, zIndex: 10000, width: '100%' }}
                isMounted={ui.selectingPerson ? true : false}
                dismountCallback={() => ui.setSelectedPerson(0)}
            >
                <PersonViewSlim goBack={() => ui.setSelectingPerson(0)}
                    personId={ui.selectedPerson}
                    loading={loading}
                    peopleView={true}
                    selectPerson={selectPerson} />
            </FadeLeft>


            <Modal
                visible={publicFocusPerson ? true : false}
                envStyle={{
                    borderRadius: 0, background: '#fff',
                    height: '100%', width: '60%',
                    minWidth: 500, maxWidth: 602,
                    ...focusedDesktopModalStyles
                }}
                bigClose={() => {
                    setPublicFocusPerson(null)
                    setPublicFocusIndex(-1)
                }}
            >
                <FocusedView
                    person={publicFocusPerson}
                    canEdit={false}
                    selectedIndex={publicFocusIndex}
                    config={widgetConfigs[selectedWidget] && widgetConfigs[selectedWidget]}
                    onSuccess={() => {
                        console.log('success')
                        setPublicFocusPerson(null)
                        setPublicFocusIndex(-1)
                    }}
                    goBack={() => {
                        setPublicFocusPerson(null)
                        setPublicFocusIndex(-1)
                    }}
                />
            </Modal>

            {toastsEl}
        </Body>
    }
    )

}

const Body = styled.div`
            flex:1;
            height:calc(100% - 105px);
            padding-bottom:80px;
            width:100%;
            overflow:auto;
            display:flex;
            flex-direction:column;
            `
const Label = styled.div`
            font-family: Roboto;
            font-style: normal;
            font-weight: bold;
            font-size: 26px;
            line-height: 40px;
            /* or 154% */
            width:204px;
            
            display: flex;
            align-items: center;
            
            /* Text 2 */
            
            color: #3C3F41;`

const Tabs = styled.div`
                        display:flex;
                        `;

interface TagProps {
    selected: boolean;
}
const Tab = styled.div<TagProps>`
            display:flex;
            padding:10px 25px;
            margin-right:35px;
            color:${p => p.selected ? '#5078F2' : '#5F6368'};
        // border-bottom: ${p => p.selected && '4px solid #618AFF'};
            cursor:pointer;
            font-weight: 500;
            font-size: 15px;
            line-height: 19px;
            background:${p => p.selected ? '#DCEDFE' : '#3C3F4100'};
            border-radius:25px;
            `;

const Link = styled.div`
            display:flex;
            justify-content:center;
            align-items:center;
            margin-left:6px;
            color:#618AFF;
            cursor:pointer;
            position:relative;
            `;