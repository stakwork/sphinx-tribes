import React, { useState } from 'react';
import styled from 'styled-components';
import { add, differenceInDays, endOfMonth, setDate, startOfMonth, sub } from 'date-fns';
import Cell from './Cell';
import { MonthsDropdown, YearDropDown } from './DropDown';

const weeks = ['S', 'M', 'T', 'W', 'T', 'F', 'S'];

type Props = {
  value?: Date;
  onChange: (date: Date) => void;
};

const CalendarWrapper = styled.div`
  width: 370px;
  height: 300px;
`;

const Grid = styled.div`
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 0px 12px 4px 12px;
`;

const DropDownContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  align-self: stretch;
  border-top: 1px solid var(--Divider-2, #dde1e5);
  border-bottom: 1px solid var(--Divider-2, #dde1e5);
  width: 100%;
  justify-content: center;
  height: 64px;
  gap: 20px;
`;

const FlexiDiv = styled.div`
  height: 48px;
  display: flex;
  justify-content: center;
  align-items: center;
  width: 172px;
`;
const ArrowBtn = styled.div`
  display: flex;
  padding: 8px;
  justify-content: center;
  align-items: center;
  gap: 10px;
`;

const MainContainer = styled.div``;

const Calendar = ({ value = new Date(), onChange }: Props) => {
  const startDate = startOfMonth(value);
  const endDate = endOfMonth(value);
  const numDays = differenceInDays(endDate, startDate) + 1;

  const prefixDays = startDate.getDay();
  const suffixDays = 6 - endDate.getDay();

  const today = new Date();
  const disableDateSelection = (date: Date) => date > today;

  const handleClickDate = (index: number) => {
    const date = setDate(value, index);
    if (!disableDateSelection(date)) {
      onChange(date);
    }
  };

  const prevYear = () => {
    const previousYear = sub(value, { years: 1 });
    if (previousYear.getFullYear() <= today.getFullYear()) {
      onChange(previousYear);
    }
  };

  const nextYear = () => {
    const nextYear = add(value, { years: 1 });
    if (nextYear.getFullYear() <= today.getFullYear()) {
      onChange(nextYear);
    }
  };

  const prevMonth = () => {
    const previousMonth = sub(value, { months: 1 });
    if (
      previousMonth.getFullYear() >= today.getFullYear() ||
      previousMonth.getMonth() >= today.getMonth()
    ) {
      onChange(previousMonth);
    }
  };

  const nextMonth = () => {
    const nextMonth = add(value, { months: 1 });
    if (
      nextMonth.getFullYear() <= today.getFullYear() ||
      nextMonth.getMonth() <= today.getMonth()
    ) {
      onChange(nextMonth);
    }
  };

  return (
    <MainContainer>
      <DropDownContainer>
        <FlexiDiv>
          <ArrowBtn onClick={prevYear}>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="25"
              viewBox="0 0 24 25"
              fill="none"
            >
              <path
                d="M15.705 7.91L14.295 6.5L8.29498 12.5L14.295 18.5L15.705 17.09L11.125 12.5L15.705 7.91Z"
                fill="#3C3F41"
              />
            </svg>
          </ArrowBtn>
          <YearDropDown value={value} onYearChange={onChange} />
          <ArrowBtn onClick={nextYear}>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="25"
              viewBox="0 0 24 25"
              fill="none"
            >
              <path
                d="M9.70498 6.5L8.29498 7.91L12.875 12.5L8.29498 17.09L9.70498 18.5L15.705 12.5L9.70498 6.5Z"
                fill="#3C3F41"
              />
            </svg>
          </ArrowBtn>
        </FlexiDiv>
        <FlexiDiv>
          <ArrowBtn onClick={prevMonth}>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="25"
              viewBox="0 0 24 25"
              fill="none"
            >
              <path
                d="M15.705 7.91L14.295 6.5L8.29498 12.5L14.295 18.5L15.705 17.09L11.125 12.5L15.705 7.91Z"
                fill="#3C3F41"
              />
            </svg>
          </ArrowBtn>
          <MonthsDropdown value={value} onMonthChange={onChange} />
          <ArrowBtn onClick={nextMonth}>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="25"
              viewBox="0 0 24 25"
              fill="none"
            >
              <path
                d="M9.70498 6.5L8.29498 7.91L12.875 12.5L8.29498 17.09L9.70498 18.5L15.705 12.5L9.70498 6.5Z"
                fill="#3C3F41"
              />
            </svg>
          </ArrowBtn>
        </FlexiDiv>
      </DropDownContainer>
      <CalendarWrapper>
        <Grid>
          {weeks.map((week: string) => (
            <Cell key={week} isEmpty={true}>
              {week}
            </Cell>
          ))}

          {Array.from({ length: prefixDays }).map((_: any, index: number) => (
            <Cell key={index} isEmpty={true} />
          ))}
          {Array.from({ length: numDays }).map((_: any, index: number) => {
            const date = index + 1;
            const isCurrentDate = date === value.getDate();

            return (
              <Cell
                key={date}
                isActive={isCurrentDate}
                onClick={() => handleClickDate(date)}
                isEmpty={false}
              >
                {date}
              </Cell>
            );
          })}

          {Array.from({ length: suffixDays }).map((_: any, index: number) => (
            <Cell key={index} isEmpty={true} />
          ))}
        </Grid>
      </CalendarWrapper>
    </MainContainer>
  );
};

export default Calendar;
