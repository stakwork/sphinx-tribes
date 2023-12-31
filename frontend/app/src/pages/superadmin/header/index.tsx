import React, { useState } from 'react';
import moment from 'moment';
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
  7: 'Last 7 days',
  30: 'Last 30 days',
  45: 'Last 45 days'
};

interface HeaderProps {
  startDate?: number;
  endDate?: number;
  setStartDate: (newDate: number) => void;
  setEndDate: (newDate: number) => void;
}
export const Header = ({ startDate, setStartDate, endDate, setEndDate }: HeaderProps) => {
  const [dateDiff, setDateDiff] = useState(7);

  const formatUnixDate = (unixDate: number, includeYear: boolean = true) => {
    const formatString = includeYear ? 'DD-MMM-YYYY' : 'DD-MMM';
    return moment.unix(unixDate).format(formatString);
  };

  const handleBackClick = () => {
    if (startDate !== undefined && endDate !== undefined) {
      const newStartDate = moment.unix(startDate).subtract(dateDiff, 'days').unix();
      const newEndDate = moment.unix(endDate).subtract(dateDiff, 'days').unix();
      setStartDate(newStartDate);
      setEndDate(newEndDate);
    }
  };

  const handleForwardClick = () => {
    if (startDate !== undefined && endDate !== undefined) {
      const newStartDate = moment.unix(startDate).add(dateDiff, 'days').unix();
      const newEndDate = moment.unix(endDate).add(dateDiff, 'days').unix();

      // Ensure the end date does not go beyond today
      const todayUnix = moment().startOf('day').unix();
      const cappedEndDate = Math.min(newEndDate, todayUnix);

      setStartDate(newStartDate);
      setEndDate(cappedEndDate);
    }
  };

  const handleDropDownChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    const selectedValue = Number(event.target.value);
    setDateDiff(selectedValue);

    if (startDate && endDate) {
      const currentEndDate = moment.unix(endDate);
      let newStartDate;

      if (selectedValue === 7) {
        newStartDate = currentEndDate.clone().subtract(7, 'days').unix();
      } else if (selectedValue === 30) {
        newStartDate = currentEndDate.clone().subtract(30, 'days').unix();
      } else if (selectedValue === 45) {
        newStartDate = currentEndDate.clone().subtract(45, 'days').unix();
      }

      // Ensure the new start date does not go beyond the current end date
      newStartDate = Math.max(
        newStartDate,
        currentEndDate.clone().subtract(selectedValue, 'days').unix()
      );

      setStartDate(newStartDate);
    }
  };

  const currentDateUnix = moment().unix();
  console.log(startDate, endDate, 'date');

  return (
    <Container>
      <AlternateWrapper>
        <LeftWrapper>
          {startDate && endDate ? (
            <>
              <ButtonWrapper>
                <ArrowButton onClick={() => handleBackClick()}>
                  <img src={arrowback} alt="" />
                </ArrowButton>
                <ArrowButton
                  disabled={
                    endDate === currentDateUnix ||
                    moment.unix(endDate).isSameOrAfter(moment().startOf('day'))
                  }
                  onClick={() => handleForwardClick()}
                >
                  <img src={arrowforward} alt="" />
                </ArrowButton>
              </ButtonWrapper>
              <Month>
                {formatUnixDate(startDate, false)} - {formatUnixDate(endDate)}
              </Month>
            </>
          ) : null}
        </LeftWrapper>
        <RightWrapper>
          <ExportButton>
            <ExportText>Export CSV</ExportText>
          </ExportButton>
          <DropDown>
            <Select value={dateDiff} onChange={handleDropDownChange}>
              {Object.entries(DateFilterObject).map(([key, value]: any) => (
                <option key={key} value={key}>
                  {value}
                </option>
              ))}
            </Select>
          </DropDown>
        </RightWrapper>
      </AlternateWrapper>
    </Container>
  );
};
