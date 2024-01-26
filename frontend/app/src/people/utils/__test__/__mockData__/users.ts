import { Person } from 'store/main';

export const users: Person[] = [
  {
    id: 1,
    pubkey: 'test_pub_key',
    alias: '',
    contact_key: 'test_owner_contact_key',
    owner_route_hint: 'test_owner_route_hint',
    unique_name: 'test 1',
    tags: [],
    photo_url: '',
    route_hint: 'test_hint:1099567661057',
    price_to_meet: 0,
    url: 'https://mockApi.com',
    description: 'description',
    verification_signature: 'test_verification_signature',
    extras: {
      email: [{ value: 'testEmail@sphinx.com' }],
      liquid: [{ value: 'none' }],
      wanted: [],
      coding_languages: [{ label: 'Typescript', value: 'Typescript' }]
    },
    owner_alias: 'test 1',
    owner_pubkey: 'test_pub_key',
    img: '/static/avatarPlaceholders/placeholder_34.jpg'
  },
  {
    id: 2,
    pubkey: 'test_pub_key',
    alias: 'test 2',
    contact_key: 'test_owner_contact_key',
    owner_route_hint: 'test_owner_route_hint',
    unique_name: 'test 2',
    tags: [],
    photo_url: '',
    route_hint: 'test_hint:1099567661057',
    price_to_meet: 0,
    url: 'https://mockApi.com',
    description: 'description',
    verification_signature: 'test_verification_signature',
    extras: {
      email: [{ value: 'testEmail@sphinx.com' }],
      liquid: [{ value: 'none' }],
      wanted: [],
      coding_languages: [{ label: 'Java', value: 'Java' }]
    },
    owner_alias: 'test 2',
    owner_pubkey: 'test_pub_key',
    img: '/static/avatarPlaceholders/placeholder_34.jpg'
  },
  {
    id: 3,
    pubkey: 'test_pub_key',
    alias: 'test 3',
    contact_key: 'test_owner_contact_key',
    owner_route_hint: 'test_owner_route_hint',
    unique_name: 'test 3',
    tags: [],
    photo_url: '',
    route_hint: 'test_hint:1099567661057',
    price_to_meet: 0,
    url: 'https://mockApi.com',
    description: 'description',
    verification_signature: 'test_verification_signature',
    extras: {
      email: [{ value: 'testEmail@sphinx.com' }],
      liquid: [{ value: 'none' }],
      wanted: []
    },
    owner_alias: 'test 3',
    owner_pubkey: 'test_pub_key',
    img: '/static/avatarPlaceholders/placeholder_34.jpg'
  }
];
