import { v4 as uuidv4 } from 'uuid';
import { MeInfo } from '../../store/ui';

export const user: MeInfo = {
  id: 1,
  pubkey: 'test_pub_key',
  uuid: uuidv4(),
  contact_key: 'test_owner_contact_key',
  owner_route_hint: 'test_owner_route_hint',
  alias: 'Vladimir',
  photo_url: '',
  github_issues: [],
  route_hint: 'test_hint:1099567661057',
  price_to_meet: 0,
  jwt: 'test_jwt',
  tribe_jwt: 'test_jwt',
  url: 'http://localhost:5002',
  description: 'description',
  verification_signature: 'test_verification_signature',
  extras: {
    email: [{ value: 'testEmail@sphinx.com' }],
    liquid: [{ value: 'none' }],
    wanted: []
  },
  owner_alias: 'Vladimir',
  owner_pubkey: 'test_pub_key',
  img: '/static/avatarPlaceholders/placeholder_34.jpg',
  twitter_confirmed: false,
  isSuperAdmin: false,
  websocketToken: 'test_websocketToken'
};
