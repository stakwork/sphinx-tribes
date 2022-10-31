import React, { useEffect } from 'react';
import { useObserver } from 'mobx-react-lite';
import { useStores } from '../store';

import { EuiHeader, EuiHeaderSection } from '@elastic/eui';

export default function Header() {
  const { ui } = useStores();

  useEffect(() => {
    if (window.location.host === 'podcasts.sphinx.chat') {
      ui.setTags(
        ui.tags.map((t) => {
          if (t.label === 'Podcast') return { ...t, checked: 'on' };
          return t;
        })
      );
    }
  });

  return useObserver(() => {
    return (
      <EuiHeader id="header">
        <div className="row" style={{ marginLeft: 15 }}>
          <EuiHeaderSection grow={false} className="col-xs-12 col-sm-12 col-md-6 col-lg-6">
            <img id="logo" src="/static/tribes_logo.svg" alt="Logo" />
          </EuiHeaderSection>
        </div>
      </EuiHeader>
    );
  });
}
