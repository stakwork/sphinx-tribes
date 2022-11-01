import { EuiText } from '@elastic/eui';
import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { formatPrice, satToUsd } from '../helpers';

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
  const [session, setSession] = useState<any>();
  let dollarUSLocale = Intl.NumberFormat('en-US');

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
        }}>
        <div
          style={{
            display: 'flex',
            alignItems: 'center'
          }}>
          <div
            style={{
              width: '28px',
              height: '33px',
              display: 'flex',
              alignItems: 'center'
            }}>
            <EuiText
              style={{
                fontFamily: 'Barlow',
                fontStyle: 'normal',
                fontWeight: '400',
                fontSize: '14px',
                lineHeight: '17px'
              }}>
              $@
            </EuiText>
          </div>
          {props.priceMin ? (
            <div
              style={{
                display: 'flex',
                alignItems: 'center',
                height: '33px',
                minWidth: '104px',
                color: '#2F7460',
                background: 'rgba(73, 201, 152, 0.15)',
                borderRadius: '2px'
              }}>
              <div
                style={{
                  minHeight: '33px',
                  minWidth: '63px',
                  display: 'flex',
                  alignItems: 'center',
                  marginLeft: '7px'
                }}>
                <EuiText
                  style={{
                    fontSize: '17px',
                    fontWeight: '700',
                    lineHeight: '20px',
                    display: 'flex',
                    alignItems: 'center'
                  }}>
                  {dollarUSLocale.format(formatPrice(props?.priceMin)).split(',').join(' ')}
                  {/* ~{' '}
                  {dollarUSLocale.format(formatPrice(props?.priceMax)).split(',').join(' ')} */}
                </EuiText>
              </div>
              <div
                style={{
                  height: '33px',
                  width: '34px',
                  display: 'flex',
                  alignItems: 'center',
                  marginLeft: '3px'
                }}>
                <EuiText
                  style={{
                    fontSize: '12px',
                    fontWeight: '400',
                    marginLeft: '6px'
                  }}>
                  SAT
                </EuiText>
              </div>
            </div>
          ) : (
            <div
              style={{
                display: 'flex',
                alignItems: 'center',
                height: '33px',
                minWidth: '104px',
                color: '#2F7460',
                background: 'rgba(73, 201, 152, 0.15)',
                borderRadius: '2px'
              }}>
              <div
                style={{
                  height: '33px',
                  minWidth: '63px',
                  display: 'flex',
                  alignItems: 'center',
                  marginLeft: '7px'
                }}>
                <EuiText
                  style={{
                    fontSize: '17px',
                    fontWeight: '700',
                    lineHeight: '20px'
                  }}>
                  {dollarUSLocale.format(formatPrice(props?.price)).split(',').join(' ')}
                </EuiText>
              </div>

              <div
                style={{
                  height: '33px',
                  width: '34px',
                  display: 'flex',
                  alignItems: 'center',
                  marginTop: '1px'
                }}>
                <EuiText
                  style={{
                    fontSize: '12px',
                    fontWeight: '400',
                    marginLeft: '6px',
                    lineHeight: '14px'
                  }}>
                  SAT
                </EuiText>
              </div>
            </div>
          )}
        </div>
        <div
          style={{
            display: 'flex',
            flexDirection: 'row',
            alignItems: 'center',
            paddingLeft: '34px'
          }}>
          {props.priceMin ? (
            <EuiText
              style={{
                fontSize: '13px',
                fontWeight: '500'
              }}>
              {satToUsd(props?.priceMin)}
              {' '}
              {/* ~ {satToUsd(props?.priceMax)} */}
              USD
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
        {session && (
          <div
            style={{
              height: '28px'
            }}>
            <EuiText
              style={{
                fontSize: '13px',
                fontWeight: '700',
                color: '#3c3f41',
                fontFamily: 'Barlow'
              }}>
              <span
                style={{
                  fontWeight: '400',
                  fontFamily: 'Barlow',
                  color: '#8E969C'
                }}>
                Est:
              </span>{' '}
              &nbsp;&nbsp;&nbsp;
              {session.value === 'Not Sure' ? (
                <span>{session.value}</span>
              ) : (
                <span>
                  <span
                    style={{
                      fontFamily: 'Roboto',
                      fontSize: '12px',
                      fontWeight: '400',
                      lineHeight: '14.06px'
                    }}>
                    {session.value.slice(0, 1)}
                  </span>
                  {session.value.slice(1)}
                </span>
              )}
            </EuiText>
          </div>
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
  padding: 0px 24px;
  color: #909baa;
  padding-top: 41px;
`;
