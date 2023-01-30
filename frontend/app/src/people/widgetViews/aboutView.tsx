import React, { useState } from 'react';
import styled from 'styled-components';
import { Divider } from '../../sphinxUI';
import QrBar from '../utils/QrBar';
import ReactMarkdown from 'react-markdown';
import { useHistory, useLocation } from 'react-router-dom';
import { renderMarkdown } from '../utils/renderMarkdown';

export default function AboutView(props: any) {
  const history = useHistory();
  const { price_to_meet, extras, twitter_confirmed, owner_pubkey } = props;
  const { twitter, github, coding_languages, tribes, repos, lightning, amboss } = extras || {};
  let tag = '';
  let githubTag = '';
  let lightningAddress = '';
  let ambossAddress = '';

  let { description } = props;

  // backend is adding 'description' to empty descriptions, short term fix
  if (description === 'description') description = '';

  const [expand, setExpand] = useState(false);

  if (twitter && twitter[0] && twitter[0].value) tag = twitter[0].value;
  if (github && github[0] && github[0].value) githubTag = github[0].value;
  if (lightning && lightning[0] && lightning[0].value) lightningAddress = lightning[0].value;
  if (amboss && amboss[0] && amboss[0].value) ambossAddress = amboss[0].value;

  const descriptionIsLong = description && description.length && description.length > 120;

  return (
    <Wrap>
      <D>
        {expand ? (
          <DExpand>{renderMarkdown(description)}</DExpand>
        ) : (
          <DCollapsed>{renderMarkdown(description)}</DCollapsed>
        )}
        {descriptionIsLong && (
          <SM onClick={() => setExpand(!expand)}>SHOW {expand ? 'LESS' : 'MORE'}</SM>
        )}
      </D>

      <Divider />
      <Row>
        <div>Price to Connect:</div>
        <div style={{ fontWeight: 'bold', color: '#000' }}>{price_to_meet}</div>
      </Row>

      <Divider />

      <QrBar value={owner_pubkey} />

      {tag && (
        <>
          <Divider />
          <Row>
            {/* <T>For Normies</T> */}
            <I>
              <div style={{ width: 4 }} />
              <Icon source={`/static/twitter2.png`} />
              <Tag>@{tag}</Tag>
              {twitter_confirmed ? (
                <Badge>VERIFIED</Badge>
              ) : (
                <Badge style={{ background: '#b0b7bc' }}>PENDING</Badge>
              )}
            </I>
          </Row>
        </>
      )}

      {githubTag && (
        <>
          <Divider />

          <Row style={{ justifyContent: 'flex-start', fontSize: 14 }}>
            <Img src={'/static/github_logo.png'} />
            <a href={`https://github.com/${githubTag}`} target="_blank" rel="noreferrer">
              {githubTag}
            </a>
          </Row>
        </>
      )}

      {coding_languages && coding_languages.length > 0 && (
        <>
          <Divider />
          <GrowRow style={{ paddingBottom: 0 }}>
            {coding_languages.map((c, i) => {
              return <CodeBadge key={i}>{c.label}</CodeBadge>;
            })}
          </GrowRow>
        </>
      )}

      {repos && repos.length > 0 && (
        <>
          <Divider />
          <T style={{ height: 20 }}>My Repos</T>
          <Grow>
            {repos.map((r, i) => {
              return (
                <ItemRow key={`${i}myrepo`} style={{ width: 'fit-content' }}>
                  <Img src={'/static/github_logo.png'} style={{ opacity: 0.6 }} />
                  <a href={`https://github.com/${r?.label}`} target="_blank" rel="noreferrer">
                    {r?.label}
                  </a>
                </ItemRow>
              );
            })}
          </Grow>
        </>
      )}

      {tribes && tribes.length > 0 && (
        <>
          <Divider />
          <T style={{ height: 20 }}>My Tribes</T>
          <Grow>
            {tribes.map((t, i) => {
              return (
                <ItemRow key={`${i}mytribe`} onClick={() => history.push(`/t/${t?.unique_name}`)}>
                  <Img src={t?.img || '/static/sphinx.png'} />
                  <div>{t?.name}</div>
                </ItemRow>
              );
            })}
          </Grow>
        </>
      )}

      {lightningAddress && (
        <>
          <Divider />
          <T style={{ height: 20 }}>Lightning Address</T>
          <Grow>
            <ItemRow style={{ width: 'fit-content' }}>{lightningAddress}</ItemRow>
          </Grow>
        </>
      )}

      {ambossAddress && (
        <>
          <Divider />
          <T style={{ height: 20 }}>Amboss Address</T>
          <Grow>
            <ItemRow style={{ width: 'fit-content' }}>{ambossAddress}</ItemRow>
          </Grow>
        </>
      )}
    </Wrap>
  );
}
const Badge = styled.div`
display:flex;
justify-content:center;
align-items:center;
margin-left:10px;
height:20px;
color:#ffffff;
background: #1DA1F2;
border-radius: 32px;
font-weight: bold;
font-size: 8px;
line-height: 9px;
padding 0 10px;
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
const ItemRow = styled.div`
  display: flex;
  align-items: center;
  cursor: pointer;
  margin-bottom: 5px;
  &:hover {
    color: #000;
  }
