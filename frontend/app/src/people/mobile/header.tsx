import React, { useState, useEffect } from 'react'
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
import { Modal, Button, Divider } from '../../sphinxUI';
import FadeLeft from '../../animated/fadeLeft';
// import EditInfo from '../edit/editInfo'
import SignIn from '../auth/signIn';

import PersonViewSlim from '../personViewSlim';
import { MeInfo } from '../../store/ui';
import api from '../../api';
import ConfirmMe from '../confirmMe';

let heartbeat

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

    const [showSignIn, setShowSignIn] = useState(false)
    const [showWelcome, setShowWelcome] = useState(false)
    const [showInitEditSelf, setShowInitEditSelf] = useState(false)
    const [showEditSelf, setShowEditSelf] = useState(false)

    async function testChallenge(chal: string) {
        try {
            const me: any = await api.get(`poll/${chal}`)
            if (me && me.pubkey) {
                ui.setMeInfo(me)
                ui.setChallenge(chal)
                setShowSignIn(false)
                setShowWelcome(true)
            }
        } catch (e) {
            console.log(e)
        }
    }

    function forceLogout() {
        ui.setMeInfo(null)
        ui.setChallenge('')
        alert('Session timeout')
    }

    async function ping() {
        let chal = ui.challenge?.challenge
        if (!chal) {
            console.log('no challenge')
            return
        }

        try {
            const me: any = await api.get(`poll/${chal}`)
            if (me && me.pubkey) {
            } else {
                // sign out
                forceLogout()
            }
        } catch (e) {
            console.log(e)
            forceLogout()
        }
    }

    useEffect(() => {
        heartbeat = setInterval(() => {
            ping()
        }, 60000)

        return function cleanup() {
            if (heartbeat) clearInterval(heartbeat)
        }
    }, [])

    useEffect(() => {
        try {
            var urlObject = new URL(window.location.href);
            var params = urlObject.searchParams;
            const chal = params.get('challenge')
            if (chal) {
                testChallenge(chal)
            }
        } catch (e) { }
    }, [])


    return useObserver(() => {
        return <>

            <EuiHeader id="header" style={{ color: '#fff' }}>
                <div className="container">
                    <Row style={{ justifyContent: 'space-between' }}>
                        <EuiHeaderSection grow={false}>
                            <Img src="/static/people_logo.svg" />
                        </EuiHeaderSection>

                        <Corner>
                            {ui.meInfo ?
                                <Imgg
                                    style={{ height: 30, width: 30, marginRight: 10 }}
                                    src={ui.meInfo.img || '/static/sphinx.png'}
                                    onClick={() => setShowEditSelf(true)} />
                                :
                                <Button
                                    icon={'account_circle'}
                                    text={'Sign in'}
                                    color='primary'
                                    onClick={() => setShowSignIn(true)}
                                />
                            }
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


            {/* you wanna login modal  */}
            <Modal
                visible={showSignIn}
                close={() => setShowSignIn(false)}
                overlayClick={() => setShowSignIn(false)}
            >
                <SignIn
                    onSuccess={() => {
                        setShowSignIn(false)
                        setShowWelcome(true)
                    }} />
            </Modal >

            {/* you logged in modal  */}
            < Modal
                visible={(ui.meInfo && showWelcome) ? true : false}>
                <div>
                    <Column>
                        <Imgg
                            style={{ height: 128, width: 128, marginBottom: 40 }}
                            src={ui.meInfo?.img || '/static/sphinx.png'} />

                        <T>
                            <div style={{ marginRight: 6 }}>Welcome</div>
                            <Name>{ui.meInfo?.owner_alias}</Name>
                        </T>

                        <Welcome>
                            Your profile is now public.
                            Connect with other people, join tribes and listen your favorite podcast!
                        </Welcome>

                        <Button
                            text={'Continue'}
                            height={48}
                            width={'100%'}
                            color={'primary'}
                            onClick={() => {
                                // switch from welcome modal to edit modal
                                setShowWelcome(false)
                                setShowEditSelf(true)
                            }}
                        />
                    </Column>
                </div>
            </Modal>

            {/* ONLY FOR FIRST TIME USER edit your info modal  */}
            {/* < Modal visible={showInitEditSelf}
                drift={40}
                fill
                close={() => setShowInitEditSelf(false)}
            >
                <div style={{
                    background: '#fff',
                    height: '100%',
                    width: '100%',
                    overflow: 'auto'
                }}>
                    
                </div>
            </Modal> */}


            < Modal visible={(ui.meInfo && showEditSelf) ? true : false}
                drift={40}
                fill
                close={() => setShowEditSelf(false)}
            >
                <PersonViewSlim goBack={() => setShowEditSelf(false)}
                    personId={ui.meInfo?.id}
                />
            </Modal>
        </>
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
const T = styled.div`
display:flex;
font-size: 26px;
line-height: 19px;
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
const Welcome = styled.div`
font-size: 15px;
line-height: 24px;
margin:20px 0 50px;
text-align: center;

/* Text 2 */

color: #3C3F41;
`


const Column = styled.div`
  width:100%;
  display:flex;
  flex-direction:column;
  justify-content:center;
  align-items:center;
  padding: 25px;
  
`
const Imgg = styled.div<ImageProps>`
                        background-image: url("${(p) => p.src}");
                        background-position: center;
                        background-size: cover;
                        width:90px;
                        height:90px;
                        border-radius: 50%;
                        `;

