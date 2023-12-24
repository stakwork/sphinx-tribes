import { StatusPillProps } from 'people/interfaces';
import React, { useEffect, useState } from 'react';
import styled from 'styled-components';

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
export default function StatusPill(props: StatusPillProps) {
  const { assignee, style } = props;

  const [assigneText, setAssigneText] = useState('');

  useEffect(() => {
    const assignedText =
      assignee && !assignee?.owner_alias
        ? 'Not assigned'
        : assignee
        ? 'Assigned to '
        : 'Completed by ';
    setAssigneText(assignedText);
  }, [assignee]);

  return (
    <div style={{ display: 'flex', ...style }}>
      <W>
        <Assignee>{assigneText}</Assignee>
      </W>
    </div>
  );
}
