import React from 'react';
import PageLoadSpinner from '..//utils/pageLoadSpinner';

export const LoaderBottom = (props: any) => (
  <PageLoadSpinner
    noAnimate
    show={props.loadingBottom}
    style={{ position: 'absolute', bottom: 0, left: 0 }}
  />
);

export const LoaderTop = (props: any) => <PageLoadSpinner show={props.loadingTop} />;
