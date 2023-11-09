import React from 'react';
import styled from 'styled-components';
import moment from 'moment';
import { useHistory } from 'react-router';
import { observer } from 'mobx-react-lite';
import { NameTagProps } from 'people/interfaces';
import { useStores } from '../../store';
import { useIsMobile } from '../../hooks';
import { colors } from '../../config/colors';

interface ImageProps {
  readonly src: string;
  iconSize?: number;
  isPaid?: boolean;
}
interface NameProps {
  textSize?: number;
  color?: string;
}

const Img = styled.div<ImageProps>`
  background-image: url('${(p: any) => p.src}');
  background-position: center;
  background-size: cover;
  height: ${(p: any) => (p.iconSize ? `${p.iconSize}px` : '16px')};
  width: ${(p: any) => (p.iconSize ? `${p.iconSize}px` : '16px')};
  border-radius: 50%;
  position: relative;
  opacity: ${(p: any) => (p.isPaid ? 0.3 : 1)};
  filter: ${(p: any) => p.isPaid && 'grayscale(100%)'};
`;

const Name = styled.div<NameProps>`
  font-family: 'Barlow';
  font-style: normal;
  font-weight: normal;
  font-size: ${(p: any) => (p.textSize ? `${p.textSize}px` : '13px')};
  color: ${(p: any) => p.color};
  line-height: 16px;
  /* or 158% */

  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  width: 11vw;
`;

const Date = styled.div`
  font-family: 'Barlow';
  font-style: normal;
  font-weight: normal;
  font-size: 13px;
  line-height: 19px;
  display: flex;
  align-items: center;
  color: #b0b7bc;
`;
interface WrapProps {
  readonly isSelected: boolean;
}

const Wrap = styled.div<WrapProps>`
  display: flex;
  align-items: center;
  cursor: ${(p: any) => !p.isSelected && 'pointer'};
  width: fit-content;
  margin-bottom: 10px;
  color: #8e969c;
  // &:hover {
  //   color: ${(p: any) => !p.isSelected && '#618AFF'};
  // }
`;

function NameTag(props: NameTagProps) {
  const { owner_alias, owner_pubkey, img, created, id, style, widget, iconSize, textSize, isPaid } =
    props;
  const { ui } = useStores();
  const color = colors['light'];

  const history = useHistory();

  const isMobile = useIsMobile();

  const isSelected = ui.selectedPerson === id ? true : false;

  function selectPerson(e: any) {
    // don't select if already selected
    if (isSelected) return;
    e.stopPropagation();

    ui.setPersonViewOpenTab(widget || '');
    ui.setSelectedPerson(id);
    ui.setSelectingPerson(id);

    if (owner_pubkey) {
      history.push(`/p/${owner_pubkey}`);
    }
  }

  let lastSeen = created ? moment.unix(created).fromNow() : '';

  // shorten lastSeen string
  if (lastSeen === 'a few seconds ago') lastSeen = 'just now';

  if (isMobile) {
    return (
      <Wrap
        isSelected={isSelected}
        onClick={(e: any) => {
          selectPerson(e);
        }}
        style={style}
      >
        {!isSelected && (
          <>
            <Img src={img || `/static/person_placeholder.png`} iconSize={iconSize} />
            <Name
              textSize={textSize}
              color={color.grayish.G250}
              style={{
                marginLeft: '10px'
              }}
            >
              {owner_alias}
            </Name>
            <div
              style={{
                height: 3,
                width: 3,
                borderRadius: '50%',
                margin: '0 6px',
                background: color.grayish.G100
              }}
            />
          </>
        )}
        <Date>{lastSeen}</Date>
      </Wrap>
    );
  }

  return (
    <Wrap isSelected={isSelected} style={style}>
      <div
        style={{
          display: 'flex',
          flexDirection: 'row'
        }}
      >
        <Img src={img || `/static/person_placeholder.png`} iconSize={32} isPaid={isPaid} />
        <div
          style={{
            display: 'flex',
            flexDirection: 'column',
            marginLeft: '14px'
          }}
        >
          <Name
            textSize={textSize}
            color={isPaid ? color.grayish.G300 : color.pureBlack}
            onClick={(e: any) => {
              selectPerson(e);
            }}
          >
            {owner_alias}
          </Name>
          <Date>{lastSeen}</Date>
          {}
        </div>
      </div>
    </Wrap>
  );
}
export default observer(NameTag);
