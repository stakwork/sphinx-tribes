import { Person } from 'store/main';
import { people } from './persons';
export const person: Person = people[0];
export const assignee: Person = people[1];

export const userAssignedTickets = [
  {
    body: {
      show: true,
      type: 'freelance_job_request',
      price: '400000',
      title: 'Trying this',
      created: 1680130875,
      assignee: {
        img: assignee.img,
        label: `${assignee.owner_alias} (${assignee.unique_name})`,
        value: assignee.owner_pubkey,
        owner_alias: assignee.owner_alias,
        owner_pubkey: assignee.owner_pubkey
      },
      description: 'Just for tesrt',
      wanted_type: 'Desktop app',
      codingLanguage: [
        {
          label: 'Typescript',
          value: 'Typescript'
        },
        {
          label: 'Node',
          value: 'Node'
        }
      ],
      one_sentence_summary: 'Trying this',
      estimate_session_length: 'More than 3 hours'
    },
    person: person
  }
];
