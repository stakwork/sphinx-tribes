import React, { useCallback, useEffect, useState } from 'react';
import styled from 'styled-components';
import { EuiLoadingSpinner } from '@elastic/eui';
import MaterialIcon from '@material/react-material-icon';
import { EuiGlobalToastList } from '@elastic/eui';
import { observer } from 'mobx-react-lite';
import NoneSpace from '../../people/utils/NoneSpace';
import { Button, Modal, SearchTextInput, Divider } from '../../components/common';
import { useStores } from '../../store';
import { useFuse, useScroll } from '../../hooks';
import { colors } from '../../config/colors';
import FadeLeft from '../../components/animated/FadeLeft';
import { useIsMobile } from '../../hooks';
import Form from '../../components/form/bounty';
import { botSchema } from '../../components/form/schema';
import Bot from './Bot';
import BotView from './BotView';
import BotSecret from './utils/BotSecret';

// avoid hook within callback warning by renaming hooks
const getFuse = useFuse;
const getScroll = useScroll;

const BotText = styled.div`
  width: 259px;
  font-family: Roboto;
  font-style: normal;
  font-weight: normal;
  font-size: 15px;
  line-height: 24px;
  /* or 160% */

  text-align: center;

  color: #3c3f41;

  margin-bottom: 30px;
`;

const Body = styled.div`
  flex: 1;
  height: calc(100% - 105px);
  padding-bottom: 80px;
  width: 100%;
  overflow: auto;
  display: flex;
  flex-direction: column;
`;
const Label = styled.div`
  font-family: Roboto;
  font-style: normal;
  font-weight: bold;
  font-size: 26px;
  line-height: 40px;
  /* or 154% */

  display: flex;
  align-items: center;

  /* Text 2 */

  color: #3c3f41;
`;

const Tabs = styled.div`
  display: flex;
`;

interface TagProps {
  selected: boolean;
}
const Tab = styled.div<TagProps>`
  display: flex;
  padding: 10px 25px;
  margin-right: 35px;
  height: 42px;
  color: ${(p: any) => (p.selected ? '#5D8FDD' : '#5F6368')};
  border: 2px solid #5f636800;
  border-color: ${(p: any) => (p.selected ? '#CDE0FF' : '#5F636800')};
  // border-bottom: ${(p: any) => p.selected && '4px solid #618AFF'};
  cursor: pointer;
  font-weight: 400;
  font-size: 15px;
  line-height: 19px;
  background: ${(p: any) => (p.selected ? '#DCEDFE' : '#3C3F4100')};
  border-radius: 25px;
`;
const Link = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  margin-left: 6px;
  color: #618aff;
  cursor: pointer;
  position: relative;
`;

interface IconProps {
  src: string;
}

const Icon = styled.div<IconProps>`
  background-image: ${(p: any) => `url(${p.src})`};
  width: 220px;
  height: 220px;
  margin: 30px;
  margin-bottom: 0px;
  background-position: center; /* Center the image */
  background-repeat: no-repeat; /* Do not repeat the image */
  background-size: contain; /* Resize the background image to cover the entire container */
  // border-radius:5px;
  overflow: hidden;
