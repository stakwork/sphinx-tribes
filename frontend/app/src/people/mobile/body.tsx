import React, { useEffect, useState } from 'react'
import styled from 'styled-components'
import { useObserver } from 'mobx-react-lite'
import { useStores } from '../../store'
import {
    EuiLoadingSpinner,
} from '@elastic/eui';
import Person from '../person'
import PersonViewSlim from '../personViewSlim'

import { useFuse, useScroll } from '../../hooks'
import { colors } from '../../colors'
import FadeLeft from '../../animated/fadeLeft';
import { useIsMobile } from '../../hooks';
import {
    useHistory,
    useLocation
} from "react-router-dom";
import Drawer from '../drawer/index'
// avoid hook within callback warning by renaming hooks
const getFuse = useFuse
const getScroll = useScroll

export default function BodyComponent() {
    const { main, ui } = useStores()
    const [loading, setLoading] = useState(false)

    const [showProfile, setShowProfile] = useState(false)

    const c = colors['light']
    const isMobile = useIsMobile()
    const history = useHistory()
    const location = useLocation()


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

        if (loading) {
            return <Body style={{ justifyContent: 'center', alignItems: 'center' }}>
                <EuiLoadingSpinner size="xl" />
            </Body>
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
                    isMounted={(ui.selectingPerson && !showProfile) ? true : false}
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

            <>
                <Drawer />
                <div style={{
                    width: '100%', padding: 16, paddingLeft: 0, display: 'flex', flexWrap: 'wrap',
                    justifyContent: 'space-around', alignItems: 'flex-start'
                }} >
                    {people.map(t => <Person {...t} key={t.id}
                        small={false}
                        selected={ui.selectedPerson === t.id}
                        select={selectPerson}
                    />)}
                </div>
            </>


            {/* selected view */}
            <FadeLeft
                withOverlay={isMobile}
                drift={40}
                overlayClick={() => ui.setSelectingPerson(0)}
                style={{ position: 'absolute', top: isMobile ? 0 : 65, right: 0, zIndex: 10000, width: '100%' }}
                isMounted={(ui.selectingPerson && !showProfile) ? true : false}
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
`
const AddWrap = styled.div`
  position:fixed;
  bottom:2rem;
  right:2rem;
  & button {
    height: 100px;
    width: 100px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  & svg {
    width:60px;
    height:60px;
  }
`
