import React, { useState } from 'react';
import styled from 'styled-components';
import moment from 'moment';
import { set } from 'date-fns';
import Calendar from './Calender';

interface Props {
  filterStartDate: (newDate: number) => void;
  filterEndDate: (newDate: number) => void;
  setShowCalendar: (show: boolean) => void;
}

const App = ({ filterStartDate, filterEndDate, setShowCalendar }: Props) => {
  const [from, setFrom] = useState(false);
  const [to, setTo] = useState(false);
  const [startDate, setStartDate] = useState<Date>();
  const [endDate, setEndDate] = useState<Date>();
  const [formInputFocused, setFormInputFocused] = useState(false);
  const [formInput2Focused, setFormInput2Focused] = useState(false);

  const Section = styled.div`
    display: flex;
    justify-content: flex-end;
    align-items: flex-end;
    gap: 8px;
    align-self: stretch;
    margin-right: 35px;
    z-index: 999;
    margin-bottom: 8px;
  `;

  const MainContainer = styled.div`
    position: absolute;
    z-index: 999;
    right: -51px;
    top: 40px;
    display: flex;
    width: 375px;
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
    border-radius: 6px;
    background: #fff;
    margin: 100px;
    box-shadow: 0px 4px 20px 0px rgba(0, 0, 0, 0.25);
  `;

  const HeaderDiv = styled.div`
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    align-self: stretch;
    // height: ${from || to ? '580px' : '180px'};
    gap: 16px;
    width: 375px;
  `;

  const FormDiv = styled.div`
    display: flex;
    justify-content: space-between;
    align-self: stretch;
    align-items: center;
    gap: 25px;
  `;
  type FormProps = {
    focused: boolean;
  };
  const FormInput = styled.input<FormProps>`
    display: flex;
    width: 113px;
    height: 40px;
    padding: 8px 16px;
    justify-content: center;
    align-items: center;
    gap: 6px;
    border-radius: 6px;
    outline: none;
    background: var(--White, #fff);
    color: var(--Placeholder-Text, var(--Disabled-Icon-color, #b0b7bc));
    text-align: center;
    font-size: 14px;
    font-style: normal;
    font-weight: 500;
    line-height: 0px;
    color: ${(props: any) =>
      props.focused
        ? 'var(--Placeholder-Text, var(--Disabled-Icon-color, #5078f2))'
        : 'var(--Placeholder-Text, var(--Disabled-Icon-color, #b0b7bc))'};

    border: ${(props: any) => (props.focused ? '1px solid #5078f2' : '1px solid #b0b7bc')};
    &::placeholder {
      color: ${(props: any) =>
        props.focused
          ? 'var(--Placeholder-Text, var(--Disabled-Icon-color, #5078f2))'
          : 'var(--Placeholder-Text, var(--Disabled-Icon-color, #b0b7bc))'};
    }
  `;
  type BtnProps = {
    color: string;
  };
  const Button = styled.button<BtnProps>`
    background-color: transparent;
    width: 54px;
    height: 40px;
    color: ${(props: any) => props.color};
    text-align: center;
    font-size: 14px;
    font-style: normal;
    font-weight: 500;
    line-height: 20px;
    letter-spacing: 0.1px;
    outline: none;
    border: none;
    font-family: Barlow;
  `;

  const Para = styled.p`
    color: var(--Text-2, var(--Hover-Icon-Color, #3c3f41));
    font-size: 14px;
    font-style: normal;
    font-weight: 500;
    line-height: 18px;
  `;

  const FlexDiv = styled.div`
    padding: 28px 0px 0px 28px;
  `;

  const Formator = styled.div`
    display: flex;
    justify-content: center;
    align-items: baseline;
    gap: 8px;
  `;

  const handleInputFocus = (parameter: string) => {
    if (parameter === 'From') {
      setTo(false);
      setFrom(true);
    }
    if (parameter === 'To') {
      setFrom(false);
      setTo(!to);
    }
  };

  const handleClick = () => {
    if (startDate && endDate) {
      let startDt = startDate;
      let endDt = endDate;
      startDt = set(startDate, { hours: 0, minutes: 0, seconds: 0 });
      if (startDt > endDt) {
        const temp = startDt;
        startDt = endDt;
        endDt = temp;
      }
      const start = moment(startDt);
      const unixStart = start.unix();
      const end = moment(endDt);
      const unixEnd = end.unix();
      filterStartDate(unixStart);
      filterEndDate(unixEnd);
      setShowCalendar(false);
    }
  };

  return (
    <>
      <MainContainer>
        <HeaderDiv>
          <FlexDiv>
            <h1
              style={{
                color: '#3C3F41',
                fontSize: '18px',
                fontWeight: '500',
                paddingBottom: '15px'
              }}
            >
              Enter Dates
            </h1>
            <FormDiv>
              <Formator>
                <Para>From</Para>
                <FormInput
                  value={startDate ? moment(startDate).format('MM/DD/YY') : ''}
                  placeholder="MM/DD/YY"
                  type="text"
                  onFocus={() => {
                    handleInputFocus('From');
                    setFormInputFocused(true);
                    setFormInput2Focused(false);
                  }}
                  focused={formInputFocused}
                />
              </Formator>
              <Formator>
                <Para>To</Para>
                <FormInput
                  value={endDate ? moment(endDate).format('MM/DD/YY') : ''}
                  placeholder="MM/DD/YY"
                  type="text"
                  onFocus={() => {
                    handleInputFocus('To');
                    setFormInput2Focused(true);
                    setFormInputFocused(false);
                  }}
                  onChange={(e: any) => setEndDate(e.target.value)}
                  focused={formInput2Focused}
                />
              </Formator>
            </FormDiv>
          </FlexDiv>
          {from && <Calendar value={startDate} onChange={setStartDate} />}
          {to && <Calendar value={endDate} onChange={setEndDate} />}
          <Section>
            <Button onClick={() => handleClick()} color="#618AFF">
              Save
            </Button>
            <Button color="#8E969C" onClick={() => setShowCalendar(false)}>
              Cancel
            </Button>
          </Section>
        </HeaderDiv>
      </MainContainer>
    </>
  );
};

export default App;
