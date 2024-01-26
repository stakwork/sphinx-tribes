import styled from 'styled-components';

interface PriceContainerProps {
  price_Text_Color?: string;
  priceBackground?: string;
  session_text_color?: string;
}
export const PriceOuterContainer = styled.div<PriceContainerProps>`
  display: flex;
  align-items: center;
  height: 33px;
  min-width: 104px;
  color: ${(p: any) => (p.price_Text_Color ? p.price_Text_Color : '')};
  background: ${(p: any) => (p.priceBackground ? p.priceBackground : '')};
  border-radius: 2px;
  .Price_inner_Container {
    min-height: 33px;
    min-width: 63px;
    display: flex;
    align-items: center;
    margin-left: 7px;
    white-space: nowrap;
  }
  .Price_Dynamic_Text {
    font-size: 17px;
    font-weight: 700;
    line-height: 20px;
    display: flex;
    align-items: center;
  }
  .Price_SAT_Container {
    height: 33px;
    width: 34px;
    display: flex;
    align-items: center;
    margin-top: 1px;
    .Price_SAT_Text {
      font-size: 12px;
      font-weight: 400;
      margin-left: 6px;
    }
  }
`;
