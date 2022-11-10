import React, { useState, useEffect } from 'react';
import { useObserver } from 'mobx-react-lite';
import { useStores } from '../../store';
import styled from 'styled-components';
import { EuiHeader, EuiHeaderSection } from '@elastic/eui';
import { useIsMobile } from '../../hooks';
import { colors } from '../../colors';
// import { languageOptions } from '../../localization';
import { useHistory, useLocation } from 'react-router-dom';
import MaterialIcon from '@material/react-material-icon';
import { Modal, Button } from '../../sphinxUI';

import SignIn from '../auth/signIn';
import api from '../../api';
import TorSaveQR from '../utils/torSaveQR';
import Select from 'react-select';
import IconButton from '../../sphinxUI/icon_button';

export default function Header() {
  const { main, ui } = useStores();
  const location = useLocation();
  const history = useHistory();
  const isMobile = useIsMobile();

  const { language } = ui;

  const c = colors['light'];

  const tabs = [
    {
      label: 'Tribes',
      name: 'tribes',
      path: '/t'
    },
    {
      label: 'People',
      name: 'people',
      path: '/p'
    },
    {
      label: 'Bounties',
      name: 'tickets',
      path: '/tickets'
    },
    {
      label: 'Bots',
      name: 'bots',
      path: '/b'
    }
  ];

  const [showWelcome, setShowWelcome] = useState(false);

  async function testChallenge(chal: string) {
    try {
      console.log('testChallenge', chal);
      const me: any = await api.get(`poll/${chal}`);
      console.log('poll succeeded', me);
      if (me && me.pubkey) {
        ui.setMeInfo(me);
        ui.setShowSignIn(false);
        setShowWelcome(true);
      }
    } catch (e) {
      console.log(e);
    }
  }

  function urlRedirect(directPathname) {
    // if route not supported, redirect
    let pass = false;
    const path = directPathname || location.pathname;
    tabs.forEach((t) => {
      if (path.includes(t.path)) pass = true;
    });
    if (!pass) {
      console.log('force fix');
      history.push('/p');
    }
  }

  useEffect(() => {
    const path = location.pathname;
    if (!path.includes('/p') && (ui.selectedPerson || ui.selectingPerson)) {
      ui.setSelectedPerson(0);
      ui.setSelectingPerson(0);
    }
  }, [location.pathname]);

  useEffect(() => {
    (async () => {
      console.log('header deeplink load');
      try {
        const urlObject = new URL(window.location.href);
        let path = location.pathname;
        const params = urlObject.searchParams;
        const chal = params.get('challenge');

        console.log('chal', chal);

        if (chal) {
          // fix url path if "/p" is not included, add challenge to proper url path
          if (!path.includes('/p')) {
            console.log('fix path!');
            path = `/p?challenge=${chal}`;
            history.push(path);
          }
          await testChallenge(chal);
        } else {
          // update self on reload
          await main.getSelf(null);
        }
        urlRedirect(path);
      } catch (e) {
        console.log('e', e);
      }
    })();
  }, []);

  function goToEditSelf() {
    if (ui.meInfo?.id) {
      history.push(`/p/${ui.meInfo.owner_pubkey}`);
      ui.setSelectedPerson(ui.meInfo.id);
      ui.setSelectingPerson(ui.meInfo.id);
    }
  }

  const headerBackground = '#1A242E';

  function renderHeader() {
    if (isMobile) {
      return (
        <EuiHeader
          id="header"
          style={{
            color: '#fff',
            background: headerBackground,
            paddingBottom: 0
          }}
        >
          <div className="container">
            <Row style={{ justifyContent: 'space-between' }}>
              <EuiHeaderSection grow={false}>
                <Img src="/static/people_logo.svg" style={{ width: 190 }} />
              </EuiHeaderSection>

              <Corner>
                <a href={'https://sphinx.chat/'} target="_blank" rel="noreferrer">
                  <Button
                    text={'Get Sphinx'}
                    color="transparent"
                    style={{ marginRight: 14, width: 85 }}
                  />
                </a>

                {ui.meInfo ? (
                  <Imgg
                    style={{
                      height: 30,
                      width: 30,
                      marginRight: 10,
                      border: '1px solid #ffffff55'
                    }}
                    src={ui.meInfo?.img || '/static/person_placeholder.png'}
                    onClick={() => {
                      goToEditSelf();
                    }}
                  />
                ) : (
                  <Button
                    icon={'account_circle'}
                    // text={'Sign in'}
                    style={{ minWidth: 38, width: 38, marginRight: 10, height: 37 }}
                    color="primary"
                    onClick={() => ui.setShowSignIn(true)}
                  />
                )}
              </Corner>
            </Row>

            <MTabs>
              {tabs &&
                tabs.map((t, i) => {
                  const { label } = t;
                  const selected = location.pathname.split('/')[1] === t.path.split('/')[1];

                  return (
                    <MTab
                      key={i}
                      selected={selected}
                      onClick={() => {
                        history.push(t.path);
                      }}
                    >
                      {label}
                    </MTab>
                  );
                })}
            </MTabs>
          </div>
        </EuiHeader>
      );
    }

    // const locOptions = languageOptions.map(lo => {
    //     return {
    //         value: lo,
    //         label: lo
    //     }
    // })

    // desktop version
    return (
      <EuiHeader
        style={{
          color: '#fff',
          width: '100%',
          height: 64,
          padding: '0 20px',
          background: headerBackground
        }}
      >
        <div
          style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            width: '100%'
          }}
        >
          <Row style={{ height: '100%', marginBottom: '-2px' }}>
            <EuiHeaderSection grow={false}>
              <Img src="/static/people_logo.svg" />
            </EuiHeaderSection>

            <Tabs>
              {tabs &&
                tabs.map((t, i) => {
                  const { label } = t;
                  const selected = location.pathname.split('/')[1] === t.path.split('/')[1];

                  return (
                    <Tab
                      key={i}
                      selected={selected}
                      onClick={() => {
                        history.push(t.path);
                      }}
                    >
                      {label}
                    </Tab>
                  );
                })}
            </Tabs>
          </Row>

          <Corner>
            <GetSphinxsBtn href={'https://sphinx.chat/'} target="_blank">
              Get Sphinx
            </GetSphinxsBtn>
            {ui.meInfo ? (
              <LoggedInBtn
                onClick={() => {
                  goToEditSelf();
                }}
              >
                <Imgg src={ui.meInfo?.img || '/static/person_placeholder.png'} />
                {ui.meInfo?.owner_alias}
              </LoggedInBtn>
            ) : (
              <LoginBtn onClick={() => ui.setShowSignIn(true)}>
                <span>Sign in</span>
                <MaterialIcon icon={'login'} style={{ fontSize: 18 }} />
              </LoginBtn>
            )}
          </Corner>
        </div>
      </EuiHeader>
    );
  }

  return useObserver(() => {
    return (
      <>
        {renderHeader()}

        {/* you wanna login modal  */}
        <Modal
          visible={ui.showSignIn}
          close={() => ui.setShowSignIn(false)}
          overlayClick={() => ui.setShowSignIn(false)}
        >
          <SignIn
            onSuccess={() => {
              ui.setShowSignIn(false);
              setShowWelcome(true);
              // if page is not /p, go to /p (people)
              const path = location.pathname;
              if (!path.includes('/p')) history.push('/p');
            }}
          />
        </Modal>

        {/* you logged in modal  */}
        <Modal visible={ui.meInfo && showWelcome ? true : false}>
          <div>
            <Column>
              <Imgg
                style={{ height: 128, width: 128, marginBottom: 40 }}
                src={ui.meInfo?.img || '/static/person_placeholder.png'}
              />
              <div
                style={{
                  position: 'absolute',
                  top: '110px',
                  right: '85px'
                }}
              >
                <img height={'32px'} width={'32px'} src="/static/badges/verfied_mark.png" alt="" />
              </div>

              <T>
                <div style={{ lineHeight: '26px' }}>
                  Welcome <Name>{ui.meInfo?.owner_alias}</Name>
                </div>
              </T>

              <Welcome>
                Your profile is now public. Connect with other people, join tribes and listen your
                favorite podcast!
              </Welcome>

              <IconButton
                text={'Continue'}
                endingIcon={'arrow_forward'}
                height={48}
                width={'100%'}
                color={'primary'}
                onClick={() => {
                  // switch from welcome modal to edit modal
                  setShowWelcome(false);
                  goToEditSelf();
                }}
                hoverColor={'#5881F8'}
                activeColor={'#5078F2'}
                shadowColor={'rgba(97, 138, 255, 0.5)'}
              />
            </Column>
          </div>
        </Modal>

        <Modal
          visible={ui?.torFormBodyQR}
          close={() => {
            ui.setTorFormBodyQR('');
          }}
        >
          <TorSaveQR
            url={ui?.torFormBodyQR}
            goBack={() => {
              ui.setTorFormBodyQR('');
            }}
          />
        </Modal>
      </>
    );
  });
}

