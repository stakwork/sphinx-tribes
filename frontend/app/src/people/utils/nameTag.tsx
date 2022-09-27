import React from 'react';
import styled from 'styled-components';
import moment from 'moment';
import { useStores } from '../../store';
import { useHistory } from 'react-router';

export default function NameTag(props) {
  const {
    owner_alias,
    owner_pubkey,
    img,
    created,
    unique_name,
    id,
    style,
    widget,
    iconSize,
    textSize
  } = props;
  const { ui } = useStores();

  const history = useHistory();

  const isSelected = ui.selectedPerson == id ? true : false;

  function selectPerson(e) {
    // don't select if already selected
    if (isSelected) return;
    e.stopPropagation();
    console.log('selectPerson', id, unique_name);
    ui.setPersonViewOpenTab(widget || '');
    ui.setSelectedPerson(id);
    ui.setSelectingPerson(id);
    if (owner_pubkey) {
      history.push(`/p/${owner_pubkey}`);
      // window.history.pushState({}, 'Sphinx Tribes', '/p/' + unique_name);
    }
  }

  let lastSeen = created ? moment.unix(created).fromNow() : '';

  // shorten lastSeen string
  if (lastSeen === 'a few seconds ago') lastSeen = 'just now';

  return (
    <Wrap
      isSelected={isSelected}
      onClick={(e) => {
        selectPerson(e);
      }}
      style={style}
    >
      {!isSelected && (
        <>
          <Img src={img || `/static/person_placeholder.png`} iconSize={iconSize} />
          <Name textSize={textSize}>{owner_alias}</Name>

          <div
            style={{
              height: 3,
              width: 3,
              borderRadius: '50%',
              margin: '0 6px',
              background: '#8E969C'
            }}
          />
        </>
      )}

      <Date>{lastSeen}</Date>
    </Wrap>
  );
}

interface ImageProps {
  readonly src: string;
  iconSize?: number;
}
interface NameProps {
  textSize?: number;
}
const Img = styled.div<ImageProps>`
  background-image: url('${(p) => p.src}');
  background-position: center;
  background-size: cover;
  height: ${(p) => (p.iconSize ? p.iconSize + 'px' : '16px')};
  width: ${(p) => (p.iconSize ? p.iconSize + 'px' : '16px')};
  border-radius: 50%;
  position: relative;
`;

const Name = styled.div<NameProps>`
  font-family: Roboto;
  font-style: normal;
  font-weight: normal;
  font-size: ${(p) => (p.textSize ? p.textSize + 'px' : '12px')};
  line-height: 19px;
  /* or 158% */
  margin-left: 5px;

  display: flex;
  align-items: center;

  /* Secondary Text 4 */
`;

const Date = styled.div`
  font-family: Roboto;
  font-style: normal;
  font-weight: normal;
  font-size: 12px;
  line-height: 19px;
  /* or 158% */

  display: flex;
  align-items: center;
`;
interface WrapProps {
  readonly isSelected: boolean;
}

const Wrap = styled.div<WrapProps>`
  display: flex;
  align-items: center;
  cursor: ${(p) => !p.isSelected && 'pointer'};
  width: fit-content;
  margin-bottom: 10px;
  color: #8e969c;
  &:hover {
    color: ${(p) => !p.isSelected && '#618AFF'};
  }
`;