`;
const Wrap = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
  overflow-x: hidden;
`;
const I = styled.div`
  display: flex;
  align-items: center;
`;

const Tag = styled.div`
  font-family: Roboto;
  font-style: normal;
  font-weight: normal;
  font-size: 14px;
  line-height: 26px;
  /* or 173% */

  display: flex;
  align-items: center;

  /* Main bottom icons */

  color: #5f6368;
`;

const Row = styled.div`
  display: flex;
  justify-content: space-between;
  height: 48px;
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

const Grow = styled.div`
  display: flex;
  justify-content: flex-start;
  flex-direction: column;
  min-height: 28px;
  font-family: Roboto;
  font-style: normal;
  font-weight: normal;
  font-size: 14px;
  line-height: 28px;
  margin-bottom: 8px;
  color: #8e969c;
`;
const T = styled.div`
  font-family: Roboto;
  font-style: normal;
  font-weight: bold;
  font-size: 10px;
  line-height: 26px;
  /* or 260% */

  letter-spacing: 0.3px;
  text-transform: uppercase;

  /* Text 2 */

  color: #3c3f41;
  margin-top: 5px;
  margin-bottom: 5px;
`;

const D = styled.div`
  margin: 5px 0 15px 0;
  font-family: Roboto;
  font-style: normal;
  font-weight: normal;
  font-size: 15px;
  line-height: 20px;
  /* or 133% */
  /* Main bottom icons */

  color: #5f6368;
`;

const SM = styled.div`
  display: flex;
  justify-content: flex-end;
  cursor: pointer;
  color: #618aff;
  font-size: 11px;
  letter-spacing: 0.3px;
  margin: 10px 0 4px;
`;

const DCollapsed = styled.div`
  color: #5f6368;
  line-height: 20px;
  padding-bottom: 5px;
  text-overflow: ellipsis;
  max-height: 180px;
  display: block;
  overflow: hidden;
  /* -webkit-line-clamp: 3;
-webkit-box-orient: vertical; */
`;

const DExpand = styled.div`
  color: #5f6368;
  line-height: 20px;
`;

interface IconProps {
  source: string;
}

const Icon = styled.div<IconProps>`
  background-image: ${(p) => `url(${p.source})`};
  width: 16px;
  height: 13px;
  margin-right: 8px;
  background-position: center; /* Center the image */
  background-repeat: no-repeat; /* Do not repeat the image */
  background-size: contain; /* Resize the background image to cover the entire container */
  border-radius: 5px;
  overflow: hidden;
`;

interface ImageProps {
  readonly src?: string;
}
const Img = styled.div<ImageProps>`
  background-image: url('${(p) => p.src}');
  background-position: center;
  background-size: cover;
  position: relative;
  width: 18px;
  height: 18px;
  margin-left: 2px;
  margin-right: 10px;
  border-radius: 5px;
`;
