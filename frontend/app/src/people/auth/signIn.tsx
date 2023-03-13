import React, { useState } from 'react';
import { useObserver } from 'mobx-react-lite';
import { useStores } from '../../store';
import styled from 'styled-components';
import { useIsMobile } from '../../hooks';
import { colors } from '../../config/colors';
import { Divider } from '../../components/common';
import ConfirmMe from '../confirmMe';
import AuthQR from './authQR';
import IconButton from '../../components/common/icon_button';

export default function SignIn(props: any) {
  const { main } = useStores();

  const c = colors['light'];
  const [showSignIn, setShowSignIn] = useState(false);

  function redirect() {
    const el = document.createElement('a');
    el.target = '_blank';
    el.href = 'https://sphinx.chat/';
    el.click();
  }

  const isMobile = useIsMobile();

  return useObserver(() => {
    return (
      <div>
        {showSignIn ? (
          <Column>
            <ConfirmMe
              onSuccess={() => {
                if (props.onSuccess) props.onSuccess();
                main.getPeople({ resetPage: true });
              }}
            />
          </Column>
        ) : (
          <>
            <Column>
              {isMobile && <Imgg src={'/static/sphinx.png'} />}

              <Name>Welcome</Name>

              <Description>Use Sphinx to login and create or edit your profile.</Description>

              {!isMobile && (
                <AuthQR
                  onSuccess={() => {
                    if (props.onSuccess) props.onSuccess();
                    main.getPeople({ resetPage: true });
                  }}
                  style={{ marginBottom: 20 }}
                />
              )}

              <IconButton
                text={'Login with Sphinx'}
                height={48}
                endingIcon={'exit_to_app'}
                width={210}
                style={{ marginTop: 20 }}
                color={'primary'}
                onClick={() => setShowSignIn(true)}
                hovercolor={'#5881F8'}
                activecolor={'#5078F2'}
                shadowcolor={'rgba(97, 138, 255, 0.5)'}
              />
            </Column>
            <Divider />
            <Column style={{ paddingTop: 0 }}>
              <Description>I don't have Sphinx!</Description>
              <IconButton
                text={'Get Sphinx'}
                endingIcon={'launch'}
                width={210}
                height={48}
                buttonType={'text'}
                style={{ color: '#83878b', marginTop: '10px', border: '1px solid #83878b' }}
                onClick={() => redirect()}
                hovercolor={'#fff'}
                activecolor={'#fff'}
                textStyle={{
                  color: '#000',
                  fontSize: '16px',
                  fontWeight: '600'
                }}
              />
            </Column>
          </>
        )}
      </div>
    );
  });
}

interface ImageProps {
  readonly src: string;
}

const Name = styled.div`
  font-style: normal;
  font-weight: 500;
  font-size: 26px;
  line-height: 19px;
  font-family: Barlow;
  /* or 73% */

  text-align: center;

  /* Text 2 */

  color: #292c33;
`;

const Description = styled.div`
  font-size: 17px;
  line-height: 20px;
  text-align: center;
  margin: 20px 0;
  font-family: Barlow;

  /* Main bottom icons */

  color: #5f6368;
`;

const Column = styled.div`
  width: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 25px;
`;
const Imgg = styled.div<ImageProps>`
  background-image: url('${(p) => p.src}');
  background-position: center;
  background-size: cover;
  margin-bottom: 20px;
  width: 90px;
  height: 90px;
  border-radius: 50%;
  position: relative;
`;
