import React from 'react'
import { useObserver } from 'mobx-react-lite'
import { useStores } from '../store'

import {
  EuiHeader,
  EuiHeaderSection,
  EuiFieldSearch,
} from '@elastic/eui';

export default function Header() {
  const { main, ui } = useStores()

  return useObserver(() => {
    return <EuiHeader id="header" >
      <div className="container">
        <div className="row">
          <EuiHeaderSection grow={false} className="col-xs-12 col-sm-12 col-md-6 col-lg-6">
            <img id="logo" src="/static/people_logo.svg" alt="Logo"
            // style={{ cursor: 'pointer' }}
            // onClick={() => {
            //   window.history.pushState({}, 'Sphinx Tribes', '/');
            //   console.log('click!')
            // }}
            />
            {/*<Title>Tribes</Title>*/}
          </EuiHeaderSection>

          <EuiHeaderSection id="header-right" side="right" className="col-xs-12 col-sm-12 col-md-6 col-lg-6">
            {/* <EuiHeaderSectionItem> */}
            <div>
              <EuiFieldSearch id="search-input"
                placeholder="Search People"
                value={ui.searchText}
                onChange={e => ui.setSearchText(e.target.value)}
                // isClearable={this.state.isClearable}
                aria-label="search"
              />
            </div>
          </EuiHeaderSection>
        </div>
      </div>
    </EuiHeader>
  })
}
