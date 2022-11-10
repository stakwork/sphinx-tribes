import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { EuiIcon } from '@elastic/eui';
import type { Props } from './propsType';
import { FieldEnv, Note } from './index';
import { SearchableSelect } from '../../sphinxUI';
import { useStores } from '../../store';

export default function SearchableSelectInput({
  error,
  note,
  name,
  type,
  label,
  options,
  value,
  handleChange,
  handleBlur,
  handleFocus,
  readOnly,
  prepend,
  extraHTML
}: Props) {
  const labeltext = label;

  const { main, ui } = useStores();

  const [opts, setOptions] = useState(options);
  const [loading, setLoading] = useState(false);
  const [search, setSearch]: any = useState('');

  useEffect(() => {
    (async () => {
      if (search) {
        setLoading(true);
        try {
          if (name === 'assignee' || name === 'recipient') {
            const p = await main.getPeopleByNameAliasPubkey(search);
            if (p && p.length) {
              const newOpts = p.map((ot) => {
                return {
                  owner_alias: ot.owner_alias,
                  owner_pubkey: ot.owner_pubkey,
                  img: ot.img,
                  value: ot.owner_pubkey,
                  label: `${ot.owner_alias} (${ot.unique_name})`
                };
              });
              setOptions(newOpts);
            }
          } else if (name === 'badge') {
            const { badgeList } = ui;

            if (badgeList && badgeList.length) {
              const newOpts = badgeList.map((ot) => {
                return {
                  img: ot.icon,
                  id: ot.id,
                  token: ot.token,
                  amount: ot.amount,
                  value: ot.asset,
                  asset: ot.asset,
                  label: `${ot.name} (${ot.amount}) `
                };
              });
              setOptions(newOpts);
            }
          }
        } catch (e) {
          console.log('e', e);
        }
        setLoading(false);
      }
    })();
  }, [search]);

  return (
    <>
      <FieldEnv label={labeltext}>
        <R>
          <SearchableSelect
            selectStyle={{ border: 'none' }}
            options={opts}
            value={value}
            loading={loading}
            onChange={(e) => {
              handleChange(e);
            }}
            onInputChange={(e) => {
              if (e) setSearch(e);
            }}
          />
          {error && (
            <E>
              <EuiIcon type="alert" size="m" style={{ width: 20, height: 20 }} />
            </E>
          )}
        </R>
      </FieldEnv>
      {note && <Note>*{note}</Note>}
      <ExtraText
        style={{ display: value && extraHTML ? 'block' : 'none' }}
        dangerouslySetInnerHTML={{ __html: extraHTML || '' }}
      />
    </>
  );
}

const ExtraText = styled.div`
padding: 2px 10px 25px 10px;
max - width: calc(100 % - 20px);
word -break: break-all;
font - size: 14px;
`;

const E = styled.div`
position: absolute;
right: 10px;
top: 0px;
display: flex;
height: 100 %;
justify - content: center;
align - items: center;
color:#45b9f6;
pointer - events: none;
user - select: none;
`;
const R = styled.div`
  position: relative;
`;
