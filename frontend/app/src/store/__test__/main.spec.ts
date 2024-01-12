import { toJS } from 'mobx';
import sinon from 'sinon';
import { people } from '../../__test__/__mockData__/persons';
import { user } from '../../__test__/__mockData__/user';
import { emptyMeInfo, uiStore } from '../ui';
import { MainStore } from '../main';
import { localStorageMock } from '../../__test__/__mockData__/localStorage';
import { TribesURL, getHost } from '../../config';
import mockBounties, { expectedBountyResponses } from '../../bounties/__mock__/mockBounties.data';
import moment from 'moment';

let fetchStub: sinon.SinonStub;

beforeAll(() => {
  fetchStub = sinon.stub(global, 'fetch');
});

afterAll(() => {
  jest.clearAllMocks();
});

describe('Main store', () => {
  beforeEach(async () => {
    uiStore.setMeInfo(user);
    localStorageMock.setItem('ui', JSON.stringify(uiStore));
  });

  afterEach(() => {
    fetchStub.reset();
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
    expect(store.peopleBounties.length).toEqual(1);
    expect(store.peopleBounties).toEqual([expectedBountyResponses[0]]);
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

  it('should set all query params, page, limit, search when fetching bounties, user logged out', async () => {
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
      resetPage: true,
      search: 'random',
      limit: 11,
      page: 2,
      sortBy: 'updatedat'
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
});
