import { toJS } from 'mobx';
import sinon from 'sinon';
import moment from 'moment';
import { people } from '../../__test__/__mockData__/persons';
import { user } from '../../__test__/__mockData__/user';
import { MeInfo, emptyMeInfo, uiStore } from '../ui';
import { MainStore } from '../main';
import { localStorageMock } from '../../__test__/__mockData__/localStorage';
import { TribesURL, getHost } from '../../config';
import mockBounties, { expectedBountyResponses } from '../../bounties/__mock__/mockBounties.data';

let fetchStub: sinon.SinonStub;
let mockApiResponseData: any[];

const origFetch = global.fetch;

beforeAll(() => {
  fetchStub = sinon.stub(global, 'fetch');
  fetchStub.returns(Promise.resolve({ status: 200, json: () => Promise.resolve({}) })); // Mock a default behavior
  mockApiResponseData = [
    { uuid: 'cm3eulatu2rvqi9o75ug' },
    { uuid: 'cldl1g04nncmf23du7kg' },
    { orgUUID: 'cmas9gatu2rvqiev4ur0' }
  ];
});

afterAll(() => {
  global.fetch = origFetch;

  sinon.restore();
});

describe('Main store', () => {
  beforeEach(async () => {
    uiStore.setMeInfo(user);
    localStorageMock.setItem('ui', JSON.stringify(uiStore));
  });

  afterEach(() => {
    fetchStub.reset();
  });

  it('should call endpoint on addOrganization', async () => {
    const mainStore = new MainStore();

    const mockApiResponse = { status: 200, message: 'success' };

    fetchStub.resolves(Promise.resolve(mockApiResponse));

    const addOrganization = {
      img: '',
      name: 'New Orgination test',
      owner_pubkey: '035f22835fbf55cf4e6823447c63df74012d1d587ed60ef7cbfa3e430278c44cce'
    };

    const expectedHeaders = {
      'Content-Type': 'application/json',
      'x-jwt': 'test_jwt'
    };

    await mainStore.addOrganization(addOrganization);

    sinon.assert.calledWith(
      fetchStub,
      `${TribesURL}/organizations`,
      sinon.match({
        method: 'POST',
        headers: expectedHeaders,
        body: JSON.stringify(addOrganization),
        mode: 'cors'
      })
    );
  });

  it('should call endpoint on addOrganization with description, github and website url', async () => {
    const mainStore = new MainStore();

    const mockApiResponse = { status: 200, message: 'success' };

    fetchStub.resolves(Promise.resolve(mockApiResponse));

    const addOrganization = {
      img: '',
      name: 'New Orgination test',
      owner_pubkey: '035f22835fbf55cf4e6823447c63df74012d1d587ed60ef7cbfa3e430278c44cce',
      description: 'My test Organization',
      github: 'https://github.com/john-doe',
      website: 'https://john.doe'
    };

    const expectedHeaders = {
      'Content-Type': 'application/json',
      'x-jwt': 'test_jwt'
    };

    await mainStore.addOrganization(addOrganization);

    sinon.assert.calledWith(
      fetchStub,
      `${TribesURL}/organizations`,
      sinon.match({
        method: 'POST',
        headers: expectedHeaders,
        body: JSON.stringify(addOrganization),
        mode: 'cors'
      })
    );
  });

  it('should call endpoint on UpdateOrganization Name', async () => {
    const mainStore = new MainStore();

    const mockApiResponse = { status: 200, message: 'success' };

    fetchStub.resolves(Promise.resolve(mockApiResponse));

    const updateOrganization = {
      id: '42',
      uuid: 'clic8k04nncuuf32kgr0',
      name: 'TEST1',
      owner_pubkey: '035f22835fbf55cf4e6823447c63df74012d1d587ed60ef7cbfa3e430278c44cce',
      img: 'https://memes.sphinx.chat/public/NVhwFqDqHKAC-_Sy9pR4RNy8_cgYuOVWgohgceAs-aM=',
      created: '2023-11-27T16:31:12.699355Z',
      updated: '2023-11-27T16:31:12.699355Z',
      show: false,
      deleted: false,
      bounty_count: 1
    };

    const expectedHeaders = {
      'Content-Type': 'application/json',
      'x-jwt': 'test_jwt'
    };

    await mainStore.updateOrganization(updateOrganization);

    sinon.assert.calledWith(
      fetchStub,
      `${TribesURL}/organizations`,
      sinon.match({
        method: 'POST',
        headers: expectedHeaders,
        body: JSON.stringify(updateOrganization),
        mode: 'cors'
      })
    );
  });

  it('should call endpoint on UpdateOrganization description, github url and website url, non mandatory fields', async () => {
    const mainStore = new MainStore();

    const mockApiResponse = { status: 200, message: 'success' };

    fetchStub.resolves(Promise.resolve(mockApiResponse));

    const updateOrganization = {
      id: '42',
      uuid: 'clic8k04nncuuf32kgr0',
      name: 'TEST1',
      owner_pubkey: '035f22835fbf55cf4e6823447c63df74012d1d587ed60ef7cbfa3e430278c44cce',
      img: 'https://memes.sphinx.chat/public/NVhwFqDqHKAC-_Sy9pR4RNy8_cgYuOVWgohgceAs-aM=',
      created: '2023-11-27T16:31:12.699355Z',
      updated: '2023-11-27T16:31:12.699355Z',
      show: false,
      deleted: false,
      bounty_count: 1,
      description: 'Update description',
      website: 'https://john.doe',
      github: 'https://github.com/john-doe'
    };

    const expectedHeaders = {
      'Content-Type': 'application/json',
      'x-jwt': 'test_jwt'
    };

    await mainStore.updateOrganization(updateOrganization);

    sinon.assert.calledWith(
      fetchStub,
      `${TribesURL}/organizations`,
      sinon.match({
        method: 'POST',
        headers: expectedHeaders,
        body: JSON.stringify(updateOrganization),
        mode: 'cors'
      })
    );
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

    fetchStub.resolves(Promise.resolve(mockApiResponse));

    const organizationUser = {
      owner_pubkey: user.owner_pubkey || '',
      org_uuid: mockApiResponseData[2]
    };

    const expectedHeaders = {
      'Content-Type': 'application/json',
      'x-jwt': 'test_jwt'
    };

    await mainStore.addOrganizationUser(organizationUser);

    sinon.assert.calledWith(
      fetchStub,
      `${TribesURL}/organizations/users/${mockApiResponseData[2]}`,
      sinon.match({
        method: 'POST',
        headers: expectedHeaders,
        body: JSON.stringify(organizationUser),
        mode: 'cors'
      })
    );
  });

  it('should call endpoint on getOrganizationUsers', async () => {
    const mainStore = new MainStore();

    const mockApiResponse = {
      status: 200,
      json: sinon.stub().resolves(mockApiResponseData.slice(0, 1))
    };

    fetchStub.resolves(Promise.resolve(mockApiResponse));

    const endpoint = `${TribesURL}/organizations/users/${mockApiResponseData[2].orgUUID}`;

    const users = await mainStore.getOrganizationUsers(mockApiResponseData[2].orgUUID);

    sinon.assert.calledWithMatch(fetchStub, endpoint, sinon.match.any);
    expect(users).toEqual(mockApiResponseData.slice(0, 1));
  });

  it('should call endpoint on getUserOrganizations', async () => {
    const mainStore = new MainStore();
    const userId = 232;
    const mockOrganizations = [
      {
        id: 42,
        uuid: 'clic8k04nncuuf32kgr0',
        name: 'TEST',
        owner_pubkey: '035f22835fbf55cf4e6823447c63df74012d1d587ed60ef7cbfa3e430278c44cce',
        img: 'https://memes.sphinx.chat/public/NVhwFqDqHKAC-_Sy9pR4RNy8_cgYuOVWgohgceAs-aM=',
        created: '2023-11-27T16:31:12.699355Z',
        updated: '2023-11-27T16:31:12.699355Z',
        show: false,
        deleted: false,
        bounty_count: 1
      },
      {
        id: 55,
        uuid: 'cmen35itu2rvqicrm020',
        name: 'Orgination name test',
        owner_pubkey: '035f22835fbf55cf4e6823447c63df74012d1d587ed60ef7cbfa3e430278c44cce',
        img: '',
        created: '2024-01-09T16:17:26.202555Z',
        updated: '2024-01-09T16:17:26.202555Z',
        show: false,
        deleted: false
      },
      {
        id: 56,
        uuid: 'cmen38itu2rvqicrm02g',
        name: 'New Orgination test',
        owner_pubkey: '035f22835fbf55cf4e6823447c63df74012d1d587ed60ef7cbfa3e430278c44cce',
        img: '',
        created: '2024-01-09T16:17:38.652072Z',
        updated: '2024-01-09T16:17:38.652072Z',
        show: false,
        deleted: false
      },
      {
        id: 49,
        uuid: 'cm7c24itu2rvqi9o7620',
        name: 'TESTing',
        owner_pubkey: '02af1ea854c7dc8634d08732d95c6057e6e08e01723da4f561d711a60aea708c00',
        img: '',
        created: '2023-12-29T12:52:34.62057Z',
        updated: '2023-12-29T12:52:34.62057Z',
        show: false,
        deleted: false
      },
      {
        id: 51,
        uuid: 'cmas9gatu2rvqiev4ur0',
        name: 'TEST_NEW',
        owner_pubkey: '03cbb9c01cdcf91a3ac3b543a556fbec9c4c3c2a6ed753e19f2706012a26367ae3',
        img: '',
        created: '2024-01-03T20:34:09.585609Z',
        updated: '2024-01-03T20:34:09.585609Z',
        show: false,
        deleted: false
      }
    ];
    const mockApiResponse = {
      status: 200,
      json: sinon.stub().resolves(mockOrganizations)
    };
    fetchStub.resolves(Promise.resolve(mockApiResponse));

    const organizationUser = await mainStore.getUserOrganizations(userId);

    sinon.assert.calledWithMatch(
      fetchStub,
      `${TribesURL}/organizations/user/${userId}`,
      sinon.match({
        method: 'GET',
        mode: 'cors',
        headers: {
          'Content-Type': 'application/json'
        }
      })
    );

    expect(organizationUser).toEqual(mockOrganizations);
  });

  it('should call endpoint on getUserOrganizationsUuid', async () => {
    const mainStore = new MainStore();
    const uuid = 'ck1p7l6a5fdlqdgmmnpg';
    const mockOrganizations = {
      id: 6,
      uuid: 'ck1p7l6a5fdlqdgmmnpg',
      name: 'Stakwork',
      owner_pubkey: '021ae436bcd40ca21396e59be8cdb5a707ceacdb35c1d2c5f23be7584cab29c40b',
      img: 'https://memes.sphinx.chat/public/_IO8M0UXltb3mbK0qso63ux86AP-2nN2Ly9uHo37Ku4=',
      created: '2023-09-14T23:14:28.821632Z',
      updated: '2023-09-14T23:14:28.821632Z',
      show: true,
      deleted: false,
      bounty_count: 8,
      budget: 640060
    };
    const mockApiResponse = {
      status: 200,
      json: sinon.stub().resolves(mockOrganizations)
    };
    fetchStub.resolves(Promise.resolve(mockApiResponse));

    const organizationUser = await mainStore.getUserOrganizationByUuid(uuid);

    sinon.assert.calledWithMatch(
      fetchStub,
      `${TribesURL}/organizations/${uuid}`,
      sinon.match({
        method: 'GET',
        mode: 'cors',
        headers: {
          'Content-Type': 'application/json'
        }
      })
    );

    expect(organizationUser).toEqual(mockOrganizations);
  });
  it('should call endpoint on getOrganizationUser', async () => {
    const mainStore = new MainStore();

    const mockApiResponse = {
      status: 200,
      json: sinon.stub().resolves({
        uuid: mockApiResponseData[0].uuid
      })
    };

    fetchStub.resolves(Promise.resolve(mockApiResponse));

    const organizationUser = await mainStore.getOrganizationUser(mockApiResponseData[0].uuid);

    sinon.assert.calledWithMatch(
      fetchStub,
      `${TribesURL}/organizations/foruser/${mockApiResponseData[0].uuid}`,
      sinon.match({
        method: 'GET',
        mode: 'cors',
        headers: {
          'x-jwt': 'test_jwt',
          'Content-Type': 'application/json'
        }
      })
    );

    expect(organizationUser).toEqual({
      uuid: mockApiResponseData[0].uuid
    });
  });

  it('should call endpoint on getOrganizationUsersCount', async () => {
    const mainStore = new MainStore();

    const mockApiResponse = {
      status: 200,
      json: sinon.stub().resolves({
        count: 2
      })
    };

    fetchStub.resolves(Promise.resolve(mockApiResponse));

    const organizationsCount = await mainStore.getOrganizationUsersCount(
      mockApiResponseData[2].orgUUID
    );

    sinon.assert.calledWithMatch(
      fetchStub,
      `${TribesURL}/organizations/users/${mockApiResponseData[2].orgUUID}/count`,
      sinon.match({
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
      json: sinon.stub().resolves({
        message: 'success'
      })
    };

    fetchStub.resolves(Promise.resolve(mockApiResponse));

    const orgUserUUID = mockApiResponseData[1].uuid;
    const deleteRequestBody = {
      org_uuid: mockApiResponseData[2].orgUUID,
      user_created: '2024-01-03T22:07:39.504494Z',
      id: 263,
      uuid: mockApiResponseData[0].uuid,
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

    sinon.assert.calledWithMatch(
      fetchStub,
      `${TribesURL}/organizations/users/${orgUserUUID}`,
      sinon.match({
        method: 'DELETE',
        mode: 'cors',
        body: JSON.stringify(deleteRequestBody),
        headers: sinon.match({
          'x-jwt': 'test_jwt',
          'Content-Type': 'application/json'
        })
      })
    );

    expect(deleteResponse.status).toBe(200);
  });

  it('should send request delete request with correct body and url', async () => {
    const url = `${TribesURL}/gobounties/pub_key/1111`;
    const allBountiesUrl = `http://${getHost()}/gobounties/all?limit=10&sortBy=created&search=&page=1&resetPage=true&Open=true&Assigned=false&Paid=false`;
    const expectedRequestOptions: RequestInit = {
      method: 'DELETE',
      mode: 'cors',
      headers: {
        'x-jwt': user.tribe_jwt,
        'Content-Type': 'application/json'
      }
    };
    fetchStub.withArgs(url, expectedRequestOptions).returns(
      Promise.resolve({
        status: 200
      }) as any
    );
    fetchStub.withArgs(allBountiesUrl, sinon.match.any).returns(
      Promise.resolve({
        status: 200,
        ok: true,
        json: (): Promise<any> => Promise.resolve([mockBounties[0]])
      }) as any
    );

    const store = new MainStore();
    await store.deleteBounty(1111, 'pub_key');

    expect(fetchStub.withArgs(url, expectedRequestOptions).calledOnce).toEqual(true);
    expect(store.peopleBounties.length).toEqual(0);
    expect(store.peopleBounties).toEqual([]);
  });

  it('should not panic if failed to delete bounty', async () => {
    const url = `${TribesURL}/gobounties/pub_key/1111`;
    const expectedRequestOptions: RequestInit = {
      method: 'DELETE',
      mode: 'cors',
      headers: {
        'x-jwt': user.tribe_jwt,
        'Content-Type': 'application/json'
      }
    };
    fetchStub.withArgs(url, expectedRequestOptions).throwsException();

    const store = new MainStore();
    await store.deleteBounty(1111, 'pub_key');

    expect(fetchStub.withArgs(url, expectedRequestOptions).calledOnce).toEqual(true);
    expect(store.peopleBounties.length).toEqual(0);
  });

  it('should not return false if asignee removed successfully', async () => {
    const url = `${TribesURL}/gobounties/assignee`;
    const expectedRequestOptions: RequestInit = {
      method: 'DELETE',
      mode: 'cors',
      body: JSON.stringify({
        owner_pubkey: 'pub_key',
        created: '1111'
      }),
      headers: {
        'x-jwt': user.tribe_jwt,
        'Content-Type': 'application/json'
      }
    };
    fetchStub.withArgs(url, expectedRequestOptions).returns(
      Promise.resolve({
        status: 200
      }) as any
    );

    const store = new MainStore();
    const res = await store.deleteBountyAssignee({ owner_pubkey: 'pub_key', created: '1111' });

    expect(fetchStub.withArgs(url, expectedRequestOptions).calledOnce).toEqual(true);
    expect(res).not.toBeFalsy();
  });

  it('should  return false if failed to remove asignee ', async () => {
    const url = `${TribesURL}/gobounties/assignee`;
    const expectedRequestOptions: RequestInit = {
      method: 'DELETE',
      mode: 'cors',
      body: JSON.stringify({
        owner_pubkey: 'pub_key',
        created: '1111'
      }),
      headers: {
        'x-jwt': user.tribe_jwt,
        'Content-Type': 'application/json'
      }
    };
    fetchStub.withArgs(url, expectedRequestOptions).throwsException();

    const store = new MainStore();
    const res = await store.deleteBountyAssignee({ owner_pubkey: 'pub_key', created: '1111' });

    expect(fetchStub.withArgs(url, expectedRequestOptions).calledOnce).toEqual(true);
    expect(res).toBeFalsy();
  });

  it('should successfully update bounty payment status', async () => {
    const url = `${TribesURL}/gobounties/paymentstatus/1111`;
    const expectedRequestOptions: RequestInit = {
      method: 'POST',
      mode: 'cors',
      headers: {
        'x-jwt': user.tribe_jwt,
        'Content-Type': 'application/json'
      }
    };
    fetchStub.withArgs(url, expectedRequestOptions).returns(
      Promise.resolve({
        status: 200
      }) as any
    );

    const store = new MainStore();
    const res = await store.updateBountyPaymentStatus(1111);

    expect(fetchStub.withArgs(url, expectedRequestOptions).calledOnce).toEqual(true);
    expect(res).not.toBeFalsy();
  });

  it('should return false if failed to update bounty status', async () => {
    const url = `${TribesURL}/gobounties/paymentstatus/1111`;
    const expectedRequestOptions: RequestInit = {
      method: 'POST',
      mode: 'cors',
      headers: {
        'x-jwt': user.tribe_jwt,
        'Content-Type': 'application/json'
      }
    };
    fetchStub.withArgs(url, expectedRequestOptions).throwsException();

    const store = new MainStore();
    const res = await store.updateBountyPaymentStatus(1111);

    expect(fetchStub.withArgs(url, expectedRequestOptions).calledOnce).toEqual(true);
    expect(res).toBeFalsy();
  });

  it('should successfully return requested bounty', async () => {
    const url = `http://${getHost()}/gobounties/id/1111`;
    fetchStub.withArgs(url, sinon.match.any).returns(
      Promise.resolve({
        status: 200,
        ok: true,
        json: () => Promise.resolve([mockBounties[0]])
      }) as any
    );

    const store = new MainStore();
    const res = await store.getBountyById(1111);

    expect(fetchStub.withArgs(url, sinon.match.any).calledOnce).toEqual(true);
    expect(res).toEqual([expectedBountyResponses[0]]);
  });

  it('should return empty array if failed to fetch bounty', async () => {
    const url = `http://${getHost()}/gobounties/id/1111`;
    fetchStub.withArgs(url, sinon.match.any).returns(
      Promise.resolve({
        status: 404,
        ok: false
      }) as any
    );

    const store = new MainStore();
    const res = await store.getBountyById(1111);

    expect(fetchStub.withArgs(url, sinon.match.any).calledOnce).toEqual(true);
    expect(res.length).toEqual(0);
  });

  it('should successfully return index of requested bounty', async () => {
    const url = `http://${getHost()}/gobounties/index/1111`;
    fetchStub.withArgs(url, sinon.match.any).returns(
      Promise.resolve({
        status: 200,
        ok: true,
        json: () => Promise.resolve(1)
      }) as any
    );

    const store = new MainStore();
    const res = await store.getBountyIndexById(1111);

    expect(fetchStub.withArgs(url, sinon.match.any).calledOnce).toEqual(true);
    expect(res).toEqual(1);
  });

  it('should return 0 if failed to fetch index', async () => {
    const url = `http://${getHost()}/gobounties/index/1111`;
    fetchStub.withArgs(url, sinon.match.any).returns(
      Promise.resolve({
        status: 400,
        ok: false
      }) as any
    );

    const store = new MainStore();
    const res = await store.getBountyIndexById(1111);

    expect(fetchStub.withArgs(url, sinon.match.any).calledOnce).toEqual(true);
    expect(res).toEqual(0);
  });

  it('should set all query params, page, limit, search, and languages when fetching bounties, user logged out', async () => {
    uiStore.setMeInfo(emptyMeInfo);
    const allBountiesUrl = `http://${getHost()}/gobounties/all?limit=10&sortBy=updatedat&search=random&page=1&resetPage=true`;
    fetchStub.withArgs(allBountiesUrl, sinon.match.any).returns(
      Promise.resolve({
        status: 200,
        ok: true,
        json: (): Promise<any> => Promise.resolve([mockBounties[0]])
      }) as any
    );

    const store = new MainStore();
    const bounties = await store.getPeopleBounties({
      resetPage: true,
      search: 'random',
      limit: 11,
      page: 1,
      sortBy: 'updatedat'
    });

    expect(store.peopleBounties.length).toEqual(1);
    expect(store.peopleBounties).toEqual([expectedBountyResponses[0]]);
    expect(bounties).toEqual([expectedBountyResponses[0]]);
  });

  it('should reset exisiting bounty if reset flag is passed, signed out', async () => {
    uiStore.setMeInfo(emptyMeInfo);
    const allBountiesUrl = `http://${getHost()}/gobounties/all?limit=10&sortBy=updatedat&search=random&page=2&resetPage=true`;
    const mockBounty = { ...mockBounties[0] };
    mockBounty.bounty.id = 2;
    fetchStub.withArgs(allBountiesUrl, sinon.match.any).returns(
      Promise.resolve({
        status: 200,
        ok: true,
        json: (): Promise<any> => Promise.resolve([{ ...mockBounty }])
      }) as any
    );

    const store = new MainStore();
    store.setPeopleBounties([expectedBountyResponses[0] as any]);
    expect(store.peopleBounties.length).toEqual(1);

    const bounties = await store.getPeopleBounties({
      resetPage: true,
      search: 'random',
      limit: 11,
      page: 2,
      sortBy: 'updatedat'
    });
    const expectedResponse = { ...expectedBountyResponses[0] };
    expectedResponse.body.id = 2;
    expect(store.peopleBounties.length).toEqual(1);
    expect(store.peopleBounties).toEqual([expectedResponse]);
    expect(bounties).toEqual([expectedResponse]);
  });

  it('should add to exisiting bounty if next page is fetched, user signed out', async () => {
    uiStore.setMeInfo(emptyMeInfo);
    const allBountiesUrl = `http://${getHost()}/gobounties/all?limit=10&sortBy=updatedat&search=random&page=2&resetPage=false`;
    const mockBounty = { ...mockBounties[0] };
    mockBounty.bounty.id = 2;
    fetchStub.withArgs(allBountiesUrl, sinon.match.any).returns(
      Promise.resolve({
        status: 200,
        ok: true,
        json: (): Promise<any> => Promise.resolve([{ ...mockBounty }])
      }) as any
    );

    const store = new MainStore();
    const bountyAlreadyPresent = { ...expectedBountyResponses[0] } as any;
    bountyAlreadyPresent.body.id = 1;
    store.setPeopleBounties([bountyAlreadyPresent]);
    expect(store.peopleBounties.length).toEqual(1);
    expect(store.peopleBounties[0].body.id).not.toEqual(2);

    const bounties = await store.getPeopleBounties({
      resetPage: false,
      search: 'random',
      limit: 11,
      page: 2,
      sortBy: 'updatedat'
    });

    const expectedResponse = { ...expectedBountyResponses[0] };
    expectedResponse.body.id = 2;
    expect(store.peopleBounties.length).toEqual(2);
    expect(store.peopleBounties[1]).toEqual(expectedResponse);
    expect(bounties).toEqual([expectedResponse]);
  });

  it('should successfully fetch people, user signed out', async () => {
    uiStore.setMeInfo(emptyMeInfo);
    const allBountiesUrl = `http://${getHost()}/people?resetPage=true&search=&limit=500&page=1&sortBy=last_login`;
    const mockPeople = { ...people[1] };
    fetchStub.withArgs(allBountiesUrl, sinon.match.any).returns(
      Promise.resolve({
        status: 200,
        ok: true,
        json: (): Promise<any> => Promise.resolve([{ ...mockPeople }])
      }) as any
    );

    const store = new MainStore();
    store.setPeople([people[0]]);
    expect(store._people.length).toEqual(1);
    expect(store._people[0]).toEqual(people[0]);

    const res = await store.getPeople({
      resetPage: true,
      search: 'random',
      limit: 11,
      page: 1,
      sortBy: 'updatedat'
    });

    expect(store._people.length).toEqual(1);
    expect(store._people[0]).toEqual(mockPeople);
    expect(res[0]).toEqual(mockPeople);
  });

  it('should hide current user, user signed in', async () => {
    const allBountiesUrl = `http://${getHost()}/people?resetPage=false&search=&limit=500&page=2&sortBy=last_login`;
    const mockPeople = { ...people[0] };
    fetchStub.withArgs(allBountiesUrl, sinon.match.any).returns(
      Promise.resolve({
        status: 200,
        ok: true,
        json: (): Promise<any> => Promise.resolve([{ ...mockPeople }])
      }) as any
    );

    const store = new MainStore();
    const res = await store.getPeople({
      resetPage: false,
      search: 'random',
      limit: 11,
      page: 2,
      sortBy: 'updatedat'
    });

    expect(store._people.length).toEqual(1);
    expect(store._people[0].hide).toEqual(true);
    expect(res).toBeTruthy();
  });

  it('should fetch and store organization bounties successfully, user signed out', async () => {
    uiStore.setMeInfo(emptyMeInfo);
    const allBountiesUrl = `http://${getHost()}/organizations/bounties/1111`;
    fetchStub.withArgs(allBountiesUrl, sinon.match.any).returns(
      Promise.resolve({
        status: 200,
        ok: true,
        json: (): Promise<any> => Promise.resolve([mockBounties[0]])
      }) as any
    );

    const store = new MainStore();
    const bounties = await store.getOrganizationBounties('1111', {
      page: 1,
      resetPage: true
    });

    expect(store.peopleBounties.length).toEqual(1);
    expect(store.peopleBounties).toEqual([expectedBountyResponses[0]]);
    expect(bounties).toEqual([expectedBountyResponses[0]]);
  });

  it('should reset exisiting organization bounty if reset flag is passed, user signed out', async () => {
    uiStore.setMeInfo(emptyMeInfo);
    const allBountiesUrl = `http://${getHost()}/organizations/bounties/1111`;
    const mockBounty = { ...mockBounties[0] };
    mockBounty.bounty.id = 2;
    fetchStub.withArgs(allBountiesUrl, sinon.match.any).returns(
      Promise.resolve({
        status: 200,
        ok: true,
        json: (): Promise<any> => Promise.resolve([{ ...mockBounty }])
      }) as any
    );

    const store = new MainStore();
    store.setPeopleBounties([expectedBountyResponses[0] as any]);
    expect(store.peopleBounties.length).toEqual(1);

    const bounties = await store.getOrganizationBounties('1111', {
      page: 1,
      resetPage: true
    });
    const expectedResponse = { ...expectedBountyResponses[0] };
    expectedResponse.body.id = 2;
    expect(store.peopleBounties.length).toEqual(1);
    expect(store.peopleBounties).toEqual([expectedResponse]);
    expect(bounties).toEqual([expectedResponse]);
  });

  it('should add to exisiting bounty if reset flag is not passed, user signed out', async () => {
    uiStore.setMeInfo(emptyMeInfo);
    const allBountiesUrl = `http://${getHost()}/organizations/bounties/1111`;
    const mockBounty = { ...mockBounties[0] };
    mockBounty.bounty.id = 2;
    fetchStub.withArgs(allBountiesUrl, sinon.match.any).returns(
      Promise.resolve({
        status: 200,
        ok: true,
        json: (): Promise<any> => Promise.resolve([{ ...mockBounty }])
      }) as any
    );

    const store = new MainStore();
    const bountyAlreadyPresent = { ...expectedBountyResponses[0] } as any;
    bountyAlreadyPresent.body.id = 1;
    store.setPeopleBounties([bountyAlreadyPresent]);
    expect(store.peopleBounties.length).toEqual(1);
    expect(store.peopleBounties[0].body.id).not.toEqual(2);

    const bounties = await store.getOrganizationBounties('1111', {
      page: 1,
      resetPage: false
    });

    const expectedResponse = { ...expectedBountyResponses[0] };
    expectedResponse.body.id = 2;
    expect(store.peopleBounties.length).toEqual(2);
    expect(store.peopleBounties[1]).toEqual(expectedResponse);
    expect(bounties).toEqual([expectedResponse]);
  });

  it('should make a succcessful bounty payment', async () => {
    const store = new MainStore();
    uiStore.setMeInfo(emptyMeInfo);
    const bounty = expectedBountyResponses[0];

    store.makeBountyPayment = jest
      .fn()
      .mockReturnValueOnce(Promise.resolve({ status: 200, message: 'success' }));

    const body = {
      id: bounty.body.id,
      websocket_token: 'test_websocket_token'
    };

    store.makeBountyPayment(body);
    expect(store.makeBountyPayment).toBeCalledWith(body);
  });

  it('it should get a s3 URL afer a successful metrics url call', async () => {
    const store = new MainStore();
    uiStore.setMeInfo(emptyMeInfo);

    store.exportMetricsBountiesCsv = jest
      .fn()
      .mockReturnValueOnce(
        Promise.resolve({ status: 200, body: 'https://test-s3url.com/metrics.csv' })
      );

    const start_date = moment().subtract(30, 'days').unix().toString();
    const end_date = moment().unix().toString();

    const body = {
      start_date,
      end_date
    };

    store.exportMetricsBountiesCsv(body);
    expect(store.exportMetricsBountiesCsv).toBeCalledWith(body);
  });

  it('I should be able to test that the signed-in user details are persisted in the local storage', async () => {
    const mockUser: MeInfo = {
      id: 20,
      pubkey: 'test_pub_key',
      uuid: mockApiResponseData[0].uuid,
      contact_key: 'test_owner_contact_key',
      owner_route_hint: 'test_owner_route_hint',
      alias: 'Owner Name',
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
      owner_alias: 'Owner Name',
      owner_pubkey: 'test_pub_key',
      img: '/static/avatarPlaceholders/placeholder_34.jpg',
      twitter_confirmed: false,
      isSuperAdmin: false,
      websocketToken: 'test_websocketToken'
    };
    uiStore.setMeInfo(mockUser);
    uiStore.setShowSignIn(true);

    localStorageMock.setItem('ui', JSON.stringify(uiStore));

    expect(uiStore.showSignIn).toBeTruthy();
    expect(uiStore.meInfo).toEqual(mockUser);
    expect(localStorageMock.getItem('ui')).toEqual(JSON.stringify(uiStore));
  });

  it('I should be able to test that when signed out the user data is deleted', async () => {
    // Shows first if signed in
    uiStore.setShowSignIn(true);
    localStorageMock.setItem('ui', JSON.stringify(uiStore));

    expect(uiStore.showSignIn).toBeTruthy();
    expect(localStorageMock.getItem('ui')).toEqual(JSON.stringify(uiStore));
    //Shows when signed out
    uiStore.setMeInfo(null);
    localStorageMock.setItem('ui', null);

    expect(localStorageMock.getItem('ui')).toEqual(null);
  });

  it('I should be able to test that signed-in user details can be displayed such as the name and pubkey', async () => {
    uiStore.setShowSignIn(true);

    expect(uiStore.meInfo?.owner_alias).toEqual(user.alias);
    expect(uiStore.meInfo?.owner_pubkey).toEqual(user.pubkey);
  });

  it('I should be able to test that a signed-in user can update their details', async () => {
    uiStore.setShowSignIn(true);
    expect(uiStore.meInfo?.alias).toEqual('Vladimir');

    user.alias = 'John';
    uiStore.setMeInfo(user);

    expect(uiStore.meInfo?.alias).toEqual('John');
  });

  it('I should be able to test that a signed-in user can make an API request without getting a 401 (unauthorized error)', async () => {
    uiStore.setShowSignIn(true);
    const loggedUrl = `http://${getHost()}/admin/auth`;
    const res = await fetchStub.withArgs(loggedUrl, sinon.match.any).returns(
      Promise.resolve({
        status: 200,
        ok: true
      }) as any
    );
    expect(res).toBeTruthy();
  });

  it('I should be able to test that when a user is signed out, a user will get a 401 error if they make an API call', async () => {
    uiStore.setMeInfo(emptyMeInfo);
    const urlNoLogged = `http://${getHost()}/admin/auth`;

    const res = await fetchStub.withArgs(urlNoLogged, sinon.match.any).returns(
      Promise.resolve({
        status: 401,
        ok: false
      }) as any
    );
    expect(res).toBeTruthy();
  });
});
