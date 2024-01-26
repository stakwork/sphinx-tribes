import React, { useState } from 'react';
import moment from 'moment';
import { EuiPopover, EuiText } from '@elastic/eui';
import MaterialIcon from '@material/react-material-icon';
import { BountyStatus } from '../../../store/main';
import paginationarrow1 from '../header/icons/paginationarrow1.svg';
import paginationarrow2 from '../header/icons/paginationarrow2.svg';
import defaultPic from '../../../public/static/profile_avatar.svg';
import copygray from '../header/icons/copygray.svg';
import { dateFilterOptions, getBountyStatus } from '../utils';
import { pageSize, visibleTabs } from '../constants.ts';
import { colors } from './../../../config/colors';
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
  BoxImage,
  DateFilterWrapper,
  DateFilterContent,
  PaginationImg
} from './TableStyle';

interface TableProps {
  bounties: Bounty[];
  startDate?: number;
  endDate?: number;
  headerIsFrozen?: boolean;
  bountyStatus?: BountyStatus;
  setBountyStatus?: React.Dispatch<React.SetStateAction<BountyStatus>>;
  dropdownValue?: string;
  sortOrder?: string;
  setDropdownValue?: React.Dispatch<React.SetStateAction<string>>;
  onChangeFilterByDate?: (option: string) => void;
  paginatePrev?: () => void;
  paginateNext?: () => void;
  currentPage: number;
  totalBounties: number;
  paginationLimit: number;
  setCurrentPage?: React.Dispatch<React.SetStateAction<number>>;
  activeTabs: number[];
  setActiveTabs: React.Dispatch<React.SetStateAction<number[]>>;
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
  bountyStatus,
  setBountyStatus,
  dropdownValue,
  headerIsFrozen,
  sortOrder,
  setDropdownValue,
  onChangeFilterByDate,
  currentPage,
  setCurrentPage,
  activeTabs,
  setActiveTabs,
  totalBounties,
  paginationLimit
}: TableProps) => {
  const [isPopoverOpen, setIsPopoverOpen] = useState<boolean>(false);
  const onButtonClick = () => setIsPopoverOpen((isPopoverOpen: any) => !isPopoverOpen);
  const closePopover = () => setIsPopoverOpen(false);

  const updateBountyStatus = (e: any) => {
    if (bountyStatus && setBountyStatus && setDropdownValue) {
      const { value } = e.target;
      getBountyStatus(value);
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
        if (setCurrentPage) setCurrentPage(nextPage);
      } else {
        nextPage = currentPage + 1;
        if (setCurrentPage) setCurrentPage(nextPage);
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

      if (setCurrentPage) setCurrentPage(currentPage - 1);
      dataNumber.pop();
      const newActivetabs = [nextPage, ...dataNumber];
      setActiveTabs(newActivetabs);
    }
  };

  const paginate = (page: number) => {
    if (setCurrentPage) {
      setCurrentPage(page);
    }
  };

  const color = colors['light'];

  return (
    <>
      <HeaderContainer freeze={!headerIsFrozen}>
        <Header>
          <BountyHeader>
            <img src={copygray} alt="" width="16.508px" height="20px" />
            <LeadingTitle>
              {totalBounties}
              <div>
                <AlternativeTitle>{bounties.length === 1 ? 'Bounty' : 'Bounties'}</AlternativeTitle>
              </div>
            </LeadingTitle>
          </BountyHeader>
          <Options>
            <FlexDiv>
              <EuiPopover
                button={
                  <DateFilterWrapper onClick={onButtonClick} color={color}>
                    <EuiText
                      className="filterText"
                      style={{
                        color: isPopoverOpen ? color.grayish.G10 : ''
                      }}
                    >
                      Sort By:
                    </EuiText>
                    <div className="image">
                      <EuiText className="filterText">
                        {sortOrder === 'desc' ? 'Newest' : 'Oldest'}
                      </EuiText>
                      <MaterialIcon
                        className="materialIconImage"
                        icon="expand_more"
                        style={{
                          color: isPopoverOpen ? color.grayish.G10 : '',
                          fontWeight: 'bold'
                        }}
                      />
                    </div>
                  </DateFilterWrapper>
                }
                panelStyle={{
                  marginTop: '3px',
                  border: 'none',
                  left: '700px',
                  maxWidth: '106px',
                  boxShadow: `0px 1px 20px ${color.black90}`,
                  background: `${color.pureWhite}`,
                  borderRadius: '6px'
                }}
                isOpen={isPopoverOpen}
                closePopover={closePopover}
                panelPaddingSize="none"
                anchorPosition="downRight"
              >
                <DateFilterContent className="CheckboxOuter" color={color}>
                  {dateFilterOptions.map((val: { [key: string]: string }) => (
                    <Options
                      onClick={() => {
                        onChangeFilterByDate?.(val.value);
                      }}
                      key={val.id}
                    >
                      {val.label}
                    </Options>
                  ))}
                </DateFilterContent>
              </EuiPopover>
            </FlexDiv>
            <FlexDiv>
              <Label>Status:</Label>
              <StyledSelect2 id="statusFilter" value={dropdownValue} onChange={updateBountyStatus}>
                <option value="All">All</option>
                <option value="Open">Open</option>
                <option value="Assigned">In Progress</option>
                <option value="Paid">Completed</option>
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
            {bounties.map((bounty: any) => {
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
              <PaginationImg
                src={paginationarrow1}
                alt="pagination arrow 1"
                onClick={() => paginatePrev()}
              />
              {activeTabs.map((page: number) => (
                <PaginationButtons
                  key={page}
                  onClick={() => paginate(page)}
                  active={page === currentPage}
                >
                  {page}
                </PaginationButtons>
              ))}
              <PaginationImg
                src={paginationarrow2}
                alt="pagination arrow 2"
                onClick={() => paginateNext()}
              />
            </PageContainer>
          ) : null}
        </FlexDiv>
      </PaginatonSection>
    </>
  );
};
