/* eslint-disable func-style */
import React from 'react';
import FavoriteButton from '../../../utils/favoriteButton';
import { Button } from '../../../../components/common';
import { sendToRedirect } from '../../../../helpers';

export const Heart = () => <FavoriteButton />;

export const AddToFavorites = (props: any) => {
  if (props.tribe && props.tribe !== 'none') {
    return (
      <Button
        text={'Add to Favorites'}
        color={'white'}
        icon={'favorite_outline'}
        iconSize={18}
        iconStyle={{ left: 14 }}
        style={{
          fontSize: 14,
          height: 48,
          width: '100%',
          marginBottom: 20,
          paddingLeft: 5
        }}
      />
    );
  }
  return <></>;
};

export const ViewGithub = (props: any) => {
  const { ticketUrl, repo, issue } = props;

  if (ticketUrl) {
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

  return <></>;
};

export const CopyLink = (props: any) => {
  const { isCopied, handleCopyUrl } = props;

  return (
    <Button
      text={isCopied ? 'Copied' : 'Copy Link'}
      color={'white'}
      icon={'content_copy'}
      iconSize={18}
      iconStyle={{ left: 14 }}
      style={{
        fontSize: 14,
        height: 48,
        width: '100%',
        marginBottom: 20,
        paddingLeft: 5
      }}
      onClick={handleCopyUrl}
    />
  );
};

export const ShareOnTwitter = (props: any) => {
  const { titleString, labels, createdURL, owner_idURL } = props;

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
        const twitterLink = `https://twitter.com/intent/tweet?text=Hey, I created a new ticket on Sphinx community.%0A${titleString} %0A&url=https://community.sphinx.chat/p?owner_id=${owner_idURL}%26created${createdURL} %0A%0A&hashtags=${
          labels && labels.map((x: any) => x.label)
        },sphinxchat`;
        sendToRedirect(twitterLink);
      }}
    />
  );
};

export const ViewTribe = (props: any) => {
  const { tribe, tribeInfo } = props;

  if (tribe && tribe !== 'none') {
    return (
      <Button
        text={'View Tribe'}
        color={'white'}
        leadingImgUrl={tribeInfo?.img || ' '}
        endingIcon={'launch'}
        iconSize={14}
        imgStyle={{ position: 'absolute', left: 10 }}
        style={{ fontSize: 14, height: 48, width: '100%', marginBottom: 20 }}
        onClick={() => {
          const profileUrl = `https://community.sphinx.chat/t/${tribe}`;
          sendToRedirect(profileUrl);
        }}
      />
    );
  }

  return <></>;
};
