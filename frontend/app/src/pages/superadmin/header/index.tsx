import React, { useState, useRef, useEffect } from 'react';
import moment from 'moment';
import { mainStore } from 'store/main';
import {
  AlternateWrapper,
  ButtonWrapper,
  ExportButton,
  ExportText,
  Month,
  ArrowButton,
  DropDown,
  LeftWrapper,
  RightWrapper,
  Container,
  Option,
  CustomButton
} from './HeaderStyles';
import arrowback from './icons/arrowback.svg';
import arrowforward from './icons/arrowforward.svg';
import expand_more from './icons/expand_more.svg';
import App from './components/Calendar/App';
interface HeaderProps {
  startDate?: number;
  endDate?: number;
  setStartDate: (newDate: number) => void;
  setEndDate: (newDate: number) => void;
}
export const Header = ({ startDate, setStartDate, endDate, setEndDate }: HeaderProps) => {
  const [showSelector, setShowSelector] = useState(false);
  const [dateDiff, setDateDiff] = useState(7);
  const [exportLoading, setExportLoading] = useState(false);
  const [showCalendar,setShowCalendar] = useState(false);
  const formatUnixDate = (unixDate: number, includeYear: boolean = true) => {
    const formatString = includeYear ? 'DD-MMM-YYYY' : 'DD-MMM';
    return moment.unix(unixDate).format(formatString);
  };

  const handleBackClick = () => {
    if (startDate !== undefined && endDate !== undefined) {
      const newStartDate = moment.unix(startDate).subtract(dateDiff, 'days').unix();
      const newEndDate = moment.unix(endDate).subtract(dateDiff, 'days').unix();
      setStartDate(newStartDate);
      setEndDate(newEndDate);
    }
  };

  const handleForwardClick = () => {
    if (startDate !== undefined && endDate !== undefined) {
      const newStartDate = moment.unix(startDate).add(dateDiff, 'days').unix();
      const newEndDate = moment.unix(endDate).add(dateDiff, 'days').unix();

      const todayUnix = moment().startOf('day').unix();
      const cappedEndDate = Math.min(newEndDate, todayUnix);

      setStartDate(newStartDate);
      setEndDate(cappedEndDate);
    }
  };

  const exportCsv = async () => {
    setExportLoading(true);
    const csvUrl = await mainStore.exportMetricsBountiesCsv({
      start_date: String(startDate),
      end_date: String(endDate)
    });

    if (csvUrl) {
      window.open(csvUrl);
    }
    setExportLoading(false);
  };

  const handleDropDownChange = (option: number) => {
    const selectedValue = Number(option);
    setDateDiff(selectedValue);

    if (startDate && endDate) {
      const currentEndDate = moment.unix(endDate);
      let newStartDate;

      if (selectedValue === 7) {
        newStartDate = currentEndDate.clone().subtract(option, 'days').unix();
      } else if (selectedValue === 30) {
        newStartDate = currentEndDate.clone().subtract(option, 'days').unix();
      } else if (selectedValue === 90) {
        newStartDate = currentEndDate.clone().subtract(option, 'days').unix();
      }
      newStartDate = Math.max(
        newStartDate,
        currentEndDate.clone().subtract(selectedValue, 'days').unix()
      );

      setStartDate(newStartDate);
    }
  };

  const currentDateUnix = moment().unix();
  const optionRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    const handleOutsideClick = (event: MouseEvent) => {
      if (optionRef.current && !optionRef.current.contains(event.target as Node)) {
        setShowSelector(!showSelector);
      }
    };

    window.addEventListener('click', handleOutsideClick);

    return () => {
      window.removeEventListener('click', handleOutsideClick);
    };
  }, [showSelector]);

  return (
    <Container>
      <AlternateWrapper>
        <LeftWrapper data-testid="leftWrapper">
          {startDate && endDate ? (
            <>
              <ButtonWrapper>
                <ArrowButton onClick={() => handleBackClick()}>
                  <img src={arrowback} alt="" />
                </ArrowButton>
                <ArrowButton
                  disabled={
                    endDate === currentDateUnix ||
                    moment.unix(endDate).isSameOrAfter(moment().startOf('day'))
                  }
                  onClick={() => handleForwardClick()}
                >
                  <img src={arrowforward} alt="" />
                </ArrowButton>
              </ButtonWrapper>
              <Month data-testid="month">
                {formatUnixDate(startDate, false)} - {formatUnixDate(endDate)}
              </Month>
            </>
          ) : null}
        </LeftWrapper>
        <RightWrapper>
          <ExportButton disabled={exportLoading} onClick={() => exportCsv()}>
            <ExportText>{exportLoading ? 'Exporting ...' : 'Export CSV'}</ExportText>
          </ExportButton>
          <DropDown
            data-testid="DropDown"
            onClick={() => {
              setShowSelector(!showSelector);
            }}
          >
            Last {dateDiff} Days
            <div>
              <img src={expand_more} alt="a" />
            </div>
            {showSelector ? (
              <Option ref={optionRef}>
                <ul>
                  <li onClick={() => handleDropDownChange(7)}>7 Days</li>
                  <li onClick={() => handleDropDownChange(30)}>30 Days</li>
                  <li onClick={() => handleDropDownChange(90)}>90 Days</li>
                  <li>
                    <CustomButton onClick={()=>setShowCalendar(!showCalendar)}>Custom</CustomButton>
                  </li>
                </ul>
              </Option>
            ) : null}
          </DropDown>
        </RightWrapper>
      </AlternateWrapper>
      {showCalendar &&<App filterStartDate={setStartDate} filterEndDate={setEndDate}/>}
    </Container>
  );
};
