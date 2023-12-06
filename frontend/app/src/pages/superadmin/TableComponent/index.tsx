

import React from 'react';
import styled from "styled-components";
import paginationarrow1 from "../Header/icons/paginationarrow1.svg"
import paginationarrow2 from "../Header/icons/paginationarrow2.svg"


const TableContainer = styled.div`
  background-color: #fff;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
 
`;

const HeaderContainer = styled.div`
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  padding-right: 40px;
  padding-left: 20px;
`;

const PaginatonSection = styled.div`
  background-color: #fff;
  height: 64px;
  flex-shrink: 0;
  align-self: stretch;
  border-radius: 8px;
  padding:1em;
`;

const Header = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

const Table = styled.table`
  width: 100%;
  border-collapse: collapse;

`;

const TableRow = styled.tr`
  border: 1px solid #ddd;
  &:nth-child(even) {
    background-color: #f9f9f9;
  }
  
`;

const TableData = styled.td`
  padding: 12px;
  text-align: left;
  white-space: wrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 300px;
  font-size: 14px;
  padding-right: 2em;
  padding-left: 2em;
`;

const TableData2 = styled.td`
  padding: 12px;
  white-space: wrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 300px;
  font-size: 14px;
  padding-right: 3em;
  padding-left: 2em;
`;

const TableHeaderData = styled.th`
  padding: 12px;
  text-align: left;
  padding-left: 26px; /* Reduce padding-left */
  color: var(--Main-bottom-icons, #5F6368);
  leading-trim: both;
  text-edge: cap;
  font-family: Barlow;
  font-size: 12px;
  font-style: normal;
  font-weight: 400;
  line-height: normal;
  letter-spacing: 0.72px;
  text-transform: uppercase;
`;

const TableHeaderDataRight = styled.th`
  padding: 12px;
  text-align: right;
  padding-right: 40px; /* Reduce padding-right */
 /* Adjust padding-left */
  color: var(--Main-bottom-icons, #5F6368);
  leading-trim: both;
  text-edge: cap;
  font-family: Barlow;
  font-size: 12px;
  font-style: normal;
  font-weight: 400;
  line-height: normal;
  letter-spacing: 0.72px;
  text-transform: uppercase;
`;

const BountyHeader = styled.div`
  background: #FFF;
  display: flex;
  height: 66px;
  justify-content: space-between;
  text-align:center;
  align-items: center;
  gap: 10px;
  padding-left: 1em;
  padding-right: 2em;
`;

const Options = styled.div`
  font-size: 15px;
  cursor: pointer;
  outline: none;
  border: none;
  display: flex;
  align-items: center;
  gap: 20px;
`;

