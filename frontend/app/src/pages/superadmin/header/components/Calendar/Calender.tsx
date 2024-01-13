import React, { useState } from "react";
import styled from "styled-components";
import {
  add,
  differenceInDays,
  endOfMonth,
  setDate,
  startOfMonth,
  sub,
} from "date-fns";
import Cell from "./Cell";
import { MonthsDropdown, YearDropDown } from "./DropDown";

const weeks = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];

type Props = {
  value?: Date;
  onChange: (date: Date) => void;
};

const CalendarWrapper = styled.div`
  width: 370px;
  height: 300px;
  padding: 10px;
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
  width: 375px;
  justify-content: center;
  gap: 72.80px;
`;

const FlexiDiv = styled.div`
  gap: 12px;
  height: 40px;
  margin-top: 20px;
  display: flex;
  justify-content: center;
  align-items: baseline;
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
    if (previousYear >= today) {
      onChange(previousYear);
    }
  };

  const nextYear = () => {
    const nextYear = add(value, { years: 1 });
    if (nextYear >= today) {
      onChange(nextYear);
    }
  };

  const prevMonth = () => {
    const previousMonth = sub(value, { months: 1 });
    if (previousMonth >= today) {
      onChange(previousMonth);
    }
  };

  const nextMonth = () => {
    const nextMonth = add(value, { months: 1 });
    if (nextMonth >= today) {
      onChange(nextMonth);
    }
  };

  return (
    <MainContainer>
      <DropDownContainer>
        <FlexiDiv>
          <p onClick={prevYear}>{"<"}</p>
          <YearDropDown value={value} onYearChange={onChange} />
          <p onClick={nextYear}>{">"}</p>
        </FlexiDiv>
        <FlexiDiv>
          <p onClick={prevMonth}>{"<"}</p>
          <MonthsDropdown value={value} onMonthChange={onChange} />
          <p onClick={nextMonth}>{">"}</p>
        </FlexiDiv>
      </DropDownContainer>
      <CalendarWrapper>
        <Grid>
          {weeks.map((week: string) => (
            <Cell key={week} className="text-xs font-bold uppercase">
              {week}
            </Cell>
          ))}

          {Array.from({ length: prefixDays }).map((_: any, index: number) => (
            <Cell key={index} />
          ))}

          {Array.from({ length: numDays }).map((_: any, index: number) => {
            const date = index + 1;
            const isCurrentDate = date === value.getDate();

            return (
              <Cell
                key={date}
                isActive={isCurrentDate}
                onClick={() => handleClickDate(date)}
              >
                {date}
              </Cell>
            );
          })}

          {Array.from({ length: suffixDays }).map((_: any, index: number) => (
            <Cell key={index} />
          ))}
        </Grid>
      </CalendarWrapper>
    </MainContainer>
  );
};

export default Calendar;
