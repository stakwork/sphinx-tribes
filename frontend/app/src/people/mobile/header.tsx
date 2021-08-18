import React, { useEffect } from 'react'
import { useObserver } from 'mobx-react-lite'
import { useStores } from '../../store'
import styled from 'styled-components'
import {
    EuiHeader,
    EuiHeaderSection,
    EuiFieldSearch,
} from '@elastic/eui';
import { useFuse } from '../../hooks'
import { colors } from '../../colors'
import { useHistory, useLocation } from 'react-router-dom'
import { Button } from '../../sphinxUI';

export default function Header() {
    const { main, ui } = useStores()

    const people = useFuse(main.people, ["owner_alias"])
    const location = useLocation()

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

    const pathname = location && location.pathname
    console.log('pathname', pathname)

    return useObserver(() => {
        return <EuiHeader id="header" style={{ color: '#fff' }}>
            <div className="container">
                <Row style={{ justifyContent: 'space-between' }}>
                    <EuiHeaderSection grow={false}>
                        <Img src="/static/people_logo.svg" />
                    </EuiHeaderSection>

                    <Corner>
                        <Button
                            icon={'account_circle'}
                            text={'Sign in'}
                            color='primary'
                        />
                    </Corner>

                    {/* {tabs.map((t, i) => {
                        const selected = pathname.includes(t.path)
                        return <Tab
                            onClick={() => {
                                if (window.history.pushState) window.history.pushState({}, 'Sphinx Tribes', t.path)
                                console.log('hi')
                            }}
                            key={i} style={{ background: selected && c.blue1 }}>
                            {t.text}
                        </Tab>
                    })} */}


                </Row>

                <EuiHeaderSection id="header-right" side="right" style={{
                    background: '#000000',
                    boxShadow: 'inset 0px 1px 2px rgba(0, 0, 0, 0.15)',
                    borderRadius: 50, overflow: 'hidden'
                }}>
                    <EuiFieldSearch id="search-input"
                        placeholder="Search for People"
                        value={ui.searchText}
                        onChange={e => ui.setSearchText(e.target.value)}
                        style={{ width: '100%', height: '100%' }}
                        aria-label="search"

                    />
                </EuiHeaderSection>
            </div>
        </EuiHeader >
    })
}

const Row = styled.div`
  display:flex;
  align-items:center;
  width:100%;
`
const Corner = styled.div`
  display:flex;
  align-items:center;
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

interface ImageProps {
    readonly src: string;
}
const Img = styled.div<ImageProps>`
    background-image: url("${(p) => p.src}");
    background-position: center;
    background-size: cover;
    height:37px;
    width:232px;
    
    position: relative;
  `;