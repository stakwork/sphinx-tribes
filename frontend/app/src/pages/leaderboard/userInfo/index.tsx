import { EuiAvatar, EuiText } from '@elastic/eui';
import MaterialIcon from '@material/react-material-icon';
import { observer } from 'mobx-react-lite';
import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import styled from 'styled-components';
import { colors } from '../../../config';
import ConnectCard from '../../../people/utils/ConnectCard';
import { useStores } from '../../../store';
import { Person } from '../../../store/main';
const UserItemContainer = styled.div`
  display: flex;
  gap: 1rem;
  align-items: center;
  & .name {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    font-size: 20px;
    font-weight: 600;
    color: rgb(60, 63, 65);
    text-decoration: none;
    .icon {
      transition: transform 0.2s ease;
      transform: rotate(45deg) scale(0.75);
    }
    &:hover {
      .icon {
        transform: rotate(0deg) scale(1);
      }
    }
  }

  & .viewProfile {
    color: ${colors.light.text2};
    display: flex;
  }
`;

export const UserInfo = observer(({ id }: { id: string }) => {
  const { main } = useStores();
  const [person, setPerson] = useState<Person>();
  const [showQR, setShowQR] = useState(false);
  useEffect(() => {
    main.getPersonByPubkey(id).then(setPerson);
  }, [id, main]);

  if (!person) {
    return (
      <UserItemContainer>
        <EuiAvatar size="xl" name={''} imageUrl={'/static/person_placeholder.png'} />
      </UserItemContainer>
    );
  }
  return (
    <UserItemContainer>
      <EuiAvatar
        size="l"
        name={person.owner_alias}
        imageUrl={person.img || '/static/person_placeholder.png'}
      />
      <div className="info">
        <EuiText className="name">
          <Link className="name" to={`/p/${person.owner_pubkey}`}>
            {person.owner_alias}
            <MaterialIcon className="icon" icon="link" />
          </Link>
        </EuiText>
      </div>

      {showQR && <ConnectCard dismiss={() => setShowQR(false)} person={person} visible={showQR} />}
    </UserItemContainer>
  );
});
