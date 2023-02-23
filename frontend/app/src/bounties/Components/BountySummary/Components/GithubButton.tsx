import React from 'react';
import { sendToRedirect } from '../../../utils/bountyUtils';
import { Button } from '../../../../sphinxUI';

export default function GithuButton(props: { ticketUrl: string; repo: string; issue: string }) {
  const { ticketUrl, repo, issue } = props;

  return (
    <Button
      text={'Original Ticket'}
      color={'white'}
      endingIcon={'launch'}
      iconSize={14}
      style={{ fontSize: 14, height: 48, width: '100%', marginBottom: 20 }}
      onClick={() => {
        const repoUrl = ticketUrl ? ticketUrl : `https://github.com/${repo}/issues/${issue}`;
        sendToRedirect(repoUrl);
      }}
    />
  );
}
