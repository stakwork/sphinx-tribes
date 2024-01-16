import React, { useCallback, useEffect, useState } from 'react';
//import styled from 'styled-components';
import { colors } from 'config';
import { EuiPopover, EuiText, EuiCheckboxGroup } from '@elastic/eui';
import { BountyStatus } from 'store/main';
import moment from 'moment';
import { useStores } from 'store';
import expand_more from '../header/icons/expand_more.svg';
import paginationarrow1 from '../header/icons/paginationarrow1.svg';
import paginationarrow2 from '../header/icons/paginationarrow2.svg';
import defaultPic from '../../../public/static/profile_avatar.svg';
import copygray from '../header/icons/copygray.svg';
import { Bounty } from './interfaces';
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
  BountyData,
  Paragraph,
  BoxImage
} from './TableStyle';

// import './styles.css';

import { FilterContainer, FlexDivStatus, StatusCheckboxItem } from './StatusStyle';

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

export const MyTable = ({
  bounties,
  startDate,
  endDate,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  headerIsFrozen = false,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  bountyStatus,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  setBountyStatus,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  dropdownValue,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  setDropdownValue
}: TableProps) => {
  const [currentPage, setCurrentPage] = useState(1);
  const pageSize = 20;
  const visibleTabs = 7;
  const [totalBounties, setTotalBounties] = useState(0);
  const [activeTabs, setActiveTabs] = useState<number[]>([]);
  const [statusFilter, setStatusFilter] = useState<string[]>([]);
  const [isPopoverOpen, setIsPopoverOpen] = useState(false);
  const [checkboxIdToSelectedMap, setCheckboxIdToSelectedMap] = useState({
    open: false,
    assigned: false,
    completed: false,
    paid: false
  });
  const { main } = useStores();
  const color = colors['light'];
  const paginationLimit = Math.floor(totalBounties / pageSize) + 1;
  const options = [
    { id: 'Open', label: 'Open', value: 'Open' },
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
  };

  const currentPageData = () => {
    const firstIndex = (currentPage - 1) * pageSize;
    const lastIndex = firstIndex + pageSize;
    return statusFilterMap().slice(firstIndex, lastIndex);
  };

  const onButtonClick = () => {
    setIsPopoverOpen(!isPopoverOpen);
  };

  const onChange = (optionId: any) => {
    const newCheckboxIdToSelectedMap = {
      ...checkboxIdToSelectedMap,
      ...{
        [optionId]: !checkboxIdToSelectedMap[optionId]
      }
    };

    setCheckboxIdToSelectedMap(newCheckboxIdToSelectedMap);

    // Check if the checkbox is being selected or deselected
    if (newCheckboxIdToSelectedMap[optionId]) {
      // If it's being selected, add the optionId to the statusFilter array
      setStatusFilter([...statusFilter, optionId]);
    } else {
      // If it's being deselected, remove the optionId from the statusFilter array
      setStatusFilter(statusFilter.filter((id: any) => id !== optionId));
    }
  };

  // const currentPageData = () => {
  //   const indexOfLastPost = currentPage * pageSize;
  //   const indexOfFirstPost = indexOfLastPost - pageSize;
  //   if (bounties) {
  //     const currentPosts = bounties.slice(indexOfFirstPost, indexOfLastPost);
  //     return currentPosts;
  //   }
  // };

  // const updateBountyStatus = (e: any) => {
  //   const { value } = e.target;
  //   if (bountyStatus && setBountyStatus && setDropdownValue) {
  //     switch (value) {
  //       case 'open': {
  //         const newStatus = { ...defaultBountyStatus, Open: true };
  //         setBountyStatus(newStatus);
  //         break;
  //       }
  //       case 'in-progress': {
  //         const newStatus = {
  //           ...defaultBountyStatus,
  //           Open: false,
  //           Assigned: true
  //         };
  //         setBountyStatus(newStatus);
  //         break;
  //       }
  //       case 'completed': {
  //         const newStatus = {
  //           ...defaultBountyStatus,
  //           Open: false,
  //           Paid: true
  //         };
  //         setBountyStatus(newStatus);
  //         break;
  //       }
  //       default: {
  //         const newStatus = {
  //           ...defaultBountyStatus,
  //           Open: false
  //         };
  //         setBountyStatus(newStatus);
  //         break;
  //       }
  //     }
  //     setDropdownValue(value);
  //   }
  // };

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
      <HeaderContainer>
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
              <EuiPopover
                button={
                  <FilterContainer onClick={onButtonClick}>
                    <FlexDivStatus>
                      <EuiText className="statusText">Status:</EuiText>
                      <EuiText className="subStatusText">
                        {statusFilter?.length === 0
                          ? 'All'
                          : statusFilter?.length === 1
                          ? statusFilter
                          : 'Multiple'}
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
                  border: '1px solid #fff'
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
                    onChange={(id: any) => onChange(id)}
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
            {currentPageData()?.map((bounty: any) => {
              // eslint-disable-next-line no-unused-vars
              const bounty_status =
                bounty?.paid && bounty?.assignee
                  ? 'paid'
                  : bounty?.assignee && !bounty?.paid
                  ? 'assigned'
                  : 'open';

              const created = moment.unix(bounty?.bounty_created).format('YYYY-MM-DD');
              const time_to_pay = bounty?.paid_date
                ? moment(bounty?.paid_date).diff(created, 'days')
                : 0;

              return (
                <TableDataRow key={bounty?.id}>
                  <BountyData className="avg">
                    <a
                      style={{ textDecoration: 'inherit', color: 'inherit' }}
                      href={`/bounty/${bounty?.bounty_id}`}
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
                    <TextInColorBox status={bounty.status} />
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
