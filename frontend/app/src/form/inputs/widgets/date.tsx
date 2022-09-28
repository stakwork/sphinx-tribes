import moment from 'moment';
import { EuiDatePicker, EuiFormRow } from '@elastic/eui';
import React, { useState } from 'react';
import { Props } from '../propsType';

export default function Date({
  label,
  note,
  value,
  name,
  handleChange,
  handleBlur,
  handleFocus,
  readOnly,
  prepend,
  extraHTML
}: Props) {
  const [startDate, setStartDate] = useState(moment());

  const handleChangeDate = (date) => {
    setStartDate(date);
    handleChange(startDate);
  };

  return (
    <EuiFormRow label={label}>
      <div>
        <EuiDatePicker selected={startDate} onChange={handleChangeDate} />
      </div>
    </EuiFormRow>
  );
}
