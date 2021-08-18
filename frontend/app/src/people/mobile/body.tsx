import React, { useEffect, useState } from 'react'
import styled from 'styled-components'
import { useObserver } from 'mobx-react-lite'
import { useStores } from '../../store'
import {
    EuiFormFieldset,
    EuiLoadingSpinner,
    EuiButtonIcon,
    EuiButton
} from '@elastic/eui';
import Person from '../person'
import PersonViewSlim from '../personViewSlim'
import EditMe from '../editMe'
import { useFuse, useScroll } from '../../hooks'
import MaterialIcon from '@material/react-material-icon';
import { colors } from '../../colors'
import FadeLeft from '../../animated/fadeLeft';
import { useIsMobile } from '../../hooks';
import {
    Switch,
    Route,
    Link,
    useHistory,
    useLocation
} from "react-router-dom";
import { Modal, Button, Divider } from '../../sphinxUI';
// avoid hook within callback warning by renaming hooks
const getFuse = useFuse
const getScroll = useScroll

export default function BodyComponent() {
    const { main, ui } = useStores()
    const [loading, setLoading] = useState(false)
    const [selectedPerson, setSelectedPerson] = useState(0)
    const [selectingPerson, setSelectingPerson] = useState(0)
    const [showProfile, setShowProfile] = useState(false)
    const c = colors['light']
    const isMobile = useIsMobile()
    const history = useHistory()
    const location = useLocation()

    console.log('history', history)
    console.log('location', location)

    function selectPerson(id: number, unique_name: string) {
        console.log('selectPerson', id, unique_name)
        setSelectedPerson(id)
        setSelectingPerson(id)
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
            if (initial && initial.unique_name === un) setSelectedPerson(initial.id)
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
        people = [...people, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}]

        return <Body>
            <div>
                {loading && <EuiLoadingSpinner size="xl" />}
                {!loading && <EuiFormFieldset style={{ width: '100%' }} >

                    {people.map(t => <Person {...t} key={t.id}
                        selected={selectedPerson === t.id}
                        select={selectPerson}
                    />)}

                </EuiFormFieldset>}
                <AddWrap>
                    {!loading && <EuiButton onClick={() => ui.setEditMe(true)} style={{ border: 'none' }}>
                        <div style={{ display: 'flex' }}>
                            <MaterialIcon
                                style={{ fontSize: 70 }}
                                icon="account_circle" aria-label="edit-me" />
                        </div>
                    </EuiButton>}
                </AddWrap>
            </div>

            <EditMe />

            <FadeLeft
                withOverlay
                drift={40}
                overlayClick={() => setSelectingPerson(0)}
                style={{ position: 'absolute', top: 0, right: 0, zIndex: 10000 }}
                isMounted={(selectingPerson && !showProfile) ? true : false}
                dismountCallback={() => setSelectedPerson(0)}
            >
                <PersonViewSlim goBack={() => setSelectingPerson(0)}
                    personId={selectedPerson}
                    loading={loading} />
            </FadeLeft>



            {/* you logged in modal  */}
            < Modal visible={false} >
                <div>
                    HIII
                </div>

            </Modal >
        </Body >
    }
    )
}


const Body = styled.div`
  flex:1;
  height:calc(100vh - 60px);
  padding-bottom:80px;
  width:100%;
  overflow:auto;
  display:flex;
`
const Column = styled.div`
  width:100%;
  display:flex;
  flex-direction:column;
  justify-content:center;
  align-items:center;
  padding: 25px;
  
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

const Name = styled.div`
font-style: normal;
font-weight: 500;
font-size: 26px;
line-height: 19px;
/* or 73% */

text-align: center;

/* Text 2 */

color: #292C33;
`;
const Description = styled.div`
font-size: 17px;
line-height: 20px;
text-align: center;
margin:20px 0;

/* Main bottom icons */

color: #5F6368;
`

interface ImageProps {
    readonly src: string;
}
const Img = styled.div<ImageProps>`
                        background-image: url("${(p) => p.src}");
                        background-position: center;
                        background-size: cover;
                        margin-bottom:20px;
                        width:90px;
                        height:90px;
                        border-radius: 50%;
                        position: relative;
                        `;