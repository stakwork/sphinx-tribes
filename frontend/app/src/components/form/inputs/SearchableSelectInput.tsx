import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { EuiIcon } from '@elastic/eui';
import { observer } from 'mobx-react-lite';
import { SearchableSelect } from '../../common';
import { useStores } from '../../../store';
import { colors } from '../../../config/colors';
import type { Props } from './propsType';
import { FieldEnv, Note } from './index';

interface styledProps {
  color?: any;
}

const ExtraText = styled.div`
  padding: 2px 10px 25px 10px;
  max-width: calc(100 % - 20px);
  word-break: break-all;
  font-size: 14px;
`;

const E = styled.div<styledProps>`
  position: absolute;
  right: 10px;
  top: 0px;
  display: flex;
  height: 100%;
  justify-content: center;
  align-items: center;
  color: ${(p: any) => p?.color && p?.color.blue3};
  pointer-events: none;
  user-select: none;
`;
const R = styled.div`
  position: relative;
`;
function SearchableSelectInput({
  error,
  note,
  name,
  label,
  options,
  value,
  handleChange,
  extraHTML
}: Props) {
  const labeltext = label;

  const { main, ui } = useStores();
  const color = colors['light'];

  const [opts, setOptions] = useState(options);
  const [loading, setLoading] = useState(false);
  const [search, setSearch]: any = useState('');
  const [isBorder, setIsBorder] = useState<boolean>(true);

  useEffect(() => {
    (async () => {
      if (search) {
        setLoading(true);
        try {
          if (name === 'assignee' || name === 'recipient') {
            const p = await main.getPeopleByNameAliasPubkey(search);
            if (p && p.length) {
              const newOpts = p.map((ot: any) => ({
                owner_alias: ot.owner_alias,
                owner_pubkey: ot.owner_pubkey,
                img: ot.img,
                value: ot.owner_pubkey,
                label: `${ot.owner_alias} (${ot.unique_name})`
              }));
              setOptions(newOpts);
            }
          } else if (name === 'badge') {
            const { badgeList } = ui;

            if (badgeList && badgeList.length) {
              const newOpts = badgeList.map((ot: any) => ({
                img: ot.icon,
                id: ot.id,
                token: ot.token,
                amount: ot.amount,
                value: ot.asset,
                asset: ot.asset,
                label: `${ot.name} (${ot.amount}) `
              }));
              setOptions(newOpts);
            }
          }
        } catch (e) {
          console.log('e', e);
        }
        setLoading(false);
      }
    })();
  }, [search, main, name, ui]);

  return (
    <>
      <FieldEnv
        color={color}
        label={labeltext}
        isTop={true}
        style={{
          border: isBorder ? `1px solid ${color.grayish.G600} ` : `1px solid ${color.pureWhite}`
        }}
      >
        <R>
          <SearchableSelect
            selectStyle={{ border: 'none' }}
            options={opts}
            value={value}
            loading={loading}
            onChange={(e: any) => {
              handleChange(e);
              setIsBorder(false);
            }}
            onInputChange={(e: any) => {
              if (e) setSearch(e);
            }}
          />
          {error && (
            <E color={color}>
              <EuiIcon type="alert" size="m" style={{ width: 20, height: 20 }} />
            </E>
          )}
        </R>
      </FieldEnv>
      {note && <Note color={color}>*{note}</Note>}
      <ExtraText
        style={{ display: value && extraHTML ? 'block' : 'none' }}
        dangerouslySetInnerHTML={{ __html: extraHTML || '' }}
      />
    </>
  );
}

export default observer(SearchableSelectInput);
