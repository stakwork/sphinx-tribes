import { SvgMaskProps } from 'people/interfaces';
import React from 'react';

export const SvgMask = (props: SvgMaskProps) => (
  <div
    {...props}
    style={{
      ...props?.svgStyle,
      width: props?.width ?? '24px',
      height: props?.height ?? '24px',
      WebkitMask: `url('${props?.src}') center center no-repeat`,
      backgroundColor: props?.bgcolor,
      maskSize: props?.size
    }}
  />
);
