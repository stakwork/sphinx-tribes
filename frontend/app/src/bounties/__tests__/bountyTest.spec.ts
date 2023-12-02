import '@testing-library/jest-dom'
import axios from 'axios';
import { mainStore } from 'store/main';
import mockBounties, { newBounty } from '../__mock__/mockBounties.data';
import mockLocalStorage from '../__mock__/mockLocalStorage.utils';


jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;


describe('Bounty Tests', () => {
  beforeAll(() => {
    Object.defineProperty(window, 'localStorage', { value: mockLocalStorage });
  });

  beforeEach(() => {
    window.localStorage.clear();
    window.localStorage.setItem('peopleBounties', JSON.stringify(mockBounties));
  });

  afterEach(() => {
    window.localStorage.clear();
  });

  it('should save and retrieve bounties from local storage', () => {
    const storedBounties = JSON.parse(window.localStorage.getItem('peopleBounties') || '[]');
    expect(storedBounties).toEqual(mockBounties);
  });

it('should fetch bounties and save to localStorage', async () => {
    mockedAxios.get.mockResolvedValue({ data: mockBounties });
    await mainStore.getPeopleBounties()
    expect(localStorage.getItem('peopleBounties')).toEqual(JSON.stringify(mockBounties));
  });

  it('should save a new bounty and persist to localStorage', async () => {
   mockedAxios.post.mockResolvedValue({ data: newBounty });
    await mainStore.saveBounty(newBounty);
    const bounties = JSON.parse(mockLocalStorage.getItem('peopleBounties'));
    expect(bounties).toContainEqual(newBounty);
  });

  it('should fetch and return a bounty matching from localStorage', async () => {
    const mockBounty = mockBounties[0];
    mockedAxios.get.mockResolvedValue({ data: mockBounty });
    const bounty = await mainStore.getBountyById(mockBounty.person.id);
    expect(bounty).toEqual(mockBounty);
  });

  it('should delete a bounty from localStorage', async () => {
    const bountyIdToDelete = mockBounties[0].person.id;
    const publicKeyToDelete = mockBounties[0].owner_id
    mockedAxios.delete.mockResolvedValue({});
    await mainStore.deleteBounty(bountyIdToDelete,publicKeyToDelete);
    const remainingBounties = JSON.parse(mockLocalStorage.getItem('peopleBounties'));
    expect(remainingBounties.some((bounty: { id: number}) => bounty.id === bountyIdToDelete)).toBe(false);
  });
});
