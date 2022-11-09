import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { useHistory } from 'react-router-dom';
import { useStores } from '../../../store';

export default function GithubStatusPill(props: any) {
  const { status, assignee, style } = props;
  const { main } = useStores();

  const [assigneText, setAssigneText] = useState('');

  const [newAssignee, setNewAssignee] = useState({ ...assignee });

  async function findUserByGithubHandle() {
    // look in database for first user with this github handle
    try {
      const p = await main.getPersonByGithubName(assignee);
      if (p) {
        let url = '';
        if (p.owner_pubkey) url = `https://community.sphinx.chat/p/${p.owner_pubkey}`;
        else url = `https://github.com/${assignee}`;
        sendToRedirect(url);
      }
    } catch (e) {
      console.log('e', e);
    }
  }

  function sendToRedirect(url) {
    const el = document.createElement('a');
    el.href = url;
    el.target = '_blank';
    el.click();
  }

  const isOpen = status === 'open' || !status;

  useEffect(() => {
    const assignedText = !newAssignee.owner_alias
      ? 'Not assigned'
      : isOpen
      ? 'Assigned to '
      : 'Completed by ';
    setAssigneText(assignedText);
  }, [isOpen, newAssignee]);

  return (
    <div style={{ display: 'flex', ...style }}>
      <Pill isOpen={isOpen}>
        <div>{isOpen ? 'Open' : 'Closed'}</div>
      </Pill>
      <W>
        <Assignee>
          {assigneText}{' '}
          {/* <Link
            onClick={(e) => {
              e.stopPropagation();
            }}
          >
            {newAssignee.owner_alias}
          </Link> */}
        </Assignee>
      </W>
    </div>
  );
}
interface PillProps {
  readonly isOpen: boolean;
}
const Pill = styled.div<PillProps>`
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 12px;
  font-weight: 300;
  background: ${(p) => (p.isOpen ? '#49C998' : '#8256D0')};
  border-radius: 30px;
  border: 1px solid transparent;
  text-transform: capitalize;
  padding: 12px 5px;
  // padding:8px;
  font-size: 12px;
  font-weight: 500;
  line-height: 20px;
  white-space: nowrap;
  border-radius: 2em;
  height: 26px;
  color: #fff;
  margin-right: 10px;
  width: 58px;
  height: 22px;
  left: 19px;
  top: 171px;

  /* Primary Green */

  border-radius: 2px;
`;

const Assignee = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 12px;
  font-weight: 300;
  color: #8e969c;

  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;

  overflow: hidden;
`;

const W = styled.div`
  display: flex;
  align-items: center;
`;

const Link = styled.span`
  font-weight: 500;
  cursor: pointer;
  &:hover {
    color: #000;
  }
`;
