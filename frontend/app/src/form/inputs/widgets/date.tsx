import moment from 'moment';
import { EuiDatePicker, EuiFormRow } from '@elastic/eui';
import React, { memo, useState } from 'react';
import { Props } from '../propsType';
import { FieldEnv } from '..';
import styled from 'styled-components';

function Date({ label, value, handleChange }: any) {
  const [startDate, setStartDate] = useState(moment(value) ?? moment());
  const [isBorder, setIsBorder] = useState<boolean>(false);

  const handleChangeDate = (date) => {
    console.log(moment(date).toISOString());
    setStartDate(date);
    handleChange(date.toISOString());
  };

  return (
    <FieldEnv label={label} isTop={true}>
      <DataPicker
        selectsEnd={true}
        selectsStart={true}
        selected={startDate}
        onChange={(e) => {
          handleChangeDate(e);
        }}
        onFocus={() => {
          setIsBorder(true);
        }}
        border={isBorder}
      />
    </FieldEnv>
  );
}
export default memo(Date);

interface datePickerProps {
  border?: boolean;
}

const DataPicker = styled(EuiDatePicker)<datePickerProps>`
  border: 1px solid ${(p) => (p.border ? '#82B4FF' : '#DDE1E5')};
  :focus {
    background-image: none;
  }
`;
