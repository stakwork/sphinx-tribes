/* eslint-disable func-style */
import React from 'react';
import { CodingLanguageLabel } from 'people/interfaces';
import FavoriteButton from '../../../utils/FavoriteButton';
import { Button } from '../../../../components/common';
import { sendToRedirect } from '../../../../helpers';
import { getTwitterLink } from './lib';

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
  const { ticket_url, repo, issue } = props;

  if (ticket_url) {
    return (
      <Button
        text={'Github Ticket'}
        color={'white'}
        endingIcon={'launch'}
        iconSize={14}
        style={{ fontSize: 14, height: 48, width: '100%', marginBottom: 20 }}
        onClick={() => {
          const repoUrl = ticket_url ? ticket_url : `https://github.com/${repo}/issues/${issue}`;
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

type ShareOnTwitterProps = {
  titleString?: string;
  labels?: Array<CodingLanguageLabel>;
  issueCreated?: number;
  ownerPubkey?: string;
};
export const ShareOnTwitter = ({
  titleString,
  labels,
  issueCreated,
  ownerPubkey
}: ShareOnTwitterProps) => {
  if (!(titleString && issueCreated && ownerPubkey)) {
    return null;
  }
  const twitterHandler = () => {
    const twitterLink = getTwitterLink({
      title: titleString,
      labels,
      issueCreated: String(issueCreated),
      ownerPubkey
    });

    sendToRedirect(twitterLink);
  };

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
      onClick={twitterHandler}
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
