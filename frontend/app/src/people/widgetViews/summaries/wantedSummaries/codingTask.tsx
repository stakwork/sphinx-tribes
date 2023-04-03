/* eslint-disable func-style */
import React from 'react';
//import { Img, Assignee } from './style';
//import { colors } from '../../../../config/colors';
//import { extractGithubIssue, extractGithubIssueFromUrl } from '../../../../helpers';
import CodingMobile from './codingMobile';

export default function CodingTask(props: any) {
  const { isMobile } = props;
  /*	

  const { ticketUrl, repo, person, issue, assigneeInfo, isMobile } = props;
  const color = colors['light'];

  const { status } = ticketUrl
    ? extractGithubIssueFromUrl(person, ticketUrl)
    : extractGithubIssue(person, repo, issue);

  let assigneeLabel: any = null;

  function sendToRedirect(url) {
    const el = document.createElement('a');
    el.href = url;
    el.target = '_blank';
    el.click();
  }

				
  if (assigneeInfo) {
    if (!isMobile) {
      assigneeLabel = (
        <div
          style={{
            display: 'flex',
            alignItems: 'center',
            fontSize: 12,
            color: color.grayish.G100,
            marginTop: isMobile ? 20 : 0,
            marginLeft: '-16px'
          }}
        >
          <Img
            src={assigneeInfo.img || '/static/person_placeholder.png'}
            style={{ borderRadius: 30 }}
          />

          <Assignee
            color={color}
            onClick={() => {
              const profileUrl = `https://community.sphinx.chat/p/${assigneeInfo.owner_pubkey}`;
              sendToRedirect(profileUrl);
            }}
            style={{ marginLeft: 3, fontWeight: 500, cursor: 'pointer' }}
          >
            {assigneeInfo.owner_alias}
          </Assignee>
        </div>
      );
    } else {
      assigneeLabel = (
        <div
          style={{
            display: 'flex',
            alignItems: 'center',
            fontSize: 12,
            color: color.grayish.G100,
            marginLeft: '16px'
          }}
        >
          <Img
            src={assigneeInfo.img || '/static/person_placeholder.png'}
            style={{ borderRadius: 30 }}
          />

          <Assignee
            color={color}
            onClick={() => {
              const profileUrl = `https://community.sphinx.chat/p/${assigneeInfo.owner_pubkey}`;
              sendToRedirect(profileUrl);
            }}
            style={{ marginLeft: 3, fontWeight: 500, cursor: 'pointer' }}
          >
            {assigneeInfo.owner_alias}
          </Assignee>
        </div>
      );
    }
  }*/

  if (isMobile) {
    return <CodingMobile />;
  }
}
