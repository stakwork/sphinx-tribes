import React from 'react';
import { render, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import '@testing-library/jest-dom/extend-expect';
import { mainStore } from '../../../../../store/main.ts';
import { OrgHeader } from '../index.tsx';
import { OrgBountyHeaderProps } from '../../../../../people/interfaces.ts';

const org_uuid = 'clf6qmo4nncmf23du7ng';
const MockProps: OrgBountyHeaderProps = {
  checkboxIdToSelectedMap: {
    Opened: false,
    Assigned: false,
    Paid: false,
    Completed: false
  },
  languageString: '',
  org_uuid,
  onChangeStatus: jest.fn()
};

describe('OrgHeader component', () => {
  beforeEach(() => {
    jest.spyOn(mainStore, 'getPeopleBounties').mockReset();
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  test('should trigger API call in response to click on status from OrgHeader', async () => {
    const { getByText } = render(<OrgHeader {...MockProps} />);

    const statusFilter = getByText('Status');
    expect(statusFilter).toBeInTheDocument();

    fireEvent.click(statusFilter);

    await waitFor(() => {
      expect(mainStore.getPeopleBounties).toHaveBeenCalledWith({
        page: 1,
        resetPage: true,
        ...MockProps.checkboxIdToSelectedMap,
        languages: MockProps.languageString,
        org_uuid
      });
    });
  });
});
