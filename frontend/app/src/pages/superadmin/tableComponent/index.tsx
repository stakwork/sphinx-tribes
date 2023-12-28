import React, { useCallback, useEffect, useState } from 'react';
import styled from 'styled-components';
import { useStores } from 'store';
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
  StyledSelect2,
  TableHeaderDataAlternative,
  TableDataRow,
  TableDataAlternative,
  BountyData
} from './TableStyle';

import './styles.css';
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
  startDate?: number;
  endDate?: number;
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
            status === 'open'
              ? '#618AFF'
              : status === 'paid'
              ? '#5F6368'
              : status === 'assigned'
              ? '#49C998'
              : status === 'completed'
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

export const MyTable = ({ bounties, startDate, endDate }: TableProps) => {
  const [currentPage, setCurrentPage] = useState(1);
  const [totalBounties, setTotalBounties] = useState(0);
  const [activeTabs, setActiveTabs] = useState<number[]>([]);
  const pageSize = 20;
  const visibleTabs = 7;

  const { main } = useStores();

  const paginationLimit = Math.floor(totalBounties / pageSize) + 1;

  const currentPageData = () => {
    const indexOfLastPost = currentPage * pageSize;
    const indexOfFirstPost = indexOfLastPost - pageSize;
    const currentPosts = bounties.slice(indexOfFirstPost, indexOfLastPost);
    return currentPosts;
  };

  const paginateNext = () => {
    const activeTab = paginationLimit > visibleTabs;
    const activePage = currentPage < totalBounties / pageSize;
    if (activePage && activeTab) {
      const dataNumber: number[] = activeTabs;

      let nextPage: number;
      if (currentPage < visibleTabs) {
        nextPage = visibleTabs + 1;
        setCurrentPage(nextPage);
      } else {
        nextPage = currentPage + 1;
        setCurrentPage(nextPage);
      }

      dataNumber.push(nextPage);
      dataNumber.shift();
    }
  };
  const paginatePrev = () => {
    const firtsTab = activeTabs[0];
    const lastTab = activeTabs[6];
    if (firtsTab > 1) {
      const dataNumber: number[] = activeTabs;
      let nextPage: number;
      if (lastTab > visibleTabs) {
        nextPage = lastTab - visibleTabs;
      } else {
        nextPage = currentPage - 1;
      }

      setCurrentPage(currentPage - 1);
      dataNumber.pop();
      const newActivetabs = [nextPage, ...dataNumber];
      setActiveTabs(newActivetabs);
    }
  };

  const getTotalBounties = useCallback(async () => {
    if (startDate && endDate) {
      const totalBounties = await main.getBountiesCountByRange(String(startDate), String(endDate));
      setTotalBounties(totalBounties);
    }
  }, [main, startDate, endDate]);

  const getActiveTabs = useCallback(() => {
    const dataNumber: number[] = [];
    for (let i = 1; i <= Math.ceil(paginationLimit); i++) {
      if (i > visibleTabs) break;
      dataNumber.push(i);
    }
    setActiveTabs(dataNumber);
  }, []);

  useEffect(() => {
    getTotalBounties();
  }, [getTotalBounties]);

  useEffect(() => {
    getActiveTabs();
  }, [getActiveTabs]);

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
              <StyledSelect2 id="statusFilter">
                <option value="All">All</option>
                <option value="open">Open</option>
                <option value="in-progress">In Progress</option>
                <option value="completed">Completed</option>
              </StyledSelect2>
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
          {totalBounties > pageSize ? (
            <PageContainer>
              <img src={paginationarrow1} alt="" onClick={() => paginatePrev()} />
              {activeTabs.map((page: number) => (
                <PaginationButtons
                  key={page}
                  onClick={() => setCurrentPage(page)}
                  active={page === currentPage}
                >
                  {page}
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
