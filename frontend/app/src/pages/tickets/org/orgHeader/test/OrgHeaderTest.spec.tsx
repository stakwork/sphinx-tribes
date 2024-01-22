import '@testing-library/jest-dom';
import { render, screen, waitFor } from '@testing-library/react';
import nock from 'nock';
import React from 'react';
import { organization } from '__test__/__mockData__/organization';
import { mockUsehistory } from '__test__/__mockFn__/useHistory';
import OrgDescription from '../OrgDescription';

beforeAll(() => {
  nock.disableNetConnect();
  mockUsehistory();
});

const updateIsPostBountyModalOpen = jest.fn();

describe('OrgDescription Component', () => {
  it('renders the component with organization information', async () => {
    const url = 'http://localhost:5002';

    nock(url).get(`/organizations/${organization.uuid}`).reply(200, {});

    render(
      <OrgDescription
        updateIsPostBountyModalOpen={updateIsPostBountyModalOpen}
        orgData={organization}
      />
    );

    await waitFor(async () => {
      expect(await screen.findByText(organization.name)).toBeInTheDocument();
      expect(screen.getByText('Post a Bounty')).toBeInTheDocument();
      expect(screen.getByText('Website')).toBeInTheDocument();
      expect(screen.getByText('Github')).toBeInTheDocument();
    });
  });
});