`;

function BotBody() {
  const { main, ui } = useStores();
  const [loading, setLoading] = useState(true);
  const [showBotCreator, setShowBotCreator] = useState(false);
  const [showCreate, setShowCreate] = useState(false);
  const [editThisBot, setEditThisBot]: any = useState(null);
  const [showSecret, setShowSecret] = useState('');
  const [selectedWidget, setSelectedWidget] = useState('top');
  const [showDropdown, setShowDropdown] = useState(false);
  const isMyBots = selectedWidget === 'mybots';

  const [toasts, setToasts]: any = useState([]);

  function addToast(name: string) {
    setToasts([
      {
        id: '1',
        title: `Deleted ${name}`
      }
    ]);
  }

  function removeToast() {
    setToasts([]);
  }

  const c = colors['light'];
  const isMobile = useIsMobile();

  const botSelectionAttribute = isMyBots ? 'id' : 'unique_name';

  function selectBot(attr: string) {
    // is mybot
    if (isMyBots) {
      const botSource = isMyBots ? main.myBots : main.bots;
      const thisBot = botSource.find((f: any) => f[botSelectionAttribute] === attr);
      setEditThisBot(thisBot);
      setShowCreate(true);
    } else {
      // is other bot
      ui.setSelectedBot(attr);
      ui.setSelectingBot(attr);
    }
  }

  const loadBots = useCallback(async () => {
    setLoading(true);

    let un = '';
    if (window.location.pathname.startsWith('/b/')) {
      un = window.location.pathname.substr(3);
    }

    await main.getBots(un);
    await main.getMyBots();
    setLoading(false);
  }, [main]);
  async function createOrSaveBot(v: any) {
    v.tags = v.tags && v.tags.map((t: any) => t.value);
    v.price_per_use = parseInt(v.price_per_use);

    const isEdit = v.id ? true : false;

    let b: any = null;

    if (isEdit) {
      // edit
      alert('Bot content cannot be updated right now. Coming soon.');
      return;
    } else {
      // create
      try {
        b = await main.makeBot(v);
        if (b) {
          setShowSecret(b.secret);
          setEditThisBot(b);
        }
        setShowCreate(false);
      } catch (e: any) {
        console.log('e', e);
        alert('Bot could not be saved.');
      }
    }
    loadBots();
  }

  async function deleteBot() {
    try {
      const r = await main.deleteBot(editThisBot.id);
      if (r) {
        addToast(editThisBot.name);
      }
    } catch (e: any) {
      console.log('e', e);
    }

    setEditThisBot(null);
    setShowCreate(false);
    loadBots();
  }

  useEffect(() => {
    loadBots();
  }, [loadBots]);

  const tabs = [
    {
      label: 'Top',
      name: 'top'
    }
  ];

  if (ui.meInfo) {
    tabs.push({
      label: 'My Bots',
      name: 'mybots'
    });
  }

  function redirect() {
    const el = document.createElement('a');
    el.target = '_blank';
    el.href = 'https://github.com/stakwork/sphinx-relay/blob/master/docs/deep/bots.md';
    el.click();
  }

  const botSource = isMyBots ? main.myBots : main.bots;

  const bs = getFuse(botSource, ['name', 'description']);
  const { n } = getScroll();
  let bots = bs.slice(0, n);

  if (!isMyBots) {
    // hide bots if not looking at your own
    bots = (bots && bots.filter((f: any) => !f.hide)) || [];
  }

  if (loading) {
    return (
      <Body style={{ justifyContent: 'center', alignItems: 'center' }}>
        <EuiLoadingSpinner size="xl" />
      </Body>
    );
  }

  const widgetLabel = selectedWidget && tabs.find((f: any) => f.name === selectedWidget);

  function renderDesktop() {
    return (
      <Body
        style={{
          background: '#f0f1f3',
          height: 'calc(100% - 65px)'
        }}
      >
        {!ui.meInfo && (
          <div style={{ marginTop: 50 }}>
            <NoneSpace
              buttonText={'Get Started'}
              buttonIcon={'arrow_forward'}
              action={() => ui.setShowSignIn(true)}
              img={'bots_nonespace.png'}
              text={'Discover Bots on Sphinx'}
              sub={'Spice up your Sphinx experience with our diverse range of Sphinx bots'}
              style={{ height: 400 }}
            />
            <Divider />
          </div>
        )}

        <div
          style={{
            width: '100%',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'flex-start',
            padding: 20,
            height: 62
          }}
        >
          <div style={{ display: 'flex', alignItems: 'center' }}>
            <Label style={{ marginRight: 46 }}>Explore</Label>

            <Tabs>
              {tabs &&
                tabs.map((t: any, i: number) => {
                  const { label } = t;
                  const selected = selectedWidget === t.name;

                  return (
                    <Tab
                      key={i}
                      selected={selected}
                      onClick={() => {
                        setSelectedWidget(t.name);
                      }}
                    >
                      {label}
                    </Tab>
                  );
                })}
            </Tabs>
          </div>

          <div style={{ display: 'flex', alignItems: 'center' }}>
            {ui.meInfo && (
              <Button
                text={'Add a Bot'}
                leadingIcon={'add'}
                height={40}
                color="primary"
                onClick={() => setShowBotCreator(true)}
              />
            )}

            <SearchTextInput
              name="search"
              type="search"
              placeholder="Search"
              value={ui.searchText}
              style={{
                width: 204,
                height: 40,
                background: c.grayish.G400,
                marginLeft: 20
              }}
              onChange={(e: any) => {
                ui.setSearchText(e);
              }}
            />
          </div>
        </div>

        <>
          <div
            style={{
              width: '100%',
              display: 'flex',
              flexWrap: 'wrap',
              height: '100%',
              justifyContent: 'flex-start',
              alignItems: 'flex-start',
              padding: 20
            }}
          >
            {bots.map((t: any) => (
              <Bot
                {...t}
                key={t.uuid}
                small={false}
                selected={ui.selectedBot === t.uuid}
                select={() => {
                  selectBot(t[botSelectionAttribute]);
                }}
              />
            ))}
          </div>
          <div style={{ height: 100 }} />
        </>

        {/* selected view */}
        <FadeLeft
          withOverlay={isMobile}
          drift={40}
          overlayClick={() => ui.setSelectingBot('')}
          style={{
            position: 'absolute',
            top: isMobile ? 0 : 64,
            right: 0,
            zIndex: 10000,
            width: '100%'
          }}
          isMounted={ui.selectingBot ? true : false}
          dismountCallback={() => ui.setSelectedBot('')}
        >
          <BotView
            goBack={() => ui.setSelectingBot('')}
            botUniqueName={ui.selectedBot}
            loading={loading}
            selectBot={(b: any) => selectBot(b[botSelectionAttribute])}
            botView={true}
          />
        </FadeLeft>
      </Body>
    );
  }

  function renderMobile() {
    return (
      <Body>
        {!ui.meInfo && (
          <div style={{ marginTop: 50 }}>
            <NoneSpace
              buttonText={'Get Started'}
              buttonIcon={'arrow_forward'}
              action={() => ui.setShowSignIn(true)}
              img={'bots_nonespace.png'}
              text={'Discover Bots on Sphinx'}
              sub={'Spice up your Sphinx experience with our diverse range of Sphinx bots'}
              style={{ background: '#fff', height: 400 }}
            />
            <Divider />
          </div>
        )}
        <div
          style={{
            width: '100%',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'flex-start',
            padding: 20,
            height: 62,
            marginBottom: 20
          }}
        >
          <Label style={{ fontSize: 20 }}>
            Explore
            <Link onClick={() => setShowDropdown(true)}>
              <div>{widgetLabel && widgetLabel.label}</div>
              <MaterialIcon icon={'expand_more'} style={{ fontSize: 18, marginLeft: 5 }} />

              {showDropdown && (
                <div
                  style={{
                    position: 'absolute',
                    top: 0,
                    left: 0,
                    zIndex: 10,
                    background: '#fff'
                  }}
                >
                  {tabs &&
                    tabs.map((t: any, i: number) => {
                      const { label } = t;
                      const selected = selectedWidget === t.name;

                      return (
                        <Tab
                          key={i}
                          style={{ borderRadius: 0, margin: 0 }}
                          selected={selected}
                          onClick={(e: any) => {
                            e.stopPropagation();
                            setShowDropdown(false);
                            setSelectedWidget(t.name);
                          }}
                        >
                          {label}
                        </Tab>
                      );
                    })}
                </div>
              )}
            </Link>
          </Label>

          <SearchTextInput
            small
            name="search"
            type="search"
            placeholder="Search"
            value={ui.searchText}
            style={{ width: 164, height: 40, border: '1px solid #DDE1E5', background: '#fff' }}
            onChange={(e: any) => {
              ui.setSearchText(e);
            }}
          />
        </div>

        <Button
          text={'Add a Bot'}
          leadingIcon={'add'}
          style={{ height: 48, minHeight: 48, margin: 20, marginTop: 0 }}
          color="primary"
          onClick={() => setShowBotCreator(true)}
        />
        <div style={{ width: '100%' }}>
          {bots.map((t: any) => (
            <Bot
              {...t}
              key={t.id}
              selected={ui.selectedBot === t.id}
              small={isMobile}
              select={() => selectBot(t[botSelectionAttribute])}
            />
          ))}
        </div>
        <FadeLeft
          withOverlay
          drift={40}
          overlayClick={() => ui.setSelectingBot('')}
          style={{ position: 'absolute', top: 0, right: 0, zIndex: 10000, width: '100%' }}
          isMounted={ui.selectingBot ? true : false}
          dismountCallback={() => ui.setSelectedBot('')}
        >
          <BotView
            goBack={() => ui.setSelectingBot('')}
            botUniqueName={ui.selectedBot}
            loading={loading}
            selectBot={(b: any) => selectBot(b[botSelectionAttribute])}
            botView={true}
          />
        </FadeLeft>
      </Body>
    );
  }

  const renderContent = isMobile ? renderMobile() : renderDesktop();

  let initialValues: any = {};

  // set initials here
  if (editThisBot) {
    initialValues = { ...editThisBot };

    initialValues.tags =
      initialValues.tags &&
      initialValues.tags.map((o: any) => ({
        value: o.value || o,
        label: o.value || o
      }));
  }

  const botEditHeader = editThisBot?.secret && (
    <div style={{ marginBottom: -50 }}>
      <BotSecret {...editThisBot} />
    </div>
  );

  const botEditHeaderFull = editThisBot?.secret && (
    <div>
      <BotSecret {...editThisBot} full />
    </div>
  );

  return (
    <>
      {renderContent}

      <div style={{ overflowY: 'auto' }}>
        <Modal
          style={{ overflowY: 'auto', height: '100%' }}
          close={() => {
            setShowBotCreator(false);
          }}
          visible={showBotCreator}
        >
          <div
            style={{
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              justifyContent: 'center'
            }}
          >
            <Icon src={'/static/bots_create.svg'} />

            <BotText>Share your awesome bot with other Sphinx chat users!</BotText>
            <Button
              text={'Add a Bot'}
              color={'primary'}
              leadingIcon={'add'}
              height={50}
              width={200}
              onClick={() => {
                setShowCreate(true);
                setShowBotCreator(false);
              }}
            />
            <div style={{ height: 20 }} />
            <Button
              text={'Learn about Bots'}
              leadingIcon={'open_in_new'}
              style={{ marginBottom: 20 }}
              height={50}
              width={200}
              onClick={() => redirect()}
            />
          </div>
        </Modal>

        <Modal
          visible={showCreate}
          close={() => {
            setShowCreate(false);
            setEditThisBot(null);
          }}
          style={{ height: '100%' }}
          envStyle={{
            height: '100%',
            borderRadius: 0,
            width: '100%',
            maxWidth: 450,
            paddingTop: editThisBot?.secret && 60
          }}
        >
          <div style={{ height: '100%', overflowY: 'auto', padding: 20 }}>
            {botEditHeader}
            <Form
              loading={loading}
              close={() => {
                setShowCreate(false);
                setEditThisBot(null);
              }}
              delete={editThisBot && deleteBot}
              onSubmit={createOrSaveBot}
              submitText={editThisBot ? 'Save' : 'Submit'}
              schema={botSchema}
              initialValues={initialValues}
            />
          </div>
        </Modal>

        <Modal
          visible={!showCreate && showSecret ? true : false}
          close={() => {
            setShowSecret('');
            setEditThisBot(null);
          }}
        >
          <div>{botEditHeaderFull}</div>
        </Modal>
      </div>

      <EuiGlobalToastList toasts={toasts} dismissToast={removeToast} toastLifeTimeMs={1000} />
    </>
  );
}

export default observer(BotBody);
