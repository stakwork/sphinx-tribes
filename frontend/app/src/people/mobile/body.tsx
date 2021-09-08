import React, { useEffect, useState } from 'react'
import styled from 'styled-components'
import { useObserver } from 'mobx-react-lite'
import { useStores } from '../../store'
import {
    EuiFieldSearch,
    EuiLoadingSpinner,
} from '@elastic/eui';
import Person from '../person'
import PersonViewSlim from '../personViewSlim'
import { useFuse, useScroll } from '../../hooks'
import { colors } from '../../colors'
import FadeLeft from '../../animated/fadeLeft';
import FirstTimeScreen from './firstTimeScreen';
import { useIsMobile } from '../../hooks';
import NoneSpace from '../utils/noneSpace';
import { Divider, SearchTextInput } from '../../sphinxUI';
import WidgetSwitchViewer from '../widgetViews/widgetSwitchViewer';

// import { SearchTextInput } from '../../sphinxUI/index'
// avoid hook within callback warning by renaming hooks
const getFuse = useFuse
const getScroll = useScroll

export default function BodyComponent() {
    const { main, ui } = useStores()
    const [loading, setLoading] = useState(false)
    const [selectedWidget, setSelectedWidget] = useState('people')

    const c = colors['light']

    const tabs = [
        {
            label: 'People',
            name: 'people',

        },
        {
            label: 'Posts',
            name: 'post',

        },
        {
            label: 'Offers',
            name: 'offer',

        },
        {
            label: 'Wanted',
            name: 'wanted',

        },
    ]
    const isMobile = useIsMobile()

    function selectPerson(id: number, unique_name: string) {
        console.log('selectPerson', id, unique_name)
        ui.setSelectedPerson(id)
        ui.setSelectingPerson(id)
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




    return useObserver(() => {
        const peeps = getFuse(main.people, ["owner_alias"])
        const { handleScroll, n, loadingMore } = getScroll()
        let people = peeps.slice(0, n)

        people = (people && people.filter(f => !f.hide)) || []

        function renderPeople() {
            const p = people && people.map(t => <Person {...t} key={t.id}
                small={false}
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

        if (isMobile) {
            return <Body>
                <div style={{ width: '100%' }} >
                    {people.map(t => <Person {...t} key={t.id}
                        selected={ui.selectedPerson === t.id}
                        small={isMobile}
                        select={selectPerson}
                    />)}
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
            </Body >
        }

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
                                selectPerson(person.id, person.unique_name)
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
                    selectPerson={selectPerson} />
            </FadeLeft>
        </Body >
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