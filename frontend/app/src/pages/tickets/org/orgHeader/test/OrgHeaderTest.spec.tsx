import '@testing-library/jest-dom';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import nock from 'nock';
import React from 'react';
import { organization } from '__test__/__mockData__/organization';
import OrgDescription from '../OrgDescription';

const updateIsPostBountyModalOpen = jest.fn();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useParams: () => ({
    uuid: 'cmg6oqitu2rnslkcjbqg',
    id: '57'
  })
}));

describe('OrgDescription Component', () => {
  it('renders the component with organization information', async () => {
    const url = 'http://localhost:5002';

    nock(url).get(`/organizations/${organization.uuid}`).reply(200, {});

    render(<OrgDescription updateIsPostBountyModalOpen={updateIsPostBountyModalOpen} />);

    await (async () => {
      expect(screen.findByText(organization.name)).toBeInTheDocument();
      expect(screen.getByText('Post a Bounty')).toBeInTheDocument();
      expect(screen.getByText('Website')).toBeInTheDocument();
      expect(screen.getByText('Github')).toBeInTheDocument();
      fireEvent.click(screen.getByText('Post a Bounty'));
    });
  });
});
