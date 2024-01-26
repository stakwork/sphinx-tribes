import { Person } from '../../store/main';

export const people: Person[] = [
  {
    id: 1,
    pubkey: 'test_pub_key',
    contact_key: 'test_owner_contact_key',
    alias: 'Vladimir',
    photo_url: '',
    route_hint: 'test_hint:1099567661057',
    price_to_meet: 0,
    url: 'https://proxy2.sphinx.chat',
    description: 'description',
    verification_signature: 'test_verification_signature',
    extras: {
      email: [{ value: 'testEmail@sphinx.com' }],
      liquid: [{ value: 'none' }],
      wanted: []
    },
    owner_alias: 'Vladimir',
    owner_pubkey: 'test_pub_key',
    unique_name: 'vladimir',
    tags: [],
    img: '',
    last_login: 1678263923
  },
  {
    id: 2,
    pubkey: 'test_pub_key_2',
    contact_key: 'test_owner_contact_key_2',
    alias: 'Raphael',
    photo_url: '',
    route_hint: 'test_hint:1099567667689',
    price_to_meet: 0,
    url: 'https://proxy2.sphinx.chat',
    description: 'description2',
    verification_signature: 'test_verification_signature_2',
    extras: {
      email: [{ value: 'testEmail2@sphinx.com' }],
      liquid: [{ value: 'none' }],
      wanted: []
    },
    owner_alias: 'Raphael',
    owner_pubkey: 'test_pub_key_2',
    unique_name: 'raphael',
    tags: [],
    img: '',
    last_login: 16782639234
  }
];

export const person: Person = people[0];