const Row = styled.div`
  display: flex;
  align-items: center;
  width: 100%;
`;
const Corner = styled.div`
  display: flex;
  align-items: center;
`;
const T = styled.div`
  display: flex;
  font-size: 26px;
  line-height: 19px;
`;

interface ImageProps {
  readonly src: string;
}
const Img = styled.div<ImageProps>`
  background-image: url('${(p) => p.src}');
  background-position: center;
  background-size: cover;
  height: 37px;
  width: 232px;

  position: relative;
`;

const Name = styled.span`
  font-style: normal;
  font-weight: 500;
  font-size: 26px;
  line-height: 19px;
  /* or 73% */

  /* Text 2 */

  color: #292c33;
`;
const Welcome = styled.div`
  font-size: 15px;
  line-height: 24px;
  margin: 20px 0 50px;
  text-align: center;

  /* Text 2 */

  color: #3c3f41;
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
  width: 90px;
  height: 90px;
  border-radius: 50%;
`;

const Tabs = styled.div`
  display: flex;
  margin-left: 20px;
  height: 100%;
`;

const MTabs = styled.div`
  display: flex;
  margin: 0 20px;
  justify-content: space-around;
`;
interface TagProps {
  selected: boolean;
}
const Tab = styled.div<TagProps>`
  display: flex;
  margin-right: 50px;
  padding: 0 8px;
  color: ${(p) => (p.selected ? '#fff' : '#6B7A8D')};
  cursor: pointer;
  font-weight: 500;
  font-size: 15px;
  line-height: 19px;
  height: 100%;
  align-items: center;
  border-bottom: ${(p) => (p.selected ? '6px solid #618AFF' : '6px solid transparent')};

  &:hover {
    color: #909baa;
  }

  &:active {
    color: #fff;
    border-color: transparent;
  }
`;

