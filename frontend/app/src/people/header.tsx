import React from 'react'
import { useObserver } from 'mobx-react-lite'
import { useStores } from '../store'
import styled from 'styled-components'
import {
  EuiHeader,
  EuiHeaderSection,
  EuiFieldSearch,
} from '@elastic/eui';
import { useFuse } from '../hooks'

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

  return useObserver(() => {

    const showDropdown = ui.searchText && window.location.pathname.startsWith('/p') ? true : false

    return <EuiHeader id="header" >
      <div className="container">
        <div className="row">
          <EuiHeaderSection grow={false} className="col-xs-12 col-sm-12 col-md-6 col-lg-6">
            <img id="logo" src="/static/people_logo.svg" alt="Logo"
            // style={{ cursor: 'pointer' }}
            // onClick={() => {
            //   window.history.pushState({}, 'Sphinx Tribes', '/');
            //   console.log('click!')
            // }}
            />
            {/*<Title>Tribes</Title>*/}
          </EuiHeaderSection>

          <EuiHeaderSection id="header-right" side="right" className="col-xs-12 col-sm-12 col-md-6 col-lg-6">
            {/* <EuiHeaderSectionItem> */}
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
        </div>
      </div>
    </EuiHeader>
  })
}

const Person = styled.div`
  display:flex;
  align-items:center;
  width:100%;
  background:#ffffff22;
  margin-bottom:3px;
  padding:5px;
  cursor:pointer;
  &:hover{
    background:#ffffff33;
  }
`


const SearchList = styled.div`
  position:absolute;
  top:40px;
  left:0px;
  display:flex;
  flex-direction:column;
  align-items:center;
  max-width:260px;
  width:100%;
`

interface ImageProps {
  readonly src: string;
}
const Img = styled.div<ImageProps>`
  background-image: url("${(p) => p.src}");
  background-position: center;
  background-size: cover;
  height: 30px;
  width: 30px;
  border-radius: 50%;
  position: relative;
  margin-right:5px;
`;

const Name = styled.div`
  font-weight:bold;
  color:#fff;
`