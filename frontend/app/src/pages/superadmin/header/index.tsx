import React from 'react';
import {
  AlternateWrapper,
  ButtonWrapper,
  ExportButton,
  ExportText,
  Month,
  ArrowButton,
  DropDown,
  LeftWrapper,
  Select,
  RightWrapper,
  Container
} from './HeaderStyles';
import arrowback from './icons/arrowback.svg';
import arrowforward from './icons/arrowforward.svg';
//import './Header.css';

const DateFilterObject = {
  7: '7 days',
  30: '30 days',
  45: '45 days'
};

export const Header = () => (
  <Container>
    <AlternateWrapper>
      <LeftWrapper>
        <ButtonWrapper>
          <ArrowButton>
            <img src={arrowback} alt="" />
          </ArrowButton>
          <ArrowButton>
            <img src={arrowforward} alt="" />
          </ArrowButton>
        </ButtonWrapper>
        <Month>01 Oct - 31 Dec 2023</Month>
      </LeftWrapper>
      <RightWrapper>
        <ExportButton>
          <ExportText>Export CSV</ExportText>
        </ExportButton>
        <DropDown>
          <Select name="" id="">
            {Object.keys(DateFilterObject).map((key: any) => (
              <option key={key} value={DateFilterObject[key]}>
                {DateFilterObject[key]}
              </option>
            ))}
          </Select>
        </DropDown>
      </RightWrapper>
    </AlternateWrapper>
  </Container>
);