const MTab = styled.div<TagProps>`
  display: flex;
  margin: 25px 5px 0;
  color: ${(p) => (p.selected ? '#fff' : '#ffffff99')};
  cursor: pointer;
  height: 30px;
  min-width: 65px;
  font-weight: 500;
  font-size: 15px;
  line-height: 19px;
  justify-content: center;
  border-bottom: ${(p) => (p.selected ? '3px solid #618AFF' : 'none')};
`;

const LoggedInBtn = styled.div`
  height: 40px;
  border-radius: 50%;
  margin-right: 20px;
  border-radius: 32px;
  background: rgba(255, 255, 255, 0.07);
  white-space: nowrap;
  padding: 0 24px 0 50px;
  display: flex;
  align-items: center;
  position: relative;

  ${Imgg} {
    width: 40px;
    height: 40px;
    position: absolute;
    left: 0;
  }

  &:hover {
    background: rgba(255, 255, 255, 0.1);
  }

  &:active {
    background: rgba(255, 255, 255, 0.13);
    ${Imgg} {
      height: 34px;
      width: 34px;
      left: 3px;
    }
  }
`;

const GetSphinxsBtn = styled.a`
  display: flex;
  flex: 1 0 auto;
  background: #618aff;
  box-shadow: 0px 2px 10px rgba(97, 138, 255, 0.5);
  padding: 0 28px;
  height: 40px;
  align-items: center;
  justify-content: center;
  margin-right: 20px;
  border-radius: 32px;
  font-weight: 600;
  font-size: 14px;
  line-height: 19px;
  color: #ffffff;
  font-family: Barlow;

  &:hover {
    background: #5881f8;
    text-decoration: none;
    color: inherit;
  }

  &:active {
    background: #5078f2;
    box-shadow: none;
  }
`;

const LoginBtn = styled.div`
  display: flex;
  flex-wrap: nowrap;
  width: 120px;
  align-items: center;
  cursor: pointer;
  color: #fff;
  font-size: 14px;
  font-weight: 500;
  line-height: 17px;
  cursor: pointer;
  margin-left: 18px;
  span {
    margin-right: 8px;
  }

  &:hover {
    color: #a3c1ff;
  }

  &:active {
    color: #82b4ff;
  }
`;

// const LanguageSelector = styled(Select)`
// min-width:150px;
// margin-right:20px;
// color:#fff !important;

// div {
//     background:#1a242e !important;
//     border:none;
// }

// div.div.div {
//     color:#fff !important;
// }

// #react-select-2-input {
//     color:#fff !important;
// }
// }
// `;
