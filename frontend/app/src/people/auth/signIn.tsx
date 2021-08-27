import React, { useState, useEffect } from 'react'
import { useObserver } from 'mobx-react-lite'
import { useStores } from '../../store'
import styled from 'styled-components'
import { useFuse, useIsMobile } from '../../hooks'
import { colors } from '../../colors'
import { Redirect, useHistory, useLocation } from 'react-router-dom'
import { Modal, Button, Divider } from '../../sphinxUI';
import FadeLeft from '../../animated/fadeLeft';
import ConfirmMe from '../confirmMe';


export default function SignIn(props: any) {
    const { main, ui } = useStores()
    // const location = useLocation()

    // function selectPerson(id: number, unique_name: string) {
    //   console.log('selectPerson', id, unique_name)
    //   setSelectedPerson(id)
    //   if (unique_name && window.history.pushState) {
    //     window.history.pushState({}, 'Sphinx Tribes', '/p/' + unique_name);
    //   }
    // }
    const c = colors['light']
    const [showSignIn, setShowSignIn] = useState(false)

    function redirect() {
        let el = document.createElement('a')
        el.target = '_blank'
        el.href = 'https://sphinx.chat/'
        el.click();
    }

    return useObserver(() => {
        return <div>
            {showSignIn ?
                <Column>
                    <ConfirmMe
                        onSuccess={() => {
                            if (props.onSuccess) props.onSuccess()
                            main.getPeople()
                        }} />
                </Column>
                :
                <>
                    <Column>
                        <Imgg src={'/static/sphinx.png'} />

                        <Name>Welcome</Name>

                        <Description>
                            Use Sphinx to login and create or edit your profile.
                        </Description>

                        <Button
                            text={'Login with Sphinx'}
                            height={60}
                            width={'100%'}
                            color={'primary'}
                            onClick={() => setShowSignIn(true)}
                        />
                    </Column>
                    <Divider />
                    <Column style={{ paddingTop: 0 }}>
                        <Description>
                            I don't have Sphinx!
                        </Description>
                        <Button
                            text={'Get Sphinx'}
                            onClick={() => redirect()}
                            height={60}
                            width={'100%'}
                            color={'widget'}
                        />
                    </Column>
                </>
            }
        </div>
    })
}


interface ImageProps {
    readonly src: string;
}

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
                        margin-bottom:20px;
                        width:90px;
                        height:90px;
                        border-radius: 50%;
                        position: relative;
                        `;