import { mainStore } from '../../store/main';
import api from '../../api';
import mockBounties, {
  newBounty,
  mockBountiesMutated,
  expectedBountyResponses
} from '../__mock__/mockBounties.data';
import { localStorageMock } from '../../__test__/__mockData__/localStorage';

jest.mock('../../api', () => ({
  __esModule: true,
  default: jest.fn()
}));

const mockedApi = api as jest.Mocked<typeof api>;

describe('Bounty Tests', () => {
  beforeAll(() => {
    Object.defineProperty(window, 'localStorage', { value: localStorageMock });
  });

  beforeEach(() => {
    window.localStorage.clear();
  });

  afterEach(() => {
    window.localStorage.clear();
  });

  it('should fetch bounties and save to localStorage', async () => {
    mockedApi.get = jest.fn().mockResolvedValue(mockBounties);
    const bounties = await mainStore.getPeopleBounties({ resetPage: true });
    expect(mockedApi.get).toHaveBeenCalledTimes(1);
    expect(bounties).toStrictEqual(mockBountiesMutated);
  });

  it('should save a new bounty and persist to localStorage', async () => {
    global.fetch = jest.fn();
    mockedApi.post = jest.fn().mockResolvedValue(newBounty);
    await mainStore.saveBounty(newBounty);
    expect(mockedApi.post).toHaveBeenCalledTimes(0);
    expect(global.fetch).toHaveBeenCalledTimes(1);

    const postRequestContent = [
      'http://localhost:5002/gobounties?token=undefined',
      {
        body: JSON.stringify(newBounty),
        headers: { 'Content-Type': 'application/json', 'x-jwt': undefined },
        method: 'POST',
        mode: 'cors'
      }
    ];

    expect(global.fetch).toHaveBeenCalledWith(...postRequestContent);
  });

  it('should fetch and return a bounty matching from localStorage', async () => {
    const mockBounty = mockBounties[0];
    mockedApi.get = jest.fn().mockResolvedValue([mockBounty]);
    const bounty = await mainStore.getBountyById(mockBounty.bounty.id);
    const expectedBountyResponse = expectedBountyResponses;
    expect(bounty).toEqual(expectedBountyResponse);
  });

  it('should delete a bounty from localStorage', async () => {
    const bountyIdToDelete = mockBounties[0].bounty.id;
    const publicKeyToDelete = mockBounties[0].owner.owner_pubkey;
    mockedApi.del = jest.fn().mockResolvedValue({});
    await mainStore.deleteBounty(bountyIdToDelete, publicKeyToDelete);

    const deleteRequestContent = [
      'http://localhost:5002/gobounties?token=undefined',
      {
        body: JSON.stringify(newBounty),
        headers: { 'Content-Type': 'application/json', 'x-jwt': undefined },
        method: 'POST',
        mode: 'cors'
      }
    ];

    expect(global.fetch).toHaveBeenCalledTimes(1);
    expect(global.fetch).toHaveBeenCalledWith(...deleteRequestContent);
    const rawDeletedBounty = localStorageMock.getItem(`bounty_${bountyIdToDelete}`);
    const deletedBounty = rawDeletedBounty ? JSON.parse(rawDeletedBounty) : null;

    expect(deletedBounty).toBeNull();
  });

  it('should fetch and persist people bounties to localStorage', async () => {
    await mainStore.getPeopleBounties();

    const peopleRequestContent = [
      'http://localhost:5002/gobounties?token=undefined',
      {
        body: JSON.stringify(newBounty),
        headers: { 'Content-Type': 'application/json', 'x-jwt': undefined },
        method: 'POST',
        mode: 'cors'
      }
    ];

    const rawStoredBounties = localStorageMock.getItem('peopleBounties');
    const storedBounties = rawStoredBounties ? JSON.parse(rawStoredBounties) : null;
    expect(storedBounties).toBeDefined();
    expect(global.fetch).toHaveBeenCalledTimes(1);
    expect(global.fetch).toHaveBeenCalledWith(...peopleRequestContent);
  });
});
