import '@testing-library/jest-dom';
import { render, screen } from '@testing-library/react';
import { CodingBountiesProps } from 'people/interfaces';
import React from 'react';
import MobileView from '../CodingBounty';

describe('MobileView component', () => {
  beforeEach(() => {
    const mockIntersectionObserver = jest.fn();
    mockIntersectionObserver.mockReturnValue({
      observe: () => null,
      unobserve: () => null,
      disconnect: () => null
    });
    window.IntersectionObserver = mockIntersectionObserver;
  });

  const defaultProps: CodingBountiesProps = {
    deliverables: 'Default Deliverables',
    description: 'Default Description',
    titleString: 'Default Title',
    nametag: <></>,
    labels: [],
    person: {
      owner_pubkey: 'DefaultOwnerPubKey',
      owner_route_hint: 'DefaultRouteHint',
      owner_alias: 'DefaultOwnerAlias'
    } as any,
    setIsPaidStatusPopOver: jest.fn(),
    creatorStep: 1,
    paid: false,
    tribe: 'Default Tribe',
    saving: 'false',
    isPaidStatusPopOver: false,
    isPaidStatusBadgeInfo: false,
    awardDetails: {
      name: 'Default Award'
    },
    isAssigned: false,
    dataValue: { a: 'a' },
    assigneeValue: false,
    assignedPerson: {
      owner_pubkey: 'DefaultOwnerPubKey',
      owner_route_hint: 'DefaultRouteHint',
      owner_alias: 'DefaultOwnerAlias'
    } as any,
    changeAssignedPerson: jest.fn(),
    sendToRedirect: jest.fn(),
    handleCopyUrl: jest.fn(),
    isCopied: false,
    replitLink: 'DefaultReplitLink',
    assigneeHandlerOpen: jest.fn(),
    setCreatorStep: jest.fn(),
    awards: ['Award1', 'Award2'],
    setExtrasPropertyAndSaveMultiple: jest.fn(),
    handleAssigneeDetails: jest.fn(),
    peopleList: [],
    setIsPaidStatusBadgeInfo: jest.fn(),
    bountyPrice: 100,
    price: 100,
    selectedAward: 'DefaultSelectedAward',
    handleAwards: jest.fn(),
    repo: 'DefaultRepo',
    issue: 'DefaultIssue',
    isMarkPaidSaved: false,
    setAwardDetails: jest.fn(),
    setBountyPrice: jest.fn(),
    owner_idURL: 'DefaultOwnerIdURL',
    createdURL: 'DefaultCreatedURL',
    created: 1234567890,
    loomEmbedUrl: 'DefaultLoomEmbedUrl',
    org_uuid: 'DefaultOrgUUID',
    id: 987654321,
    localPaid: 'UNKNOWN' as any,
    setLocalPaid: jest.fn(),
    isMobile: false,
    actionButtons: false,
    assigneeLabel: {},
    setExtrasPropertyAndSave: jest.fn(),
    setIsModalSideButton: jest.fn(),
    setIsExtraStyle: jest.fn(),
    coding_languages: ['language'],
    type: '',
    badgeRecipient: '',
    fromBountyPage: '',
    wanted_type: '',
    one_sentence_summary: '',
    github_description: '',
    show: false,
    formSubmit: jest.fn(),
    ticket_url: '',
    assignee: undefined as any,
    title: ''
  };

  it('should render titleString on the screen', () => {
    render(<MobileView {...defaultProps} titleString="Test Title" />);
    const titleElement = screen.getByText('Test Title');
    expect(titleElement).toBeInTheDocument();
  });

  it('should render description on the screen', () => {
    render(<MobileView {...defaultProps} description="Test Description" />);
    const descriptionElement = screen.getByText('Test Description');
    expect(descriptionElement).toBeInTheDocument();
  });

  it('should render deliverables on the screen', () => {
    render(<MobileView {...defaultProps} deliverables="Test Deliverables" />);
    const deliverablesElement = screen.getByText('Test Deliverables');
    expect(deliverablesElement).toBeInTheDocument();
  });

  it('I can help button is rendered on the screen', () => {
    render(<MobileView {...defaultProps} />);

    const iCanHelp = screen.getByText('I can help');
    expect(iCanHelp).toBeInTheDocument();
  });

  it('correct price is rendered on the screen', () => {
    render(<MobileView {...defaultProps} />);

    const bountyPrice = screen.getByText(`${defaultProps.bountyPrice}`);
    expect(bountyPrice).toBeInTheDocument();
  });

  it('share with twitter button is rendered on the screen', () => {
    render(<MobileView {...defaultProps} />);

    const iCanHelp = screen.getByText('Share to Twitter');
    expect(iCanHelp).toBeInTheDocument();
  });
});
