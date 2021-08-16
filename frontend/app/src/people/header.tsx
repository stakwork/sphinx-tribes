import React, { useEffect } from 'react'
import { useObserver } from 'mobx-react-lite'
import { useStores } from '../store'
import styled from 'styled-components'
import {
  EuiHeader,
  EuiHeaderSection,
  EuiFieldSearch,
} from '@elastic/eui';
import { useFuse } from '../hooks'
import { colors } from '../colors'
import { useHistory } from 'react-router-dom'

export default function Header() {
  const { main, ui } = useStores()

  const people = useFuse(main.people, ["owner_alias"])

  // function selectPerson(id: number, unique_name: string) {
  //   console.log('selectPerson', id, unique_name)
  //   setSelectedPerson(id)
  //   if (unique_name && window.history.pushState) {
  //     window.history.pushState({}, 'Sphinx Tribes', '/p/' + unique_name);
  //   }
  // }
  const c = colors['light']

  const tabs = [
    {
      text: 'Tribes',
      path: '/t/'
    },
    {
      text: 'People',
      path: '/p/'
    }
  ]

  const pathname = window.location.pathname
  console.log('pathname', pathname)

  return useObserver(() => {

    const showDropdown = ui.searchText && window.location.pathname.startsWith('/p') ? true : false

    return <EuiHeader id="header" style={{ color: '#fff' }}>
      <div className="container">
        <Row>
          <EuiHeaderSection grow={false} style={{ marginRight: 20 }}>
            <img id="logo" src="/static/people_logo.svg" alt="Logo"
            />
          </EuiHeaderSection>

          {tabs.map((t, i) => {
            console.log('pathname', pathname)
            const selected = pathname.includes(t.path)
            return <Tab
              onClick={() => {
                if (window.history.pushState) window.history.pushState({}, 'Sphinx Tribes', t.path)
                console.log('hi')
              }}
              key={i} style={{ background: selected && c.blue1 }}>
              {t.text}
            </Tab>
          })}
          {/* <EuiHeaderSection id="header-right" side="right" className="col-xs-12 col-sm-12 col-md-6 col-lg-6">
          <div style={{ position: 'relative' }}>
            <EuiFieldSearch id="search-input"
              placeholder="Search People"
              value={ui.searchText}
              onChange={e => ui.setSearchText(e.target.value)}
              // isClearable={this.state.isClearable}
              aria-label="search"
            />
            {showDropdown &&
              <SearchList>
                {people.map(t => <Person
                  onClick={() => {
                    ui.setSearchText('')
                  }}
                  key={t.owner_alias}>
                  <Img src={t.img || '/static/sphinx.png'} />
                  <Name>{t.owner_alias}</Name>
                </Person>)}
              </SearchList>
            }
          </div>
          </EuiHeaderSection> 
          */}
        </Row>
      </div>
    </EuiHeader >
  })
}

const Row = styled.div`
  display:flex;
  align-items:center;
  width:100%;
`
const Tab = styled.div`
  margin-left:10px;
  display:flex;
  justify-content:center;
  align-items:center;
  width:150px;
  padding:10px;
  height:32px;
  width:92px;
  border-radius: 5px;
  font-weight: 500;
  font-size: 13px;
  cursor:pointer;
`