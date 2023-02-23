import React from 'react';
import { sendToRedirect } from '../../../utils/bountyUtils';
import { Button } from '../../../../sphinxUI';

export default function GithuButton(props: {
  title: string;
  owner_idURL: string;
  createdURL: string;
  labels: string[];
}) {
  const { title, owner_idURL, createdURL, labels } = props;

  return (
    <Button
      text={'Share to Twitter'}
      color={'white'}
      icon={'share'}
      iconSize={18}
      iconStyle={{ left: 14 }}
      style={{
        fontSize: 14,
        height: 48,
        width: '100%',
        marginBottom: 20,
        paddingLeft: 5
      }}
      onClick={() => {
        const twitterLink = `https://twitter.com/intent/tweet?text=Hey, I created a new ticket on Sphinx community.%0A${title} %0A&url=https://community.sphinx.chat/p?owner_id=${owner_idURL}%26created${createdURL} %0A%0A&hashtags=${
          labels && labels.map((x: any) => x.label)
        },sphinxchat`;
        sendToRedirect(twitterLink);
      }}
    />
  );
}
