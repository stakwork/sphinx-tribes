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
});

afterAll(() => {
  global.fetch = origFetch;
  jest.clearAllMocks();
});

const getOrganizationUsersEndpoint = (orgUUID: string): string => {
  return `https://community.sphinx.chat/organizations/users/${orgUUID}`;
};

describe('Main store', () => {
  beforeEach(async () => {
    uiStore.setMeInfo(user);
    localStorageMock.setItem('ui', JSON.stringify(uiStore));
  });

  it('should call endpoint on saveBounty', () => {
    const mainStore = new MainStore();
    mainStore.saveBounty = jest
      .fn()
      .mockReturnValueOnce(Promise.resolve({ status: 200, message: 'success' }));
    const bounty = {
      body: {
        title: 'title',
        description: 'description',
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
    };
    mainStore.saveBounty(bounty);
    expect(mainStore.saveBounty).toBeCalledWith({
      body: bounty.body
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
    };
    mainStore.saveProfile(person);

    expect(toJS(uiStore.meInfo)).toEqual(user);
    expect(localStorageMock.getItem('ui')).toEqual(JSON.stringify(uiStore));
  });
  it('should call endpoint on addOrganizationUser', async () => {
    const mainStore = new MainStore();

    const mockApiResponse = { status: 200, message: 'success' };

    mockFetch.mockReturnValueOnce(Promise.resolve(mockApiResponse));

    const organizationUser = {
      owner_pubkey: user.owner_pubkey || '',
      org_uuid: 'cmas9gatu2rvqiev4ur0'
    };

    const expectedHeaders = {
      'Content-Type': 'application/json',
      'x-jwt': 'test_jwt'
    };

    await mainStore.addOrganizationUser(organizationUser);

    expect(mockFetch).toBeCalledWith(
      'https://community.sphinx.chat/organizations/users/cmas9gatu2rvqiev4ur0',
      {
        method: 'POST',
        headers: expectedHeaders,
        body: JSON.stringify(organizationUser),
        mode: 'cors'
      }
    );
  });

  it('should call endpoint on getOrganizationUsers', async () => {
    const mainStore = new MainStore();

    const mockApiResponse = {
      status: 200,
      json: jest
        .fn()
        .mockResolvedValue([{ uuid: 'cm3eulatu2rvqi9o75ug' }, { uuid: 'cldl1g04nncmf23du7kg' }])
    };

    mockFetch.mockReturnValueOnce(Promise.resolve(mockApiResponse));

    const orgUUID = 'cmas9gatu2rvqiev4ur0';

    const endpoint = getOrganizationUsersEndpoint(orgUUID);

    const users = await mainStore.getOrganizationUsers(orgUUID);

    expect(mockFetch).toBeCalledWith(endpoint, expect.anything());

    expect(users).toEqual([{ uuid: 'cm3eulatu2rvqi9o75ug' }, { uuid: 'cldl1g04nncmf23du7kg' }]);
  });

  it('should call endpoint on getOrganizationUser', async () => {
    const mainStore = new MainStore();

    const mockApiResponse = {
      status: 200,
      json: jest.fn().mockResolvedValue({
        uuid: 'cm3eulatu2rvqi9o75ug'
      })
    };

    mockFetch.mockReturnValueOnce(Promise.resolve(mockApiResponse));

    const userUUID = 'cm3eulatu2rvqi9o75ug';

    const organizationUser = await mainStore.getOrganizationUser(userUUID);

    expect(mockFetch).toBeCalledWith(
      `https://community.sphinx.chat/organizations/foruser/${userUUID}`,
      expect.objectContaining({
        method: 'GET',
        mode: 'cors',
        headers: {
          'x-jwt': 'test_jwt',
          'Content-Type': 'application/json'
        }
      })
    );

    expect(organizationUser).toEqual({
      uuid: 'cm3eulatu2rvqi9o75ug'
    });
  });

  it('should call endpoint on getOrganizationUsersCount', async () => {
    const mainStore = new MainStore();

    const mockApiResponse = {
      status: 200,
      json: jest.fn().mockResolvedValue({
        count: 2
      })
    };

    mockFetch.mockReturnValueOnce(Promise.resolve(mockApiResponse));

    const orgUUID = 'cmas9gatu2rvqiev4ur0';

    const organizationsCount = await mainStore.getOrganizationUsersCount(orgUUID);

    expect(mockFetch).toBeCalledWith(
      `https://community.sphinx.chat/organizations/users/${orgUUID}/count`,
      expect.objectContaining({
        method: 'GET',
        mode: 'cors'
      })
    );

    expect(organizationsCount).toEqual({ count: 2 });
  });

  it('should call endpoint on deleteOrganizationUser', async () => {
    const mainStore = new MainStore();

    const mockApiResponse = {
      status: 200,
      json: jest.fn().mockResolvedValue({
        message: 'success'
      })
    };

    mockFetch.mockReturnValueOnce(Promise.resolve(mockApiResponse));

    const orgUserUUID = 'cldl1g04nncmf23du7kg';
    const deleteRequestBody = {
      org_uuid: 'cmas9gatu2rvqiev4ur0',
      user_created: '2024-01-03T22:07:39.504494Z',
      id: 263,
      uuid: 'cm3eulatu2rvqi9o75ug',
      owner_pubkey: '02af1ea854c7dc8634d08732d95c6057e6e08e01723da4f561d711a60aea708c00',
      owner_alias: 'Nayan',
      unique_name: 'nayan',
      description: 'description',
      tags: [],
      img: '',
      created: '2023-12-23T14:31:49.963009Z',
      updated: '2023-12-23T14:31:49.963009Z',
      unlisted: false,
      deleted: false,
      last_login: 1704289377,
      owner_route_hint:
        '03a6ea2d9ead2120b12bd66292bb4a302c756983dc45dcb2b364b461c66fd53bcb:1099519819777',
      owner_contact_key:
        'MIIBCgKCAQEAugvVYqgIIBmpLCjmaBhLi6GfxssrdM74diTlKpr+Qr/0Er1ND9YQ3HUveaI6V5DrBunulbSEZlIXIqVSLm2wobN4iAqvoGGx1aZ13ByOJLjINjD5nA9FnfAJpvcMV/gTDQzQL9NHojAeMx1WyAlhIILdiDm9zyCJeYj1ihC660xr6MyVjWn9brJv47P+Bq2x9AWPufYMMgPH7GV1S7KkjEPMbGCdUvUZLs8tzzKtNcABCHBQKOcBNG/D4HZcCREMP90zj8/NUzz9x92Z5zuvJ0/eZVF91XwyMtThrJ+AnrXWv7AEVy63mu9eAO3UYiUXq2ioayKBgalyos2Mcs9DswIDAQAB',
      price_to_meet: 0,
      new_ticket_time: 0,
      twitter_confirmed: false,
      extras: {},
      github_issues: {}
    };

    const deleteResponse = await mainStore.deleteOrganizationUser(deleteRequestBody, orgUserUUID);

    expect(mockFetch).toBeCalledWith(
      `https://community.sphinx.chat/organizations/users/${orgUserUUID}`,
      expect.objectContaining({
        method: 'DELETE',
        mode: 'cors',
        body: JSON.stringify(deleteRequestBody),
        headers: expect.objectContaining({
          'x-jwt': 'test_jwt',
          'Content-Type': 'application/json'
        })
      })
    );

    expect(deleteResponse.status).toBe(200);
  });
});
