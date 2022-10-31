import { EuiText } from '@elastic/eui';
import React from 'react';
import styled from 'styled-components';
import { formatPrice, satToUsd } from '../helpers';

const BountyPrice = (props) => {
  return (
    <>
      <PriceContainer
        style={{
          ...props.style
        }}>
        <div
          style={{
            display: 'flex',
            alignItems: 'center'
          }}>
          <div
            style={{
              width: '26px'
            }}>
            <EuiText
              style={{
                fontSize: '14px',
                fontWeight: '400'
              }}>
              $@
            </EuiText>
          </div>
          {props.priceMin ? (
            <div
              style={{
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                height: '33px',
                width: '104px',
                color: '#2F7460',
                background: '#e4f7f0',
                borderRadius: '2px'
              }}>
              <EuiText
                style={{
                  fontSize: '17px',
                  fontWeight: '700',
                  lineHeight: '20.4px'
                }}>
                {formatPrice(props?.priceMin)} ~ {formatPrice(props?.priceMax)}
              </EuiText>
              <EuiText
                style={{
                  fontSize: '12px',
                  fontWeight: '400',
                  marginLeft: '6px'
                }}>
                SAT
              </EuiText>
            </div>
          ) : (
            <div
              style={{
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                height: '33px',
                width: '104px',
                color: '#2F7460',
                background: '#e4f7f0',
                borderRadius: '2px'
              }}>
              <EuiText
                style={{
                  fontSize: '17px',
                  fontWeight: '700',
                  lineHeight: '20.4px'
                }}>
                {formatPrice(props?.price)}
              </EuiText>
              <EuiText
                style={{
                  fontSize: '12px',
                  fontWeight: '400',
                  marginLeft: '6px'
                }}>
                SAT
              </EuiText>
            </div>
          )}
        </div>
        <div
          style={{
            display: 'flex',
            flexDirection: 'row',
            alignItems: 'center',
            paddingLeft: '36px'
          }}>
          {props.priceMin ? (
            <EuiText
              style={{
                fontSize: '13px',
                fontWeight: '500'
              }}>
              {satToUsd(props?.priceMin)} ~ {satToUsd(props?.priceMax)} USD
            </EuiText>
          ) : (
            <EuiText
              style={{
                fontSize: '13px',
                fontWeight: '500'
              }}>
              {satToUsd(props?.price)} USD{' '}
            </EuiText>
          )}
        </div>
        {props.sessionLength && (
          <EuiText
            style={{
              fontSize: '13px',
              fontWeight: '700',
              color: '#3c3f41'
            }}>
            <span
              style={{
                fontWeight: '400',
                color: '#8E969C'
              }}>
              Est:
            </span>{' '}
            &nbsp;
            {props.sessionLength}
          </EuiText>
        )}
      </PriceContainer>
    </>
  );
};

export default BountyPrice;

const PriceContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  padding: 0px 25px;
  color: #909baa;
  padding-top: 40px;
`;
