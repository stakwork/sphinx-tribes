import { useIsMobile, usePerson } from 'hooks';
import { observer } from 'mobx-react-lite';
import RenderWidgets from 'people/widgetViews/renderWidgets';
import { widgetConfigs } from 'people/utils/constants';
import React, { useCallback, useEffect } from 'react';
import { Route, Switch, useHistory, useLocation, useRouteMatch } from 'react-router-dom';
import { useStores } from 'store';
import styled from 'styled-components';
import { Wanted } from './Wanted';

const tabs = widgetConfigs;
export const TabsPages = observer(() => {
  const location = useLocation();
  const { ui } = useStores();
  const { url, path } = useRouteMatch();
  const history = useHistory();
  const isMobile = useIsMobile();
  const personId = ui.selectedPerson;
  const { person, canEdit } = usePerson(personId);

  const tabsNames = Object.keys(tabs).filter((name: any) => {
    if (name === 'about' && !isMobile) {
      return false;
    }
    return true;
  });
  const changeTabHandler = useCallback(
    (tabName: any) => {
      history.replace({
        pathname: `${url}/${tabName}`
      });
    },
    [history, url]
  );

  useEffect(() => {
    const tabSelected = tabsNames.some((name: any) => location.pathname.includes(name));
    if (!tabSelected) {
      changeTabHandler(tabsNames[0]);
    }
  }, [changeTabHandler, location.pathname, tabsNames]);

  const fullSelectedWidget = (name: any) => person?.extras?.[name];

  return (
    <Container isMobile={isMobile}>
      <Tabs
        style={{
          background: '#fff',
          padding: '0 20px',
          borderBottom: 'solid 1px #ebedef',
          boxShadow: canEdit ? '0px 2px 0px rgba(0, 0, 0, 0.07)' : '0px 2px 6px rgba(0, 0, 0, 0.07)'
        }}
      >
        {tabs &&
          tabsNames.map((name: any, i: number) => {
            const t = tabs[name];
            const { label } = t;
            const selected = location.pathname.includes(name);
            const hasExtras = !!person?.extras?.[name]?.length;
            const count: any = hasExtras
              ? person.extras[name].filter((f: any) => {
                  if ('show' in f) {
                    // show has a value
                    if (!f.show) return false;
                  }
                  // if no value default to true
                  return true;
                }).length
              : null;

            return (
              <Tab
                key={i}
                style={{ height: 64, alignItems: 'center' }}
                selected={selected}
                onClick={() => {
                  changeTabHandler(name);
                }}
              >
                {label}
                {count > 0 && <Counter>{count}</Counter>}
              </Tab>
            );
          })}
      </Tabs>
      <Switch>
        {tabsNames.map((name: any) => (
          <Route key={name} path={`${path}${name}`}>
            <div
              style={{
                padding: 20,
                height: 'calc(100% - 63px)',
                overflowY: 'auto',
                position: 'relative'
              }}
            >
              <div
                style={{
                  display: 'flex',
                  alignItems: 'flex-start',
                  justifyContent:
                    fullSelectedWidget && fullSelectedWidget.length > 0 ? 'flex-start' : 'center',
                  flexWrap: 'wrap',
                  minHeight: '100%'
                }}
              >
                {name === 'wanted' && <Wanted />}
                <RenderWidgets widget={name} />
              </div>
            </div>
          </Route>
        ))}
      </Switch>
    </Container>
  );
});

const Container = styled.div<{ isMobile: boolean }>`
  flex-grow: 1;
	margin: ${(p: any) => (p.isMobile ? '0 -20px' : '0')};
`;

const Tabs = styled.div`
  display: flex;
  width: 100%;
  align-items: center;
  overflow-x: auto;
  ::-webkit-scrollbar {
    display: none;
  }
`;
interface TagProps {
  selected: boolean;
}

const Tab = styled.div<TagProps>`
  display: flex;
  padding: 10px;
  margin-right: 25px;
	color: ${(p: any) => (p.selected ? '#292C33' : '#8E969C')};
  border-bottom: ${(p: any) => p.selected && '4px solid #618AFF'};
  cursor: hover;
  font-weight: 500;
  font-size: 15px;
  line-height: 19px;
  cursor: pointer;
`;

const Counter = styled.div`
  font-family: Roboto;
  font-style: normal;
  font-weight: normal;
  font-size: 11px;
  line-height: 19px;
  margin-bottom: -3px;
  /* or 173% */
  margin-left: 8px;

  display: flex;
  align-items: center;

  /* Placeholder Text */

  color: #b0b7bc;
`;
