import { EuiText } from '@elastic/eui';
import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { colors } from '../colors';
import { DollarConverter, satToUsd } from '../helpers';

const Session = [
  {
    label: 'Less than 1 hour',
    value: '< 1 hrs'
  },
  {
    label: 'Less than 3 hours',
    value: '< 3 hrs'
  },
  {
    label: 'More than 3 hours',
    value: '> 3 hrs'
  },
  {
    label: 'Not sure yet',
    value: 'Not Sure'
  }
];

const BountyPrice = (props) => {
  const color = colors['light'];
  const [session, setSession] = useState<any>();

  useEffect(() => {
    let res;
    if (props.sessionLength) {
      res = Session?.find((value: any) => {
        return props?.sessionLength === value.label;
      });
    }
    setSession(res);
  }, [props]);

  return (
    <>
      <PriceContainer
        style={{
          ...props.style
        }}
      >
        <div
          style={{
            display: 'flex',
            alignItems: 'center'
          }}
        >
          <div className="PriceStaticTextContainer">
            <EuiText className="PriceStaticText">$@</EuiText>
          </div>
          {props.priceMin ? (
            <PriceOuterContainer
              price_Text_Color={color.primaryColor.P300}
              priceBackground={color.primaryColor.P100}
            >
              <div className="Price_inner_Container">
                <EuiText className="Price_Dynamic_Text">{DollarConverter(props?.priceMin)}</EuiText>
              </div>
              <div className="Price_SAT_Container">
                <EuiText className="Price_SAT_Text">SAT</EuiText>
              </div>
            </PriceOuterContainer>
          ) : (
            <PriceOuterContainer
              price_Text_Color={color.primaryColor.P300}
              priceBackground={color.primaryColor.P100}
            >
              <div className="Price_inner_Container">
                <EuiText className="Price_Dynamic_Text">{DollarConverter(props?.price)}</EuiText>
              </div>

              <div className="Price_SAT_Container">
                <EuiText className="Price_SAT_Text">SAT</EuiText>
              </div>
            </PriceOuterContainer>
          )}
        </div>
        <USDContainer>
          {props.priceMin ? (
            <EuiText className="USD_Price">
              {satToUsd(props?.priceMin)}
              USD
            </EuiText>
          ) : (
            <EuiText className="USD_Price">{satToUsd(props?.price)} USD </EuiText>
          )}
        </USDContainer>
        {session && (
          <SessionContainer session_text_color={color.grayish.G10}>
            <EuiText className="Session_Text">
              <span
                className="EST_Text"
                style={{
                  color: color.grayish.G100
                }}
              >
                Est:
              </span>{' '}
              &nbsp;&nbsp;&nbsp;
              {session.value === 'Not Sure' ? (
                <span>{session.value}</span>
              ) : (
                <span>
                  <span className="EST_Value">{session.value.slice(0, 1)}</span>
                  {session.value.slice(1)}
                </span>
              )}
            </EuiText>
          </SessionContainer>
        )}
      </PriceContainer>
    </>
  );
};

export default BountyPrice;

interface PriceContainerProps {
  price_Text_Color?: string;
  priceBackground?: string;
  session_text_color?: string;
}

const PriceContainer = styled.div<PriceContainerProps>`
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  padding: 0px 24px;
  color: #909baa;
  padding-top: 41px;
  .PriceStaticTextContainer {
    width: 28px;
    height: 33px;
    display: flex;
    align-items: center;
  }
  .PriceStaticText {
    font-family: Barlow;
    font-style: normal;
    font-weight: 400;
    font-size: 14px;
    line-height: 17px;
  }
`;

const PriceOuterContainer = styled.div<PriceContainerProps>`
  display: flex;
  align-items: center;
  height: 33px;
  min-width: 104px;
  color: ${(p) => (p.price_Text_Color ? p.price_Text_Color : '')};
  background: ${(p) => (p.priceBackground ? p.priceBackground : '')};
  border-radius: 2px;
  .Price_inner_Container {
    min-height: 33px;
    min-width: 63px;
    display: flex;
    align-items: center;
    margin-left: 7px;
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

const USDContainer = styled.div`
  display: flex;
  flex-direction: row;
  align-items: center;
  padding-left: 34px;
  .USD_Price {
    font-size: 13px;
    font-weight: 500;
  }
`;

const SessionContainer = styled.div<PriceContainerProps>`
  height: 28px;
  .Session_Text {
    font-size: 13px;
    font-weight: 700;
    color: ${(p) => (p.session_text_color ? p.session_text_color : '')};
    font-family: Barlow;
  }
  .EST_Text {
    font-weight: 400;
    font-family: Barlow;
  }
  .EST_Value {
    font-family: Roboto;
    font-size: 12px;
    font-weight: 400;
    line-height: 14.06px;
  }
`;
