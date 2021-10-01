import React from 'react'
import styled from "styled-components";
import MaterialIcon from '@material/react-material-icon';

export default function GithubStatusPill(props: any) {
    const { status, assignee } = props

    const isOpen = status === 'open'

    if (!status) return <div />

    return <div style={{ display: 'flex' }}>
        <Pill isOpen={isOpen}>
            <MaterialIcon icon={isOpen ? "arrow_circle_up" : "check_circle_outline"} />
            <div>
                {status}
            </div>
        </Pill>
        <Assignee>
            {assignee || 'Not assigned'}
        </Assignee>
    </div>

}
interface PillProps {
    readonly isOpen: boolean;
}
const Pill = styled.div<PillProps>`
display: flex;
justify-content:center;
align-items:center;
height:20px;
font-size:12px;
font-weight:300;
background:${p => p.isOpen ? '#006d32' : '#c93c37'}; //#26a641
border-radius:30px;

border: 1px solid transparent;
text-transform: capitalize;
padding: 12px 8px;
font-size: 14px;
font-weight: 500;
line-height: 20px;
white-space: nowrap;
border-radius: 2em;
height:32px;
color:#fff;
margin-right:10px;
`;

const Assignee = styled.div`
display: flex;
justify-content:center;
align-items:center;
font-size:12px;
font-weight:300;
`
