import React, { useEffect } from 'react';

import { EuiHeader, EuiHeaderSection } from '@elastic/eui';
import { observer } from 'mobx-react-lite';
import { useStores } from '../../store';

function Header() {
  const { ui } = useStores();

  useEffect(() => {
    if (window.location.host === 'podcasts.sphinx.chat') {
      ui.setTags(
        ui.tags.map((t: any) => {
          if (t.label === 'Podcast') return { ...t, checked: 'on' };
          return t;
        })
      );
    }
  });

  return (
    <EuiHeader id="header">
      <div className="row" style={{ marginLeft: 15 }}>
        <EuiHeaderSection grow={false} className="col-xs-12 col-sm-12 col-md-6 col-lg-6">
          <img id="logo" src="/static/tribes_logo.svg" alt="Logo" />
        </EuiHeaderSection>
      </div>
    </EuiHeader>
  );
}
export default observer(Header);