const StyledSelect = styled.select`
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


const LeadingTitle =styled.h2`
  color: var(--Primary-Text-1, var(--Press-Icon-Color, #292C33));
  leading-trim: both;
  text-edge: cap;
  font-family: Barlow;
  font-size: 15px;
  font-style: normal;
  font-weight: 600;
  display:flex;
  gap:6px;
  line-height: normal;
`
const AlternativeTitle =styled.h2`
  color: var(--Main-bottom-icons, #5F6368);
  leading-trim: both;
  text-edge: cap;
  font-family: Barlow;
  font-size: 15px;
  font-style: normal;
  font-weight: 500;
  line-height: normal;
`


interface Bounty {
  id: number;
  title: string;
  date: string;
  dtgp: number;
  assignee: string;
  assigneeImage: string;
  provider: string;
  providerImage: string;
  organization: string;
  organizationImage: string;
  status: string;
}

const Label = styled.label`
  margin-top:6px;
  color: var(--Main-bottom-icons, var(--Hover-Icon-Color, #5F6368));
  font-family: Barlow;
  font-size: 15px;
  font-style: normal;
  font-weight: 500;
  line-height: normal;
`  


const FlexDiv = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 2px;
`

interface PaginationButtonProps {
  selected: boolean;
}

const PaginationButtons = styled.button<PaginationButtonProps>`
  border-radius: 3px;
  background: ${(props:any) => (props.selected ? 'var(--Primary-blue, #618AFF)' : 'transparent')};
  width: 34px;
  height: 34px;
  flex-shrink: 0;
  outline:none;
  border:none;
  color:${(props:any) => (props.selected ? 'white' : 'gray')};
`

interface TableProps {
  bounties: Bounty[];
  selectedButtonIndex: number;
}

interface ImageWithTextProps {
    image?: string;
    text: string;
  }

export const ImageWithText = ({ image, text }: ImageWithTextProps) => {
    const BoxImage = styled.div`
      display: flex;
      width: 162px;
      align-items: center;
      text-align:center;
      gap: 6px;
    `;
    const Paragraph =styled.div`
      margin-top:2px;
    `
    return (
      <>
        <BoxImage>
          <img
            src={image}
            style={{
              width: "30px",
              height: "30px",
              borderRadius: "50%",
              marginRight: "10px",
            }}
            alt={text}
          />
          <Paragraph>{text}</Paragraph>
        </BoxImage>
      </>
    );
  };
  
  interface TextInColorBoxProps {
    status: string;
  }


  export const TextInColorBox = ({ status }: TextInColorBoxProps) => (
      <>
      <div  style={{
              display:"flex",
              justifyContent:"flex-end",
            
            }}
            >
           <p
              style={{
                color: "#fff",
                paddingRight: "10px",
                paddingLeft: "10px",
                width: "max-content",
                textAlign: "right",
                backgroundColor:
                  status === "open"
                    ? "#618AFF"
                    : status === "paid"
                    ? "#5F6368"
                    : status === "assigned" // Fix the value here
                    ? "#49C998"
                    : status === "completed"
                    ? "#9157F6"
                    : "transparent", // Add a default value or handle other cases
                borderRadius: "2px",
              }}
            >
            {status}
           </p>
      </div>
      </>
    );


export const MyTable: React.FC<TableProps> = ({bounties, selectedButtonIndex  }:TableProps) => (
  
    <>
      <HeaderContainer>
        <Header>
          <BountyHeader>
            
            <LeadingTitle>  {bounties.length} <AlternativeTitle> {bounties.length === 1 ? "Bounty" : "Bounties"}</AlternativeTitle> </LeadingTitle>
  
          </BountyHeader>
          <Options>
            <FlexDiv>
            <Label>
              Sort By:
            </Label>
              <StyledSelect id="sortBy">
                <option value="date">Date</option>
                <option value="assignee">Assignee</option>
                <option value="status">Status</option>
              </StyledSelect>
            </FlexDiv>
            <FlexDiv>
            <Label>
              Status:
            </Label>
            <StyledSelect id="statusFilter">
              <option value="open">Open</option>
              <option value="in-progress">In Progress</option>
              <option value="completed">Completed</option>
            </StyledSelect>
            </FlexDiv>
          </Options>
        </Header>
      </HeaderContainer>
      <TableContainer>
        <Table>
          <TableRow>
            <TableHeaderData>Bounty</TableHeaderData>
            <TableHeaderData>Date</TableHeaderData>
            <TableHeaderData>#DTGP</TableHeaderData>
            <TableHeaderData>Assignee</TableHeaderData>
            <TableHeaderData>Provider</TableHeaderData>
            <TableHeaderData>Organization</TableHeaderData>
            <TableHeaderDataRight>Status</TableHeaderDataRight>
          </TableRow>
          <tbody>
            {bounties?.map((bounty:any) => (
              <TableRow key={bounty?.id}>
                <TableData>{bounty?.title}</TableData>
                <TableData>{bounty?.date}</TableData>
                <TableData>{bounty?.dtgp}</TableData>
                <TableData>
                  <ImageWithText
                    text={bounty?.assignee}
                    image={bounty?.assigneeImage}
                  />
                </TableData>
                <TableData>
                  <ImageWithText
                    text={bounty?.provider}
                    image={bounty?.providerImage}
                  />
                </TableData>
                <TableData>
                  <ImageWithText
                    text={bounty?.organization}
                    image={bounty?.organizationImage}
                  />
                </TableData>
                <TableData2>
                  <TextInColorBox status={bounty?.status} />
                </TableData2>
              </TableRow>
            ))}
          </tbody>
        </Table>
      </TableContainer>
      <PaginatonSection>
      <FlexDiv>
        <img src={paginationarrow1} alt="" />
        <PaginationButtons selected={selectedButtonIndex === 1}>1</PaginationButtons>
        <PaginationButtons selected={selectedButtonIndex === 2}>2</PaginationButtons>
        <PaginationButtons selected={selectedButtonIndex === 3}>3</PaginationButtons>
        <PaginationButtons selected={selectedButtonIndex === 4}>4</PaginationButtons>
        <PaginationButtons selected={selectedButtonIndex === 5}>5</PaginationButtons>
        <PaginationButtons selected={selectedButtonIndex === 6}>6</PaginationButtons>
        <PaginationButtons selected={selectedButtonIndex === 7}>7</PaginationButtons>
        <img src={paginationarrow2} alt="" />
      </FlexDiv>
      </PaginatonSection>
    </>
  );


