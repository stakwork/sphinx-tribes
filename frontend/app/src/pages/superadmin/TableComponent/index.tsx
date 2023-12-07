

import React,{useState} from 'react';
import styled from "styled-components";
import paginationarrow1 from "../Header/icons/paginationarrow1.svg"
import paginationarrow2 from "../Header/icons/paginationarrow2.svg"
import copy from "../Header/icons/copy.svg"
import {
  TableContainer,
  HeaderContainer,
  PaginatonSection,
  Header,
  Table,
  TableRow,
  TableData,
  TableDataCenter,
  TableData2,
  TableHeaderData,
  TableHeaderDataCenter,
  TableHeaderDataRight,
  BountyHeader,
  Options,
  StyledSelect,
  LeadingTitle,
  AlternativeTitle,
  Label,
  FlexDiv,
  PaginationButtons,
  PageContainer,
  // Import other styled components as needed
} from "./TableStyle"

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

interface TableProps {
  bounties: Bounty[];
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


export const MyTable = ({
  bounties,
  }:TableProps) => {
    
    const [currentPage , setCurrentPage] = useState(1);
    const pageSize = 10;

    const dataNumber: number[] = [];
    
    for (let i = 1; i <= Math.ceil(bounties.length / pageSize); i++) {
      dataNumber.push(i);
    }

  
    const currentPageData = () => {
      const indexOfLastPost = currentPage * pageSize;
      const indexOfFirstPost = indexOfLastPost - pageSize;
      const currentPosts = bounties.slice(indexOfFirstPost, indexOfLastPost);
      return currentPosts;
    };

    const paginateNext = () => {
      console.log("clicked")
      if (currentPage < bounties?.length / pageSize) {
        setCurrentPage(currentPage + 1);
      }
    };
    const paginatePrev = () => {
      console.log("clicked")
      if (currentPage > 1) {
        setCurrentPage(currentPage - 1);
      }
    };
  
    return (
      <>
      <HeaderContainer>
        <Header>
          <BountyHeader>
             <img src={copy} alt="" />
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
            <TableHeaderDataCenter>#DTGP</TableHeaderDataCenter>
            <TableHeaderData>Assignee</TableHeaderData>
            <TableHeaderData>Provider</TableHeaderData>
            <TableHeaderData>Organization</TableHeaderData>
            <TableHeaderDataRight>Status</TableHeaderDataRight>
          </TableRow>
          <tbody>
            {currentPageData()?.map((bounty:any) => (
              <TableRow key={bounty?.id}>
                <TableData>{bounty?.title}</TableData>
                <TableData>{bounty?.date}</TableData>
                <TableDataCenter>{bounty?.dtgp}</TableDataCenter>
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
      {bounties.length > pageSize ? 
        <PageContainer>
          <img src={paginationarrow1} alt="" onClick={() => paginatePrev()} />
          {dataNumber.map((number: number) => (
            <PaginationButtons
              key={number}
              onClick={() => setCurrentPage(number)}
              active={number === currentPage}
            >
              {number}
            </PaginationButtons>
          ))}
          <img src={paginationarrow2} alt="" onClick={() => paginateNext()} />
        </PageContainer>
        : null }
    
      </FlexDiv>
      </PaginatonSection>
      </>
    )
  }


