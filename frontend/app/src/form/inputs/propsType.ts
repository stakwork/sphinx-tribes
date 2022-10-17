import React from 'react';

export interface Props {
  value: any;
  label: string;
  type?: string;
  handleChange: any;
  handleBlur: any;
  handleFocus: any;
  readOnly: boolean;
  prepend?: string;
  extraHTML?: string;
  note?: string;
  options?: any[];
  name: string;
  error: string;
  borderType?: 'bottom' | 'outline';
  imageIcon?: boolean;

  notProfilePic?: boolean;
}
