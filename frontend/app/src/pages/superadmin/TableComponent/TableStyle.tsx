
import styled from 'styled-components';

export const TableContainer = styled.div`
  background-color: #fff;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
`;

export const HeaderContainer = styled.div`
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  padding-right: 40px;
  padding-left: 20px;
`;

export const PaginatonSection = styled.div`
  background-color: #fff;
  height: 64px;
  flex-shrink: 0;
  align-self: stretch;
  border-radius: 8px;
  padding: 1em;
`;

export const Header = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

export const Table = styled.table`
  width: 100%;
  border-collapse: collapse;
`;

export const TableRow = styled.tr`
  border: 1px solid #ddd;
  &:nth-child(even) {
    background-color: #f9f9f9;
  }
`;

export const TableData = styled.td`
  padding: 12px;
  text-align: left;
  white-space: wrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 400px;
  font-size: 14px;
  padding-right: 2em;
  padding-left: 2em;
`;

export const TableDataCenter = styled.td`
  padding: 12px;
  text-align: center;
  white-space: wrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 300px;
  font-size: 14px;
  padding-right: 2em;
  padding-left: 2em;
`;

export const TableData2 = styled.td`
  white-space: wrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-size: 14px;
  padding-right: 3em;
  padding-top: 22px;
  display: flex;
  justify-content: end;
`;

export const TableHeaderData = styled.th`
  padding: 12px;
  text-align: left;
  padding-left: 26px;
  color: var(--Main-bottom-icons, #5F6368);
  font-family: Barlow;
  font-size: 12px;
  font-style: normal;
  font-weight: 400;
  line-height: normal;
  letter-spacing: 0.72px;
  text-transform: uppercase;
`;

export const TableHeaderDataCenter = styled.th`
  padding: 12px;
  text-align: center;
  color: var(--Main-bottom-icons, #5F6368);
  font-family: Barlow;
  font-size: 12px;
  font-style: normal;
  font-weight: 400;
  line-height: normal;
  letter-spacing: 0.72px;
  text-transform: uppercase;
`;

export const TableHeaderDataRight = styled.th`
  padding: 12px;
  text-align: right;
  padding-right: 40px;
  color: var(--Main-bottom-icons, #5F6368);
  font-family: Barlow;
  font-size: 12px;
  font-style: normal;
  font-weight: 400;
  line-height: normal;
  letter-spacing: 0.72px;
  text-transform: uppercase;
`;

export const BountyHeader = styled.div`
  background: #FFF;
  display: flex;
  height: 66px;
  justify-content: space-between;
  text-align: center;
  align-items: center;
  gap: 10px;
  padding-left: 1em;
  padding-right: 2em;
`;

export const Options = styled.div`
  font-size: 15px;
  cursor: pointer;
  outline: none;
  border: none;
  display: flex;
  align-items: center;
  gap: 20px;
`;

export const StyledSelect = styled.select`
  color: var(--Text-2, var(--Hover-Icon-Color, #3C3F41));
  font-family: Barlow;
  font-size: 15px;
  font-style: normal;
  font-weight: 500;
  line-height: normal;
  border-radius: 4px;
  cursor: pointer;
  outline: none;
  border: none;
`;

export const LeadingTitle = styled.h2`
  color: var(--Primary-Text-1, var(--Press-Icon-Color, #292C33));
 
  font-family: Barlow;
  font-size: 15px;
  font-style: normal;
  font-weight: 600;
  display: flex;
  gap: 6px;
  line-height: normal;
  margin-top: 15px;
`;

export const AlternativeTitle = styled.h2`
  color: var(--Main-bottom-icons, #5F6368);
  font-family: Barlow;
  font-size: 15px;
  font-style: normal;
  font-weight: 500;
  line-height: normal;
`;

export const Label = styled.label`
  margin-top: 6px;
  color: var(--Main-bottom-icons, var(--Hover-Icon-Color, #5F6368));
  font-family: Barlow;
  font-size: 15px;
  font-style: normal;
  font-weight: 500;
  line-height: normal;
`;

export const FlexDiv = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 2px;
`;

interface PaginationButtonsProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
    active: boolean;
  }

export const PaginationButtons = styled.button<PaginationButtonsProps>`
  border-radius: 3px;
  width: 30px;
  height: 30px;
  flex-shrink: 0;
  outline: none;
  border: none;
  text-align: center;
  margin: 10px;
  background: ${(props: any) => (props.active ? 'var(--Active-blue, #618AFF)' : 'white')};
  color: ${(props: any) => (props.active ? 'white' : 'black')};
`;

export const PageContainer = styled.div`
  display: flex;
  align-items: center;
`;