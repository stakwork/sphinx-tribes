import { CodingLanguageLabel } from 'people/interfaces';

type Props = {
  title: string;
  labels?: Array<CodingLanguageLabel>;
  ownerPubkey: string;
  issueCreated: string;
};
export const getTwitterLink = ({ title, issueCreated, ownerPubkey, labels }: Props) => {
  const bountyParams = {
    owner_id: ownerPubkey,
    created: issueCreated
  };
  const origin = window.location.origin.includes('localhost')
    ? 'https://community.sphinx.chat'
    : window.location.origin;

  const bountyUrl = new URL('/bounties', origin);

  for (const key in bountyParams) {
    bountyUrl.searchParams.append(key, bountyParams[key]);
  }

  const params = {
    text: `I just created a bounty on Sphinx Community: ${title}\n`,
    url: `${bountyUrl}\n`,
    hashtags: [...(labels ?? []), { label: 'sphinxchat', value: '' }]
      .map((x: CodingLanguageLabel) => x.label)
      .join(',')
  };

  const link = new URL('https://twitter.com/intent/tweet');

  for (const key in params) {
    link.searchParams.append(key, params[key]);
  }

  return link.toString();
};
