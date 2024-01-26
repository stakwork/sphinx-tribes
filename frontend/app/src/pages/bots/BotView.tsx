import React from 'react';
import styled from 'styled-components';
import { useHistory } from 'react-router';
import { observer } from 'mobx-react-lite';
import { useStores } from '../../store';

import { Button, Divider, IconButton } from '../../components/common';
import { useIsMobile } from '../../hooks';
import Bot from './Bot';
import BotBar from './utils/BotBar';
import { BotViewProps } from './interfaces';

const BotList = styled.div`
  display: flex;
  flex-direction: column;
  background: #ffffff;
  width: 265px;
`;

const DBack = styled.div`
  height: 64px;
  min-height: 64px;
  display: flex;
  align-items: center;
  background: #ffffff;
  box-shadow: 0px 0px 6px rgba(0, 0, 0, 0.07);
  z-index: 2;
`;

const Panel = styled.div`
  position: relative;
  background: #ffffff;
  color: #000000;
  margin-bottom: 10px;
  padding: 20px;
  box-shadow: 0px 0px 3px rgb(0 0 0 / 29%);
`;
const Content = styled.div`
  display: flex;
  flex-direction: column;

  width: 100%;
  height: 100%;
  align-items: center;
  color: #000000;
  background: #fff;
`;

const Head = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: 100%;
`;

const Name = styled.div`
  width: 100%;

  text-align: center;

  color: #3c3f41;
  margin: 30px 0;
  font-style: normal;
  font-weight: 600;
  font-size: 24px;
  line-height: 28px;
  text-align: center;
  color: #3c3f41;
`;

const Sleeve = styled.div``;

const Value = styled.div`
  font-weight: bold;
  color: #3c3f41;
`;

const RowWrap = styled.div`
  display: flex;
  justify-content: space-between;
  min-height: 48px;
  align-items: center;
  font-family: Roboto;
  font-style: normal;
  font-weight: normal;
  font-size: 15px;
  line-height: 48px;
  width: 100%;
  color: #8e969c;
`;

interface ImageProps {
  readonly src: string;
}
const Img = styled.div<ImageProps>`
  background-image: url('${(p: any) => p.src}');
  background-position: center;
  background-size: cover;
  width: 150px;
  height: 150px;
  border-radius: 16px;
  position: relative;
  display: flex;
  align-items: flex-end;
  justify-content: flex-end;
`;

const GrowRow = styled.div`
  display: flex;
  justify-content: flex-start;
  flex-wrap: wrap;
  min-height: 48px;
  padding: 10px 0;
  align-items: center;
  font-family: Roboto;
  font-style: normal;
  font-weight: normal;
  font-size: 15px;
  line-height: 48px;
  /* identical to box height, or 320% */

  display: flex;
  align-items: center;

  /* Secondary Text 4 */

  color: #8e969c;
`;

const CodeBadge = styled.div`
display:flex;
justify-content:center;
align-items:center;
margin-right:10px;
height:26px;
color:#5078F2;
background: #DCEDFE;
border-radius: 32px;
font-weight: bold;
font-size: 12px;
line-height: 13px;
padding 0 10px;
margin-bottom:10px;
`;

function BotView(props: BotViewProps) {
  const { botUniqueName, selectBot, loading, goBack } = props;

  const { main } = useStores();

  const history = useHistory();

  const bot: any =
    main.bots && main.bots.length && main.bots.find((f: any) => f.unique_name === botUniqueName);

  const { name, unique_name, description, img, owner_pubkey, owner_alias, tags, price_per_use } =
    bot || {};

  // FOR BOT VIEW
  const bots: any = main.bots && main.bots.length && main.bots.filter((f: any) => !f.hide);

  const isMobile = useIsMobile();

  if (loading) return <div>Loading...</div>;

  const head = (
    <Head>
      {!isMobile && <div style={{ height: 35 }} />}
      <Img src={img || '/static/bot_placeholder.png'} />
      <RowWrap>
        <Name>{name}</Name>
      </RowWrap>
      <RowWrap style={{ marginBottom: 40 }}>
        <BotBar value={unique_name} />
      </RowWrap>

      <RowWrap>
        <div>Price per use</div>
        <Value>{price_per_use}</Value>
      </RowWrap>
      <Divider />
      <RowWrap>
        <div>Creator</div>
        <Value
          style={{ cursor: 'pointer', color: '#5078F2' }}
          onClick={() => history.push(`/p/${owner_pubkey}`)}
        >
          {owner_alias || ''}
        </Value>
      </RowWrap>

      {tags && tags.length > 0 && (
        <div style={{ width: '100%' }}>
          <Divider style={{ marginBottom: 6 }} />
          <GrowRow style={{ paddingBottom: 0 }}>
            {tags.map((c: any, i: number) => (
              <CodeBadge key={i}>{c}</CodeBadge>
            ))}
          </GrowRow>
        </div>
      )}
    </Head>
  );

  function renderMobileView() {
    return (
      <div
        style={{
          display: 'flex',
          flexDirection: 'column',
          width: '100%',
          overflow: 'auto',
          height: '100%'
        }}
      >
        <Panel style={{ paddingBottom: 0, paddingTop: 80 }}>
          <div
            style={{
              position: 'absolute',
              top: 20,
              left: 0,
              display: 'flex',
              justifyContent: 'space-between',
              width: '100%',
              padding: '0 20px'
            }}
          >
            <IconButton onClick={goBack} icon="arrow_back" />
            <div />
          </div>

          {/* profile photo */}
          {head}
        </Panel>

        <Sleeve style={{ padding: 20 }}>
          {description}
          <div style={{ height: 60 }} />
        </Sleeve>
      </div>
    );
  }

  function renderDesktopView() {
    return (
      <div
        style={{
          display: 'flex',
          width: '100%',
          height: '100%'
        }}
      >
        <BotList>
          <DBack>
            <Button color="clear" leadingIcon="arrow_back" text="Back" onClick={goBack} />
          </DBack>

          <div style={{ width: '100%', overflowY: 'auto' }}>
            {bots.map((t: any) => (
              <Bot
                {...t}
                key={t.uuid}
                selected={botUniqueName === t.unique_name}
                hideActions={true}
                small={true}
                select={() => selectBot(t)}
              />
            ))}
          </div>
        </BotList>

        <div
          style={{
            width: 364,
            minWidth: 364,
            overflowY: 'auto',
            position: 'relative',
            background: '#ffffff',
            color: '#000000',
            padding: 40,
            height: '100%',
            borderLeft: '1px solid #F2F3F5',
            borderRight: '1px solid #F2F3F5',
            boxShadow: '1px 0px 6px -2px rgba(0, 0, 0, 0.07)'
          }}
        >
          {/* profile photo */}
          {head}
          {/* Here's where the details go */}
        </div>

        <div
          style={{
            width: 'calc(100% - 628px)',
            minWidth: 250
          }}
        >
          <div
            style={{
              padding: 62,
              height: 'calc(100% - 63px)',
              overflowY: 'auto',
              position: 'relative'
            }}
          >
            <Sleeve
              style={{
                display: 'flex',
                alignItems: 'flex-start',
                flexWrap: 'wrap'
              }}
            >
              {description}
            </Sleeve>
            <div style={{ height: 60 }} />
          </div>
        </div>
      </div>
    );
  }

  return <Content>{isMobile ? renderMobileView() : renderDesktopView()}</Content>;
}

export default observer(BotView);
