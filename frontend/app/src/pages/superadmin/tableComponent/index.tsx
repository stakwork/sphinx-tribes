import React, { useCallback, useEffect, useState } from 'react';
import { useStores } from 'store';
import { BountyStatus, defaultBountyStatus } from 'store/main';
import moment from 'moment';
import paginationarrow1 from '../header/icons/paginationarrow1.svg';
import paginationarrow2 from '../header/icons/paginationarrow2.svg';
import defaultPic from '../../../public/static/profile_avatar.svg';
import copygray from '../header/icons/copygray.svg';
import { Bounty } from './interfaces.ts';

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
  BountyData,
  Paragraph,
  BoxImage
} from './TableStyle';

//import './styles.css';

interface TableProps {
  bounties: Bounty[];
  startDate?: number;
  endDate?: number;
  headerIsFrozen?: boolean;
  bountyStatus?: BountyStatus;
  setBountyStatus?: React.Dispatch<React.SetStateAction<BountyStatus>>;
  dropdownValue?: string;
  setDropdownValue?: React.Dispatch<React.SetStateAction<string>>;
  paginatePrev?: () => void;
  paginateNext?: () => void;
}

interface ImageWithTextProps {
  image?: string;
  text: string;
}

export const ImageWithText = ({ image, text }: ImageWithTextProps) => (
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
        data-testid="bounty-status"
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

export const MyTable = ({
  bounties,
  startDate,
  endDate,
  bountyStatus,
  setBountyStatus,
  dropdownValue,
  headerIsFrozen,
  setDropdownValue
}: TableProps) => {
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
    if (bounties) {
      const currentPosts = bounties.slice(indexOfFirstPost, indexOfLastPost);
      return currentPosts;
    }
  };

  const updateBountyStatus = (e: any) => {
    const { value } = e.target;
    if (bountyStatus && setBountyStatus && setDropdownValue) {
      switch (value) {
        case 'open': {
          const newStatus = { ...defaultBountyStatus, Open: true };
          setBountyStatus(newStatus);
          break;
        }
        case 'in-progress': {
          const newStatus = {
            ...defaultBountyStatus,
            Open: false,
            Assigned: true
          };
          setBountyStatus(newStatus);
          break;
        }
        case 'completed': {
          const newStatus = {
            ...defaultBountyStatus,
            Open: false,
            Paid: true
          };
          setBountyStatus(newStatus);
          break;
        }
        default: {
          const newStatus = {
            ...defaultBountyStatus,
            Open: false
          };
          setBountyStatus(newStatus);
          break;
        }
      }
      setDropdownValue(value);
    }
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
  }, [paginationLimit]);

  useEffect(() => {
    getTotalBounties();
  }, [getTotalBounties]);

  useEffect(() => {
    getActiveTabs();
  }, [getActiveTabs]);

  return (
    <>
      <HeaderContainer freeze={!headerIsFrozen}>
        <Header>
          <BountyHeader>
            <img src={copygray} alt="" width="16.508px" height="20px" />
            <LeadingTitle>
              {bounties.length}
              <div>
                <AlternativeTitle>{bounties.length === 1 ? 'Bounty' : 'Bounties'}</AlternativeTitle>
              </div>
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
              <StyledSelect2 id="statusFilter" value={dropdownValue} onChange={updateBountyStatus}>
                <option value="all">All</option>
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
          <TableRow freeze={!headerIsFrozen}>
            <TableHeaderData>Bounty</TableHeaderData>
            <TableHeaderData>Date</TableHeaderData>
            <TableHeaderDataCenter>#DTGP</TableHeaderDataCenter>
            <TableHeaderData>Assignee</TableHeaderData>
            <TableHeaderData>Provider</TableHeaderData>
            <TableHeaderDataAlternative>Organization</TableHeaderDataAlternative>
            <TableHeaderDataRight>Status</TableHeaderDataRight>
          </TableRow>
          <tbody>
            {currentPageData()?.map((bounty: any) => {
              const bounty_status =
                bounty?.paid && bounty.assignee
                  ? 'paid'
                  : bounty.assignee && !bounty.paid
                  ? 'assigned'
                  : 'open';

              const created = moment.unix(bounty.bounty_created).format('YYYY-MM-DD');
              const time_to_pay = bounty.paid_date
                ? moment(bounty.paid_date).diff(created, 'days')
                : 0;

              return (
                <TableDataRow key={bounty?.id}>
                  <BountyData className="avg">
                    <a
                      style={{ textDecoration: 'inherit', color: 'inherit' }}
                      href={`/bounty/${bounty.bounty_id}`}
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      {bounty?.title}
                    </a>
                  </BountyData>
                  <TableData>{created}</TableData>
                  <TableDataCenter>{time_to_pay}</TableDataCenter>
                  <TableDataAlternative>
                    <ImageWithText
                      text={bounty?.assignee}
                      image={bounty?.assignee_img || defaultPic}
                    />
                  </TableDataAlternative>
                  <TableDataAlternative className="address">
                    <ImageWithText
                      text={bounty?.owner_pubkey}
                      image={bounty?.providerImage || defaultPic}
                    />
                  </TableDataAlternative>
                  <TableData className="organization">
                    <ImageWithText
                      text={bounty?.organization}
                      image={bounty?.organization_img || defaultPic}
                    />
                  </TableData>
                  <TableData3>
                    <TextInColorBox status={bounty_status} />
                  </TableData3>
                </TableDataRow>
              );
            })}
          </tbody>
        </Table>
      </TableContainer>
      <PaginatonSection>
        <FlexDiv>
          {totalBounties > pageSize ? (

            <PageContainer role="pagination">

              <img src={paginationarrow1} alt="pagination arrow 1" onClick={() => paginatePrev()} />
              {activeTabs.map((page: number) => (
                <PaginationButtons
                  key={page}
                  onClick={() => setCurrentPage(page)}
                  active={page === currentPage}
                >
                  {page}
                </PaginationButtons>
              ))}
              <img src={paginationarrow2} alt="pagination arrow 2" onClick={() => paginateNext()} />
            </PageContainer>
          ) : null}
        </FlexDiv>
      </PaginatonSection>
    </>
  );
};
