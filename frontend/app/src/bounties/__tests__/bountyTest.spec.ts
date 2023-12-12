import { mainStore } from '../../store/main';
import api from '../../api';
import mockBounties, { newBounty, mockBountiesMutated } from '../__mock__/mockBounties.data';
import mockLocalStorage from '../__mock__/mockLocalStorage.utils';

jest.mock('../../api', () => ({
  __esModule: true,
  default: jest.fn()
}));

const mockedApi = api as jest.Mocked<typeof api>;
//const mockedUiStore = uiStore as jest.Mocked<typeof uiStore>;

describe('Bounty Tests', () => {
  beforeAll(() => {
    Object.defineProperty(window, 'localStorage', { value: mockLocalStorage });
  });

  beforeEach(() => {
    window.localStorage.clear();

    //window.localStorage.setItem('peopleBounties', JSON.stringify(mockBounties));
  });

  afterEach(() => {
    window.localStorage.clear();
  });

  //it('should save and retrieve bounties from local storage', () => {
  //  const storedBounties = JSON.parse(window.localStorage.getItem('peopleBounties') || '[]');
  //  expect(storedBounties).toEqual(mockBounties);
  //}, 8000);

  it('should fetch bounties and save to localStorage', async () => {
    mockedApi.get = jest.fn().mockResolvedValue(mockBounties);
    const bounties = await mainStore.getPeopleBounties({ resetPage: true });
    expect(mockedApi.get).toHaveBeenCalledTimes(1);
    expect(bounties).toStrictEqual(mockBountiesMutated);
  });

  it('should save a new bounty and persist to localStorage', async () => {
    //mockedApi.post.mockResolvedValue({ data: newBounty });
    global.fetch = jest.fn();
    mockedApi.post = jest.fn().mockResolvedValue(newBounty);
    await mainStore.saveBounty(newBounty);
    expect(mockedApi.post).toHaveBeenCalledTimes(0);
    expect(global.fetch).toHaveBeenCalledTimes(1);

    const postRequestContent = [
      'http://localhost:5002/gobounties?token=undefined',
      {
        body: '{"person":{"id":10,"name":"Alice Johnson"},"body":"Bounty details here","org_uuid":"org1","title":"Fix Backend Bug","description":"Fix a critical bug in the backend","owner_id":"owner123","created":1610000000,"show":true,"assignee":{"id":"dev123","name":"Jane Smith"},"wanted_type":"bugfix","type":"backend","price":"1000","codingLanguage":"JavaScript","estimated_session_length":"3 hours","bounty_expires":"2023-12-31","commitment_fee":50,"coding_languages":[]}',
        headers: { 'Content-Type': 'application/json', 'x-jwt': undefined },
        method: 'POST',
        mode: 'cors'
      }
    ];

    expect(global.fetch).toHaveBeenCalledWith(...postRequestContent);
    //const bounties = JSON.parse(mockLocalStorage.getItem('peopleBounties'));
    //expect(window.localStorage.getItem('peopleBounties')).toContainEqual([]);
  });

  it('should fetch and return a bounty matching from localStorage', async () => {
    const mockBounty = mockBounties[0];
    mockedApi.get = jest.fn().mockResolvedValue([mockBounty]);
    const bounty = await mainStore.getBountyById(mockBounty.bounty.id);
    const expectedBountyResponse = [
      {
        body: {
          assigned_hours: 0,
          assignee: '',
          award: '',
          bounty_expires: '',
          coding_languages: ['Lightning'],
          commitment_fee: 0,
          created: 1702398254,
          deliverables: 'test',
          description: 'test',
          estimated_completion_date: '',
          estimated_session_length: '',
          github_description: false,
          id: 892,
          one_sentence_summary: '',
          org_uuid: 'ck13rgua5fdkhph1po4g',
          owner_id: '03bfe6723c06fb2b7546df1e8ca1a17ae5c504615da32c945425ccbe8d3ca6260d',
          paid: false,
          paid_date: null,
          price: '1',
          show: true,
          ticket_url: '',
          title: 'test',
          tribe: 'None',
          type: 'freelance_job_request',
          updated: '2023-12-12T16:24:14.585187Z',
          wanted_type: 'Mobile development'
        },
        organization: {
          img: 'https://memes.sphinx.chat/public/l-_K9mJatGvz16Ixw1lPHtG9Om8QWtZtiRS_aIQme9c=',
          name: 'test after jwt',
          uuid: 'ck13rgua5fdkhph1po4g'
        },
        person: {
          created: '2021-11-15T19:09:46.356248Z',
          deleted: false,
          description: 'This is the real kevkevin I like to code alot',
          extras: null,
          github_issues: null,
          id: 38,
          img: 'https://memes.sphinx.chat/public/1jPjjDsrwRoBPzahpWjR8QE5DjQ726MsUSrCOjflSmo=',
          last_login: 1702065128,
          new_ticket_time: 0,
          owner_alias: 'kevkevin',
          owner_contact_key:
            'MIIBCgKCAQEAxtkwpxx8RjdVhgzx4oUYkmJQttvFzwI+lCWYgngMi/4o8OgUF9eVvW8zSY0t9A1KEY2MdEOTGjv9QiesoN7hmkgTdUqDQd1LIsU4vBtwPVWyJs0d6VEdMySN9veN68S7Fu+S20e5gygj17X8cffoEwLNDPi0dsTgojAC/uggE98zJvHmEd/Ob/W3ADQD68DQErCejvqXK2557GtsDNo35iIN9KlOPLRmvG3S/oV4pIyj5Z/6uMEXlok2b/mtvP0E4ClMP77j9QPs7mQarQ03XM0iRC2Ru/Qg/xWBTeqmYv5zfD8hmtzakBVyMSrHNZKZjSnURVNVpFaEXoB4wBqcvQIDAQAB',
          owner_pubkey: '',
          owner_route_hint: '',
          price_to_meet: 21,
          tags: [],
          twitter_confirmed: true,
          unique_name: 'umbreltest',
          unlisted: false,
          updated: '2023-12-04T19:48:05.641056Z',
          uuid: 'cd9dm5ua5fdtsj2c2nbg',
          wanteds: []
        }
      }
    ];
    expect(bounty).toEqual(expectedBountyResponse);
  });

  //it('should delete a bounty from localStorage', async () => {
  //  global.fetch = jest.fn();
  //				console.log(uiStore)
  //				const meInfoSpy = jest.spyOn(uiStore.uiStore, "meInfo", "get")

  //								meInfoSpy.mockReturnValue({
  //      id: 1,
  //      pubkey: '',
  //      uuid: '',
  //      owner_pubkey: '',
  //      owner_route_hint: '',
  //      photo_url: '',
  //      alias: '',
  //      img: '',
  //      owner_alias: '',
  //      github_issues: [],
  //      route_hint: '',
  //      contact_key: '',
  //      price_to_meet: 0,
  //      jwt: '',
  //      tribe_jwt: '',
  //      url: '',
  //      description: '',
  //      verification_signature: '',
  //      twitter_confirmed: false,
  //      extras: {},
  //      isSuperAdmin: false,
  //      websocketToken: ''
  //    })
  //  const bountyIdToDelete = mockBounties[0].bounty.id;
  //  const publicKeyToDelete = mockBounties[0].owner.owner_pubkey;
  //  mockedApi.del = jest.fn().mockResolvedValue({});
  //  await mainStore.deleteBounty(bountyIdToDelete, publicKeyToDelete);

  //  const deleteRequestContent = [
  //    'https://people-test.sphinx.chat/gobounties?token=undefined',
  //    {
  //      body: '{"person":{"id":10,"name":"Alice Johnson"},"body":"Bounty details here","org_uuid":"org1","title":"Fix Backend Bug","description":"Fix a critical bug in the backend","owner_id":"owner123","created":1610000000,"show":true,"assignee":{"id":"dev123","name":"Jane Smith"},"wanted_type":"bugfix","type":"backend","price":"1000","codingLanguage":"JavaScript","estimated_session_length":"3 hours","bounty_expires":"2023-12-31","commitment_fee":50,"coding_languages":[]}',
  //      headers: { 'Content-Type': 'application/json', 'x-jwt': undefined },
  //      method: 'POST',
  //      mode: 'cors'
  //    }
  //  ];

  //  expect(global.fetch).toHaveBeenCalledTimes(1);
  //  expect(global.fetch).toHaveBeenCalledWith(...deleteRequestContent);
  //  //const remainingBounties = JSON.parse(mockLocalStorage.getItem('peopleBounties'));
  //  //expect(remainingBounties.some((bounty: { id: number }) => bounty.id === bountyIdToDelete)).toBe(
  //  //  false
  //  //);
  //});
});
