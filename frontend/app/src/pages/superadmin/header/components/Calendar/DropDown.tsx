import { setMonth, setYear } from 'date-fns';
import React, { ChangeEvent, useState } from 'react';
import styled from 'styled-components';

type YearDropDownProps = {
  onYearChange: (date: Date) => void;
  value: Date;
};

const DropDown = styled.select`
  outline: none;
  border: none;
`;
const List = styled.div`
  display: flex;
  padding: 10px 4px 10px 8px;
  justify-content: center;
  align-items: center;
  gap: 8px;
`;
const YearDropDown = ({ value, onYearChange }: YearDropDownProps) => {
  const currentYear = new Date().getFullYear();
  const years = Array.from({ length: 3 }, (_: number, index: number) => currentYear - index);
  const [selectedYear, setSelectedYear] = useState<string>(value.getFullYear().toString());

  const changeYear = (event: ChangeEvent<HTMLSelectElement>) => {
    const selectedYear = event.target.value;
    const year: Date = setYear(value, parseInt(selectedYear, 10));
    setSelectedYear(selectedYear);
    onYearChange(year);
  };

  return (
    <List>
      <DropDown id="year" value={selectedYear} onChange={changeYear}>
        <option value="" disabled hidden>
          Year
        </option>
        {years.map((year: number) => (
          <option key={year} value={year.toString()}>
            {year}
          </option>
        ))}
      </DropDown>
    </List>
  );
};

type MonthsDropdownProps = {
  value: Date;
  onMonthChange: (date: Date) => void;
};

const MonthsDropdown = ({ value, onMonthChange }: MonthsDropdownProps) => {
  const months = [
    'January',
    'February',
    'March',
    'April',
    'May',
    'June',
    'July',
    'August',
    'September',
    'October',
    'November',
    'December'
  ];
  const [selectedMonth, setSelectedMonth] = useState<string>((value.getMonth() + 1).toString());

  const changeMonth = (event: ChangeEvent<HTMLSelectElement>) => {
    const selectedMonth = event.target.value;
    const month = setMonth(value, parseInt(selectedMonth, 10) - 1);
    setSelectedMonth(selectedMonth);
    onMonthChange(month);
  };

  return (
    <List>
      <DropDown id="month" value={selectedMonth} onChange={changeMonth}>
        <option value="" disabled hidden>
          Month
        </option>
        {months.map((month: String, index: number) => (
          <option key={index} value={(index + 1).toString()}>
            {month}
          </option>
        ))}
      </DropDown>
    </List>
  );
};

export { YearDropDown, MonthsDropdown };
