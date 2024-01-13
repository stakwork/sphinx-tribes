import React, { useState } from "react";
import styled from "styled-components";
import moment from "moment";
import Calendar from "./Calender";

interface Props {
  filterStartDate: (newDate: number) => void;
  filterEndDate: (newDate: number) => void;
}

const App = ({filterStartDate,filterEndDate}:Props) => {
  const [from, setFrom] = useState(false);
  const [to, setTo] = useState(false);
  const [startDate, setStartDate] = useState<Date>();
  const [endDate, setEndDate] = useState<Date>();

  const Section = styled.div`
    display: flex;
    justify-content: flex-end;
    align-items: flex-end;
    gap: 8px;
    align-self: stretch;
    margin-right: 35px;
    z-index: 999;
  `;

  const MainContainer = styled.div`
    position: absolute;
    z-index: 999;
    right: 0;
    top: 70px;
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
    gap: 20px;
    align-self: stretch;
    height: ${from || to ? "580px" : "180px"};
  `;

  const FormDiv = styled.div`
    display: flex;
    justify-content: space-between;
    align-self: stretch;
    align-items: center;
    gap: 25px;
  `;

  const FormInput = styled.input`
    display: flex;
    width: 113px;
    height: 40px;
    padding: 8px 16px;
    justify-content: center;
    align-items: center;
    gap: 6px;
    border-radius: 6px;
    outline: none;
    border: 1px solid var(--Input-Outline-1, #d0d5d8);
    background: var(--White, #fff);
    color: var(--Placeholder-Text, var(--Disabled-Icon-color, #b0b7bc));
    text-align: center;
    font-size: 14px;
    font-style: normal;
    font-weight: 500;
    line-height: 0px;
    &:focus {
      color: var(--Placeholder-Text, var(--Disabled-Icon-color, #5078f2));
      border: 1px solid var(--Input-Outline-1, #5078f2);

      &::placeholder {
        color: var(--Placeholder-Text, var(--Disabled-Icon-color, #5078f2));
      }
    }
  `;

  const Button = styled.button`
    background-color: transparent;

    width: 54px;
    height: 40px;
    padding: 1px;
    color: ${(props: any) => props.color};
    text-align: center;
    font-size: 14px;
    font-style: normal;
    font-weight: 500;
    line-height: 20px;
    letter-spacing: 0.1px;
    outline: none;
    border: none;
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
    if (parameter === "From") {
      setTo(false);
      setFrom(!from);
    }
    if (parameter === "To") {
      setFrom(false);
      setTo(!to);
    }
  };

  const handleClick = () => {
    console.log("called");
    if (startDate && endDate) {
      console.log(endDate)
      const start = moment(startDate);
      const unixStart = start.unix();

      console.log("Unix Timestamp:", unixStart);
   

      const end = moment(endDate);
      const unixEnd = end.unix();
      console.log("Unix end:", unixEnd);
      
      
      filterStartDate(unixStart);
      filterEndDate(unixEnd)
    }

  };
  

  return (
    <>
      <MainContainer>
        <HeaderDiv>
          <FlexDiv>
            <h1 style={{ color: "#3C3F41", fontSize: "18px", fontWeight: "500", paddingBottom: "15px" }}>Enter Dates</h1>
            <FormDiv>
              <Formator>
                <Para>From</Para>
                <FormInput value={startDate ? moment(startDate).format("DD/MM/YY") : ""} placeholder="dd/mm/yy" type="text" readOnly onFocus={() => handleInputFocus("From")} />
              </Formator>
              <Formator>
                <Para>To</Para>
                <FormInput value={endDate ? moment(endDate).format("DD/MM/YY") : ""} placeholder="dd/mm/yy" type="text" readOnly onFocus={() => handleInputFocus("To")} />
              </Formator>
            </FormDiv>
          </FlexDiv>
          {from && <Calendar value={startDate} onChange={setStartDate} />}
          {to && <Calendar value={endDate} onChange={setEndDate} />}
          <Section>
            <Button onClick={()=>handleClick()} color="#618AFF">Save</Button>
            <Button color="#8E969C">Cancel</Button>
          </Section>
        </HeaderDiv>
      </MainContainer>
    </>
  );
};

export default App;
