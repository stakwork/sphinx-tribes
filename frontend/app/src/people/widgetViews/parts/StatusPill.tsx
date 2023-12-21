import { GithubStatusPillProps } from 'people/interfaces';
import React, { useEffect, useState } from 'react';
import styled from 'styled-components';

interface PillProps {
  readonly isOpen: boolean;
}
const Pill = styled.div<PillProps>`
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 12px;
  font-weight: 300;
  background: ${(p: any) => (p.isOpen ? '#49C998' : '#8256D0')};
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
export default function GithubStatusPill(props: GithubStatusPillProps) {
  const { status, assignee, style } = props;

  const [assigneText, setAssigneText] = useState('');

  const isOpen = status === 'open' || !status;

  useEffect(() => {
    const assignedText =
      assignee && !assignee?.owner_alias
        ? 'Not assigned'
        : isOpen
        ? 'Assigned to '
        : 'Completed by ';
    setAssigneText(assignedText);
  }, [isOpen, assignee]);

  return (
    <div style={{ display: 'flex', ...style }}>
      <Pill isOpen={isOpen}>{/* <div>{isOpen ? 'Open' : 'Closed'}</div> */}</Pill>
      <W>
        <Assignee>{assigneText}</Assignee>
      </W>
    </div>
  );
}
