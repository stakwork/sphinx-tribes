import '@testing-library/jest-dom';
import React from 'react';
import { render, screen } from '@testing-library/react';
import { setupStore } from '__test__/__mockData__/setupStore';
import { mockUsehistory } from '__test__/__mockFn__/useHistory';
import nock from 'nock';
import { user } from '__test__/__mockData__/user';
import { person } from '__test__/__mockData__/persons';
import { Organization } from 'store/main';
import { mainStore } from '../../../store/main';
import OrganizationView from '../OrganizationView';

beforeAll(() => {
  nock.disableNetConnect();
  setupStore();
  mockUsehistory();
});

const organizations: Organization[] = [
  {
    bounty_count: 0,
    created: '2024-01-03T20:34:09.585609Z',
    deleted: false,
    id: '51',
    img: '',
    name: 'TEST_NEW',
    owner_pubkey: '03cbb9c01cdcf91a3ac3b543a556fbec9c4c3c2a6ed753e19f2706012a26367ae3',
    show: false,
    updated: '2024-01-03T20:34:09.585609Z',
    uuid: 'cmas9gatu2rvqiev4ur0'
  },
  {
    bounty_count: 0,
    created: '2024-01-03T20:34:09.585609Z',
    deleted: false,
    id: '52',
    img: '',
    name: 'TEST_SECOND',
    owner_pubkey: '03cbb9c01cdcf91a3ac3b543a556fbec9c4c3c2a6ed753e19f2706012a26367ae3',
    show: false,
    updated: '2024-01-03T20:34:09.585609Z',
    uuid: 'cmas9gatu2rvqiev4ur0'
  }
];

describe('OrganizationView Component', () => {
  nock(user.url).get('/person/id/1').reply(200, {});
  it('renders organization names correctly', async () => {
    jest.spyOn(mainStore, 'getUserRoles').mockReturnValue(Promise.resolve([]));
    jest.spyOn(mainStore, 'getOrganizationUser').mockReturnValue(Promise.resolve(person as any));
    mainStore.setOrganizations(organizations);

    render(<OrganizationView person={person} />);

    const organizationName = screen.getByText(organizations[0].name);
    const secondOrganization = screen.getByText(organizations[1].name);
    expect(organizationName).toBeInTheDocument();
    expect(secondOrganization).toBeInTheDocument();
  });

  it('renders view bounties button correctly', async () => {
    jest.spyOn(mainStore, 'getUserRoles').mockReturnValue(Promise.resolve([]));
    jest.spyOn(mainStore, 'getOrganizationUser').mockReturnValue(Promise.resolve(person as any));
    mainStore.setOrganizations([organizations[0]]);

    render(<OrganizationView person={person} />);

    const viewBountiesBtn = screen.getByRole('button', { name: 'View Bounties open_in_new_tab' });
    expect(viewBountiesBtn).toBeInTheDocument();
  });

  it('should not render manage bounties button if user does not have access', async () => {
    jest.spyOn(mainStore, 'getUserRoles').mockReturnValue(Promise.resolve([]));
    jest.spyOn(mainStore, 'getOrganizationUser').mockReturnValue(Promise.resolve({} as any));
    mainStore.setOrganizations([organizations[0]]);

    render(<OrganizationView person={person} />);

    const manageButton = screen.queryAllByRole('button', { name: 'Manage' });
    expect(manageButton.length).toBe(0);
  });

  it('renders manage bounties button if user is owner correctly', async () => {
    jest.spyOn(mainStore, 'getUserRoles').mockReturnValue(Promise.resolve([]));
    jest.spyOn(mainStore, 'getOrganizationUser').mockReturnValue(Promise.resolve(person as any));
    const userOrg = {
      ...organizations[0],
      owner_pubkey: person.owner_pubkey
    };
    mainStore.setOrganizations([userOrg]);

    render(<OrganizationView person={person} />);

    const manageButton = screen.getByRole('button', { name: 'Manage' });
    expect(manageButton).toBeInTheDocument();
  });
});
