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
  const bountyUrl = new URL('https://community.sphinx.chat/tickets');

  for (const key in bountyParams) {
    bountyUrl.searchParams.append(key, bountyParams[key]);
  }

  const params = {
    text: `Hey, I created a new ticket on Sphinx community. ${title}`,
    url: `${bountyUrl}`,
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
