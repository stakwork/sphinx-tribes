import React from 'react';
import { mount } from 'enzyme';
import { TicketModalPage } from './TicketModalPage';

describe('prevArrHandler', () => {
  const initialState = {
    main: {
      peopleBounties: [{ body: { id: 1 }, person: {} }],
      getBountyById: jest.fn(),
      getBountyIndexById: jest.fn(),
      getOrganizationBounties: jest.fn()
    },
    modals: { setStartupModal: jest.fn() },
    ui: { meInfo: true },
    response: [{ body: { id: 1 }, person: {} }],
    bountyId: '1',
    connectPersonBody: {},
    publicFocusIndex: 0,
    removeNextAndPrev: false,
    visible: true,
    isDeleted: false,
    search: { owner_id: '1', created: '1' },
    history: { goBack: jest.fn() },
    getUuidFromUrl: jest.fn(() => '1'),
    directionHandler: jest.fn(),
    setResponse: jest.fn(),
    uuid: '1'
  };

  const component = mount(<TicketModalPage {...initialState} setConnectPerson={jest.fn()} />);

  it('should call directionHandler with the correct arguments when called', () => {
    component.instance().prevArrHandler();
    expect(component.instance().directionHandler).toHaveBeenCalledWith({}, { id: 0 });
  });

  it('should not call directionHandler when index is less than or equal to 0', () => {
    component.instance().publicFocusIndex = -1;
    component.instance().prevArrHandler();
    expect(component.instance().directionHandler).not.toHaveBeenCalled();
  });

  it('should not call directionHandler when there are no more items to navigate to', () => {
    component.instance().publicFocusIndex = 1;
    component.instance().prevArrHandler();
    expect(component.instance().directionHandler).not.toHaveBeenCalled();
  });
});
