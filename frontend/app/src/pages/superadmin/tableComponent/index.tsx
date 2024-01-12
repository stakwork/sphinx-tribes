import React, { useState } from 'react';
import styled from 'styled-components';
import { colors } from 'config';
import { EuiPopover, EuiText, EuiCheckboxGroup } from '@elastic/eui';
import expand_more from '../header/icons/expand_more.svg';
import paginationarrow1 from '../header/icons/paginationarrow1.svg';
import paginationarrow2 from '../header/icons/paginationarrow2.svg';
import copygray from '../header/icons/copygray.svg';
import {
  TableContainer,
  HeaderContainer,
  PaginatonSection,
  Header,
  Table,
  TableRow,
  TableData,
  TableDataCenter,
  TableData3,
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
  TableHeaderDataAlternative,
  TableDataRow,
  TableDataAlternative,
  BountyData
} from './TableStyle';

import './styles.css';

import { FilterContainer, FlexDivStatus, StatusCheckboxItem } from './StatusStyle';

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
    text-align: center;
    gap: 6px;
  `;
  const Paragraph = styled.div`
    margin-top: 2px;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
    max-width: 200px;
  `;
  return (
    <>
      <BoxImage>
        <img
          src={image}
          style={{
            width: '30px',
            height: '30px',
            borderRadius: '50%',
            marginRight: '10px'
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
    <div
      style={{
        display: 'flex',
        justifyContent: 'flex-end'
      }}
    >
      <p
        style={{
          color: '#fff',
          textTransform: 'uppercase',
          paddingRight: '10px',
          paddingLeft: '10px',
          width: 'max-content',
          textAlign: 'right',
          backgroundColor:
            status === 'Open'
              ? '#618AFF'
              : status === 'Paid'
              ? '#5F6368'
              : status === 'Assigned'
              ? '#49C998'
              : status === 'Completed'
              ? '#9157F6'
              : 'transparent',
          borderRadius: '2px',
          marginBottom: '0'
        }}
      >
        {status}
      </p>
    </div>
  </>
);

export const MyTable = ({ bounties }: TableProps) => {
  const [currentPage, setCurrentPage] = useState(1);
  const pageSize = 10;
  const [statusFilter, setStatusFilter] = useState<string[]>([]);
  const [isPopoverOpen, setIsPopoverOpen] = useState(false);
  const [checkboxIdToSelectedMap, setCheckboxIdToSelectedMap] = useState({
    open: false,
    assigned: false,
    completed: false,
    paid: false
  });

  const color = colors['light'];

  const options = [
    { id: 'Open', label: 'Open', value: 'Open'},
    { id: 'Assigned', label: 'Assigned', value: 'Assigned' },
    { id: 'Completed', label: 'Completed', value: 'Completed' },
    { id: 'Paid', label: 'Paid', value: 'Paid' }
  ];

  const dataNumber: number[] = [];

  for (let i = 1; i <= Math.ceil(bounties.length / pageSize); i++) {
    dataNumber.push(i);
  }

  const statusFilterMap = () => {
    if (statusFilter.length === 0) {
      return bounties;
    }
    return bounties.filter((bounty: any) => statusFilter.includes(bounty.status));
  }

  const currentPageData = () => {
    const firstIndex = (currentPage - 1) * pageSize;
    const lastIndex = firstIndex + pageSize;
    return statusFilterMap().slice(firstIndex, lastIndex);
  };
  

  const paginateNext = () => {
    console.log('clicked');
    if (currentPage < bounties?.length / pageSize) {
      setCurrentPage(currentPage + 1);
    }
  };
  const paginatePrev = () => {
    console.log('clicked');
    if (currentPage > 1) {
      setCurrentPage(currentPage - 1);
    }
  };

  const onButtonClick = () => {
    setIsPopoverOpen(!isPopoverOpen);
  };

  const onChange = (optionId:any) => {
    const newCheckboxIdToSelectedMap = {
      ...checkboxIdToSelectedMap,
      ...{
        [optionId]: !checkboxIdToSelectedMap[optionId],
      },
    };
  
    setCheckboxIdToSelectedMap(newCheckboxIdToSelectedMap);
  
    // Check if the checkbox is being selected or deselected
    if (newCheckboxIdToSelectedMap[optionId]) {
      // If it's being selected, add the optionId to the statusFilter array
      setStatusFilter([...statusFilter, optionId]);
    } else {
      // If it's being deselected, remove the optionId from the statusFilter array
      setStatusFilter(statusFilter.filter((id : any) => id !== optionId));
    }
  };

  return (
    <>
      <HeaderContainer>
        <Header>
          <BountyHeader>
            <img src={copygray} alt="" width="16.508px" height="20px" />
            <LeadingTitle>
              {' '}
              {bounties.length}{' '}
              <AlternativeTitle> {bounties.length === 1 ? 'Bounty' : 'Bounties'}</AlternativeTitle>{' '}
            </LeadingTitle>
          </BountyHeader>
          <Options>
            <FlexDiv>
              <Label>Sort By:</Label>
              <StyledSelect id="sortBy">
                <option value="date">Date</option>
                <option value="assignee">Assignee</option>
                <option value="status">Status</option>
              </StyledSelect>
            </FlexDiv>
            <FlexDiv>
              <Label>Status:</Label>
              <EuiPopover
                button={
                  <FilterContainer onClick={onButtonClick}>                                       
                    <FlexDivStatus>
                    <EuiText
                      className="statusText"
                    >
                      Status:
                    </EuiText>
                    <EuiText
                      className="subStatusText"
                    >
                      {statusFilter.length === 0 ? 'All' : statusFilter.length === 1 ? statusFilter : 'Multiple'}
                    </EuiText>
                    <img src={expand_more} alt="" width="20px" height="20px" />
                    </FlexDivStatus>
                    
                  </FilterContainer>
                }
                panelStyle={{
                  maxWidth: '162px',
                  maxHeight: '168px',
                  borderRadius: '6px',
                  fontSize: '15px',
                  lineHeight: '18px',
                  fontWeight: '500',
                  color: '#5F6368',
                  fontFamily: 'Barlow',
                  border: '1px solid #fff',

                }}
                isOpen={isPopoverOpen}
                closePopover={() => setIsPopoverOpen(false)}
                panelClassName="yourClassNameHere"
                panelPaddingSize="none"
                anchorPosition="downCenter"
              >
               <StatusCheckboxItem color={color}>
               <EuiCheckboxGroup
                  options={options}
                  onChange={(id:any) => onChange(id)}
                  idToSelectedMap={checkboxIdToSelectedMap}
               />
               </StatusCheckboxItem>
              </EuiPopover>      
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
            <TableHeaderDataAlternative>Organization</TableHeaderDataAlternative>
            <TableHeaderDataRight>Status</TableHeaderDataRight>
          </TableRow>
          <tbody>
            {currentPageData()?.map((bounty: any) => (
              <TableDataRow key={bounty?.id}>
                <BountyData className="avg">{bounty?.title}</BountyData>
                <TableData>{bounty?.date}</TableData>
                <TableDataCenter>{bounty?.dtgp}</TableDataCenter>
                <TableDataAlternative>
                  <ImageWithText text={bounty?.assignee} image={bounty?.assigneeImage} />
                </TableDataAlternative>
                <TableDataAlternative className="address">
                  <ImageWithText text={bounty?.provider} image={bounty?.providerImage} />
                </TableDataAlternative>
                <TableData className="organization">
                  <ImageWithText text={bounty?.organization} image={bounty?.organizationImage} />
                </TableData>
                <TableData3>
                  <TextInColorBox status={bounty?.status} />
                </TableData3>
              </TableDataRow>
            ))}
          </tbody>
        </Table>
      </TableContainer>
      <PaginatonSection>
        <FlexDiv>
          {bounties.length > pageSize ? (
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
          ) : null}
        </FlexDiv>
      </PaginatonSection>
    </>
  );
};