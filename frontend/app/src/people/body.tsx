import React, {useEffect, useState} from 'react'
import styled from 'styled-components'
import { useObserver } from 'mobx-react-lite'
import { useStores } from '../store'
import {
  EuiFormFieldset,
  EuiLoadingSpinner,
  EuiButtonIcon
} from '@elastic/eui';
import Person from './person'
import EditMe from './editMe'
import {useFuse, useScroll} from '../hooks'

export default function BodyComponent() {
  const { main, ui } = useStores()
  const [loading,setLoading] = useState(false)
  const [selectedPerson, setSelectedPerson] = useState(0)

  function selectPerson(id:number, unique_name:string){
    setSelectedPerson(id)
  }

  async function loadPeople(){
    setLoading(true)
    const ps = await main.getPeople('')
    setLoading(false)
  }
  useEffect(()=>{
    loadPeople()
  }, [])

  return useObserver(() => {
    const peeps = useFuse(main.people, ["owner_alias"])
    const {handleScroll, n, loadingMore} = useScroll()
    const people = peeps.slice(0,n)

    return <Body id="main">
      <Column className="main-wrap">
        {loading && <EuiLoadingSpinner size="xl" />}
        {!loading && <EuiFormFieldset style={{width:'100%'}} className="container">
          <div className="row">
            {people.map(t=> <Person {...t} key={t.id}
              selected={selectedPerson===t.id}
              select={selectPerson}
            />)}
          </div>
        </EuiFormFieldset>}
        <AddWrap>
          {!loading && <EuiButtonIcon 
            onClick={()=> ui.setEditMe(true)}
            iconType="plusInCircleFilled"
            iconSize="l"
            size="m"
            aria-label="edit-me"
          />}
        </AddWrap>
      </Column>

      <EditMe />

    </Body>
  }
)}

const Body = styled.div`
  flex:1;
  height:calc(100vh - 90px);
  padding-bottom:80px;
  width:100%;
  overflow:auto;
  background:#272c4b;
  display:flex;
  flex-direction:column;
  align-items:center;
`
const Column = styled.div`
  display:flex;
  flex-direction:column;
  align-items:center;
  max-width:900px;
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