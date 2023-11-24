import { toJS } from 'mobx';
import { user } from '../../__test__/__mockData__/user';
import { uiStore } from '../ui';
import { MainStore } from '../main';
import { localStorageMock } from '../../__test__/__mockData__/localStorage';

const mockFetch = jest.fn();
const mockHeaders = jest.fn();

const origFetch = global.fetch;

beforeAll(() => {
    global.fetch = mockFetch;
    global.Headers = mockHeaders;
})

afterAll(() => {
    global.fetch = origFetch;
    jest.clearAllMocks();
});

describe('Main store', () => {
  beforeEach(async () => {
    uiStore.setMeInfo(user)
    localStorageMock.setItem('ui', JSON.stringify(uiStore))
  });

  it("should call endpoint on saveBounty", () => {
    const mainStore = new MainStore();
    const bounty = {
      title: "title",
      description: "description",
      amount: 100,
      owner_pubkey: user.owner_pubkey,
      owner_alias: user.alias,
      owner_contact_key: user.contact_key,
      owner_route_hint: user.route_hint ?? '',
      extras: user.extras,
      price_to_meet: user.price_to_meet,
      img: user.img,
      tags: [],
      route_hint: user.route_hint
    }
    mainStore.saveBounty(bounty)


    expect(mockFetch).toBeCalledTimes(1);
    expect(mockFetch).toBeCalledWith(`https://people.sphinx.chat/gobounties?token=test_jwt`, {
        method: 'POST',
        mode: 'cors',
        headers: {
            'x-jwt': 'test_jwt',
            'Content-Type': 'application/json'
        },
        body: expect.anything(),
    });
  });


  it('should save user profile', async () => {
    const mainStore = new MainStore();
    const person = {
      owner_pubkey: user.owner_pubkey,
      owner_alias: user.alias,
      owner_contact_key: user.contact_key,
      owner_route_hint: user.route_hint ?? '',
      description: user.description,
      extras: user.extras,
      price_to_meet: user.price_to_meet,
      img: user.img,
      tags: [],
      route_hint: user.route_hint
    }
    mainStore.saveProfile(person)

    expect(toJS(uiStore.meInfo)).toEqual(user);
    expect(localStorageMock.getItem('ui')).toEqual(JSON.stringify(uiStore));
  });
});

