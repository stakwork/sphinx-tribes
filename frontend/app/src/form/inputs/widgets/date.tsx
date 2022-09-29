import moment from 'moment';
import { EuiDatePicker, EuiFormRow } from '@elastic/eui';
import React, { memo, useState } from 'react';
import { Props } from '../propsType';
import { FieldEnv } from '..';

function Date({ label, value, handleChange }: any) {
  const [startDate, setStartDate] = useState(moment(value) ?? moment());

  const handleChangeDate = (date) => {
    console.log(moment(date).toISOString());
    setStartDate(date);
    handleChange(date.toISOString());
  };

  return (
    <FieldEnv label={label}>
      <EuiDatePicker
        selectsEnd={true}
        selectsStart={true}
        selected={startDate}
        onChange={(e) => handleChangeDate(e)}
      />
    </FieldEnv>
  );
}
export default memo(Date);
