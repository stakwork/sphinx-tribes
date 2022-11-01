import { EuiText } from '@elastic/eui';
import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import NameTag from '../people/utils/nameTag';

const LanguageObject = [
  {
    label: 'Lightning',
    border: '1px solid rgba(184, 37, 95, 0.1)',
    background: 'rgba(184, 37, 95, 0.1)',
    color: '#B8255F'
  },
  {
    label: 'Javascript',
    border: '1px solid rgba(219, 64, 53, 0.1)',
    background: 'rgba(219, 64, 53, 0.1)',
    color: '#DB4035'
  },
  {
    label: 'Typescript',
    border: '1px solid rgba(255, 153, 51, 0.1)',
    background: ' rgba(255, 153, 51, 0.1)',
    color: '#FF9933'
  },
  {
    label: 'Node',
    border: '1px solid rgba(255, 191, 59, 0.1)',
    background: 'rgba(255, 191, 59, 0.1)',
    color: '#FFBF3B'
  },
  {
    label: 'Golang',
    border: '1px solid rgba(175, 184, 59, 0.1)',
    background: 'rgba(175, 184, 59, 0.1)',
    color: '#AFB83B'
  },
  {
    label: 'Swift',
    border: '1px solid rgba(126, 204, 73, 0.1)',
    background: 'rgba(126, 204, 73, 0.1)',
    color: '#7ECC49'
  },
  {
    label: 'Kotlin',
    border: '1px solid rgba(41, 148, 56, 0.1)',
    background: 'rgba(41, 148, 56, 0.1)',
    color: '#299438'
  },
  {
    label: 'MySQL',
    border: '1px solid rgba(106, 204, 188, 0.1)',
    background: 'rgba(106, 204, 188, 0.1)',
    color: '#6ACCBC'
  },
  {
    label: 'PHP',
    border: '1px solid rgba(21, 143, 173, 0.1)',
    background: 'rgba(21, 143, 173, 0.1)',
    color: '#158FAD'
  },
  {
    label: 'R',
    border: '1px solid rgba(64, 115, 255, 0.1)',
    background: 'rgba(64, 115, 255, 0.1)',
    color: '#4073FF'
  },
  {
    label: 'C#',
    border: '1px solid rgba(136, 77, 255, 0.1)',
    background: 'rgba(136, 77, 255, 0.1)',
    color: '#884DFF'
  },
  {
    label: 'C++',
    border: '1px solid rgba(175, 56, 235, 0.1)',
    background: 'rgba(175, 56, 235, 0.1)',
    color: '#AF38EB'
  },
  {
    label: 'Java',
    border: '1px solid rgba(235, 150, 235, 0.1)',
    background: 'rgba(235, 150, 235, 0.1)',
    color: '#EB96EB'
  },
  {
    label: 'Rust',
    border: '1px solid rgba(224, 81, 148, 0.1)',
    background: 'rgba(224, 81, 148, 0.1)',
    color: '#E05194'
  },
  {
    label: 'No-code',
    border: '1px solid rgba(255, 141, 133, 0.1)',
    background: 'rgba(255, 141, 133, 0.1)',
    color: '#FF8D85'
  }
];

const BountyDescription = (props) => {
  const [dataValue, setDataValue] = useState([]);
  useEffect(() => {
    let res;
    if (props.codingLanguage.length > 0) {
      res = LanguageObject?.filter((value) => {
        return props.codingLanguage?.find((val) => {
          return val.label === value.label;
        });
      });
    }
    setDataValue(res);
  }, [props.codingLanguage]);

  return (
    <>
      <BountyDescriptionContainer style={{ ...props.style }}>
        <Header>
          <div
            style={{
              display: 'flex',
              flexDirection: 'column'
            }}>
            <NameTag {...props} iconSize={32} isPaid={props?.isPaid} />
          </div>
        </Header>
        <Description>
          <div
            style={{
              width: '481px',
              minHeight: '64px',
              display: 'flex',
              alignItems: 'center'
            }}>
            <EuiText
              style={{
                fontSize: '17px',
                lineHeight: '20px',
                fontWeight: '500',
                color: props.isPaid ? '#5F6368' : '#3C3F41',
                display: 'flex',
                alignItems: 'center'
              }}>
              {props.title}
            </EuiText>
          </div>
        </Description>
        <LanguageContainer>
          {dataValue &&
            dataValue?.length > 0 &&
            dataValue?.map((lang: any, index) => {
              return (
                <CodingLabels
                  key={index}
                  border={props.isPaid ? '#f0f2f2' : lang?.border}
                  color={props.isPaid ? '#B0B7BC' : lang?.color}
                  background={props.isPaid ? '#f7f8f8' : lang?.background}>
                  <EuiText
                    style={{
                      fontSize: '13px',
                      fontWeight: '500',
                      textAlign: 'center',
                      fontFamily: 'Barlow',
                      lineHeight: '16px'
                    }}>
                    {lang?.label}
                  </EuiText>
                </CodingLabels>
              );
            })}
        </LanguageContainer>
      </BountyDescriptionContainer>
    </>
  );
};

export default BountyDescription;

interface codingLangProps {
  background?: string;
  border?: string;
  color?: string;
}

const BountyDescriptionContainer = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;
  min-width: 519px;
  max-width: 519px;
  padding-left: 17px;
`;

const Header = styled.div`
  display: flex;
  flex-direction: row;
  align-item: center;
  height: 32px;
  margin-top: 16px;
`;

const Description = styled.div`
  display: flex;
  flex-direction: row;
  align-item: center;
`;

const LanguageContainer = styled.div`
  display: flex;
  flex-wrap: wrap;
  width: 100%;
  margin-top: 9px;
`;

const CodingLabels = styled.div<codingLangProps>`
  padding: 0px 8px;
  border: ${(p) => (p.border ? p?.border : '1px solid #000')};
  color: ${(p) => (p.color ? p?.color : '#000')};
  background: ${(p) => (p.background ? p?.background : '#fff')};
  border-radius: 4px;
  overflow: hidden;
  max-height: 22.75px;
  min-height: 22.75px;
  display: flex;
  flex-direction: row;
  align-items: center;
  margin-right: 4px;
`;
