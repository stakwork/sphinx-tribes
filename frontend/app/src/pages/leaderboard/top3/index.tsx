import { EuiAvatar, EuiText } from '@elastic/eui';
import { PriceOuterContainer } from 'components/common';
import MaterialIcon from '@material/react-material-icon';
import { colors } from 'config';
import { DollarConverter } from 'helpers';
import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { useStores } from 'store';
import { Person } from 'store/main';
import styled from 'styled-components';
import { LeaderItem } from '../store';

const color = colors.light;
const ItemContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  & .name {
    padding: 0 0.5rem;
    font-size: 20px;
    font-weight: 600;
    color: rgb(60, 63, 65);
    cursor: pointer;
    text-decoration: none;
  }
`;

const Top3Container = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: flex-end;
`;
const colorsDispatcher = {
  1: '#FFD700',
  2: '#C0C0C0',
  3: '#CD7F32'
};
const Podium = styled.div<{ place: number }>`
  --height: 300px;
  width: 100%;
  height: ${(p: any) => `calc(var(--height) / ${p.place * 1.5})`};
  background-color: ${(p: any) => colorsDispatcher[p.place]};
`;

type ItemProps = LeaderItem & {
  place: number;
};

const Item = ({ owner_pubkey, total_sats_earned, place }: ItemProps) => {
  const { main } = useStores();
  const [person, setPerson] = useState<Person>();

  useEffect(() => {
    main.getPersonByPubkey(owner_pubkey).then(setPerson);
  }, [owner_pubkey, main]);

  return (
    <ItemContainer
      style={{
        order: place === 1 ? 2 : place === 3 ? 3 : 1
      }}
    >
      <EuiAvatar
        size="xl"
        name={person?.owner_alias || ''}
        imageUrl={person?.img || '/static/person_placeholder.png'}
      />
      <div>
        <EuiText textAlign="center" className="name">
          {!!person?.owner_alias && (
            <Link className="name" to={`/p/${person.owner_pubkey}`}>
              {person.owner_alias}
              <MaterialIcon className="icon" icon="link" />
            </Link>
          )}
        </EuiText>
        <PriceOuterContainer
          price_Text_Color={color.primaryColor.P300}
          priceBackground={color.primaryColor.P100}
        >
          <div className="Price_inner_Container">
            <EuiText className="Price_Dynamic_Text">{DollarConverter(total_sats_earned)}</EuiText>
          </div>
          <div className="Price_SAT_Container">
            <EuiText className="Price_SAT_Text">SAT</EuiText>
          </div>
        </PriceOuterContainer>
      </div>
      <Podium place={place} />
    </ItemContainer>
  );
};

export const Top3 = () => {
  const { leaderboard } = useStores();
  return (
    <Top3Container>
      {leaderboard.top3.map((item: any, index: number) => (
        <Item place={index + 1} key={item.owner_pubkey} {...item} />
      ))}
    </Top3Container>
  );
};
