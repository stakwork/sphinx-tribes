import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import MobileView from './CodingBounty';

describe('MobileView component', () => {
  const defaultProps = {
    deliverables: 'Default Deliverables',
    description: 'Default Description',
    titleString: 'Default Title',
    nametag: <></>,
    labels: [],
    person: {
      owner_pubkey: 'DefaultOwnerPubKey',
      owner_route_hint: 'DefaultRouteHint',
      owner_alias: 'DefaultOwnerAlias',
    } as any,
    setIsPaidStatusPopOver: jest.fn(),
    creatorStep: 1,
    paid: false,
    tribe: 'Default Tribe',
    saving: 'false',
    isPaidStatusPopOver: false,
    isPaidStatusBadgeInfo: false,
    awardDetails: {
      name: 'Default Award',
      // Add other awardDetails properties as needed
    },
    isAssigned: false,
    dataValue: {},
    assigneeValue: false,
    assignedPerson: {
        owner_pubkey: 'DefaultOwnerPubKey',
        owner_route_hint: 'DefaultRouteHint',
        owner_alias: 'DefaultOwnerAlias',
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
    isMobile: true,
    actionButtons: false,
    assigneeLabel: {},
    setExtrasPropertyAndSave : jest.fn(),
    setIsModalSideButton : jest.fn(),
    setIsExtraStyle : jest.fn(),
    coding_languages:['language'],
    type :'',
     badgeRecipient: '',
      fromBountyPage:'',
       wanted_type:'',
        one_sentence_summary:'',
         github_description:'',
          show:false,
           formSubmit:jest.fn()
    
  };

  it('should render titleString on the screen', () => {
    render(<MobileView   ticket_url={''} assignee={defaultProps.person as any} title={''} {...defaultProps} titleString="Test Title" />);
    const titleElement = screen.getByText('Test Title');
    expect(titleElement).toBeInTheDocument();
  });

  it('should render description on the screen', () => {
    render(<MobileView   ticket_url={''} assignee={defaultProps.person as any} title={''} {...defaultProps} description="Test Description" />);
    const descriptionElement = screen.getByText('Test Description');
    expect(descriptionElement).toBeInTheDocument();
  });

  it('should render deliverables on the screen', () => {
    render(<MobileView   ticket_url={''} assignee={defaultProps.person as any} title={''} {...defaultProps} deliverables="Test Deliverables" />);
    const deliverablesElement = screen.getByText('Test Deliverables');
    expect(deliverablesElement).toBeInTheDocument();
  });

  // Add more test cases as needed for other functionality

  // Example test for button click
//   it('should call handleSetAsPaid when the "Set as Paid" button is clicked', () => {
//     const handleSetAsPaidMock = jest.fn();
//     render(<MobileView {...defaultProps} handleSetAsPaid={handleSetAsPaidMock} />);
//     const setAsPaidButton = screen.getByText('Set as Paid');
//     userEvent.click(setAsPaidButton);
//     expect(handleSetAsPaidMock).toHaveBeenCalled();
//   });

  // Add more test cases for other user interactions
});
