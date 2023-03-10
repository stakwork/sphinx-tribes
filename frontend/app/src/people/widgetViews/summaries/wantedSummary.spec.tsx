
import React from 'react'
import { render, screen } from "@testing-library/react";
import OfferSummary from "./offerSummary";
import '@testing-library/jest-dom'
import WantedSummary from './wantedSummary';

describe('OfferSummary', () => {
  test('', () => {
    const props = {
      formSubmit: jest.fn(() => {}),
      ReCallBounties: jest.fn(() => {}),
      deleteAction: jest.fn(() => {}),
      editAction: jest.fn(() => {}),
    }
    render(<WantedSummary
      // {...item}
      ReCallBounties={props.ReCallBounties}
      formSubmit={props.formSubmit}
      deleteAction={props?.deleteAction}
      editAction={props?.editAction}
      // person={person}
      // personBody={props?.personBody}
      // fromBountyPage={props?.fromBountyPage}
      // extraModalFunction={props?.extraModalFunction}
      // deletingState={props?.deletingState}
      // setIsModalSideButton={props?.setIsModalSideButton}
      // setIsExtraStyle={props?.setIsExtraStyle}
           />);

    // const title = screen.getByText('title');
    // expect(title).toBeInTheDocument()
  })
});