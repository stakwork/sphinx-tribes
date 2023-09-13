import { MeInfo } from '../../store/ui';

export const user: MeInfo = {
  id: 1,
  pubkey: 'test_pub_key',
  contact_key: 'test_owner_contact_key',
  alias: 'Vladimir',
  photo_url: '',
  route_hint: 'test_hint:1099567661057',
  price_to_meet: 0,
  jwt: 'test_jwt',
  tribe_jwt: 'test_jwt',
  url: 'https://mockApi.com',
  description: 'description',
  verification_signature: 'test_verification_signature',
  extras: {
    email: [{ value: 'testEmail@sphinx.com' }],
    liquid: [{ value: 'none' }],
    wanted: []
  },
  owner_alias: 'Vladimir',
  owner_pubkey: 'test_pub_key',
  img: '',
  twitter_confirmed: false,
  isSuperAdmin: false
};
