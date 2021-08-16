import React, { useEffect, useState } from 'react'
import styled from 'styled-components'
import { useObserver } from 'mobx-react-lite'
import { useStores } from '../store'
import {
  EuiFormFieldset,
  EuiLoadingSpinner,
  EuiButtonIcon,
  EuiButton
} from '@elastic/eui';
import Person from './person'
import Drawer from './drawer/index'
import PersonView from './personView'
import PersonViewSlim from './personViewSlim'
import EditMe from './editMe'
import { useFuse, useScroll } from '../hooks'
import MaterialIcon from '@material/react-material-icon';
import { colors } from '../colors'
import FadeLeft from '../animated/fadeLeft';

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

    if (selectedPerson && showProfile) {
      return <Body>
        <Column className="main-wrap">
          <PersonView goBack={() => {
            setSelectedPerson(0)
            window.history.pushState({}, 'Sphinx Tribes', '/');
          }}
            personId={selectedPerson}
            loading={loading} />
        </Column>
      </Body>
    }

    return <Body>
      <Drawer />

      <Column className="main-wrap">
        {loading && <EuiLoadingSpinner size="xl" />}
        {!loading && <EuiFormFieldset style={{ width: '100%' }} >
          <Row>
            {people.map(t => <Person {...t} key={t.id}
              selected={selectedPerson === t.id}
              select={selectPerson}
            />)}
          </Row>
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
      </Column>
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
    </Body>
  }
  )
}


const Body = styled.div`
  flex:1;
  height:calc(100vh - 90px);
  padding-bottom:80px;
  width:100%;
  overflow:auto;
  display:flex;
`
const Column = styled.div`
  display:flex;
  // flex-direction:column;
  // align-items:center;
  max-width:900px;
  flex-wrap:wrap;
  width:100%;
`

const Row = styled.div`
  display:flex;
  flex-wrap:wrap;
  width:100%;
  
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