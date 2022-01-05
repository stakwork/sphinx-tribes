import React from 'react'
import styled from "styled-components";
import MaterialIcon from '@material/react-material-icon';

export default function GithubStatusPill(props: any) {
    const { status, assignee, style } = props

    const isOpen = status === 'open'

    return <div style={{ display: 'flex', ...style }}>
        <Pill isOpen={isOpen}>
            <MaterialIcon style={{
                // marginRight: 2,
                fontSize: 14
            }} icon={isOpen ? "arrow_circle_up" : "check_circle_outline"} />
            {/* <div>
                {status || 'Open'}
            </div> */}
        </Pill>
        <Assignee>
            {(assignee && `Assigned to ${assignee}`) || 'Not Assigned'}
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
font-size:12px;
font-weight:300;
background:${p => p.isOpen ? '#347d39' : '#8256d0'}; //#26a641
border-radius:30px;
border: 1px solid transparent;
text-transform: capitalize;
padding: 12px 5px;
// padding:8px;
font-size: 12px;
font-weight: 500;
line-height: 20px;
white-space: nowrap;
border-radius: 2em;
height:26px;
color:#fff;
margin-right:5px;
`;

const Assignee = styled.div`
display: flex;
justify-content:center;
align-items:center;
font-size:12px;
font-weight:300;
color:#8E969C;
`
