import { EuiText } from '@elastic/eui';
import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { colors } from '../colors';
import { LanguageObject } from '../people/utils/language_label_style';
import NameTag from '../people/utils/nameTag';

const BountyDescription = (props: any) => {
  const color = colors['light'];
  const [dataValue, setDataValue] = useState([]);
  const [replitLink, setReplitLink] = useState('');
  const [descriptionImage, setDescriptionImage] = useState('');
  // const [descriptionLoomVideo, setDescriptionLoomVideo] = useState(props?.loomEmbedUrl);

  useEffect(() => {
    if (props.description) {
      const found = props?.description.match(/(https?:\/\/.*\.(?:png|jpg|jpeg|gif))/);
      setReplitLink(
        props?.description.match(
          /https?:\/\/(www\.)?[replit]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)/
        )
      );
      setDescriptionImage(found && found.length > 0 && found[0]);
    }
  }, [props]);

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
          <div className="NameContainer">
            <NameTag {...props} iconSize={32} isPaid={props?.isPaid} />
          </div>
        </Header>
        <Description>
          <div
            className="DescriptionContainer"
            style={{
              width: descriptionImage ? '334px' : '481px'
            }}>
            <EuiText
              className="DescriptionTitle"
              style={{
                color: props.isPaid ? color.grayish.G50 : color.grayish.G10
              }}>
              {props.title.slice(0, descriptionImage ? 80 : 120)}
              {props.title.length > 80 ? '...' : ''}
            </EuiText>
          </div>
          {descriptionImage && (
            <div className="DescriptionImage">
              <img
                src={descriptionImage}
                alt={'desc'}
                style={{ objectFit: 'cover' }}
                height={'100%'}
                width={'100%'}
              />
            </div>
          )}

          {/* 
          
          // TODO : add loom video - unable to add because some not supported features of loom video player.

          {props?.loomEmbedUrl && (
            <div
              style={{
                height: '64px',
                width: '130px',
                marginLeft: '17px',
                marginRight: '18px',
                borderRadius: '4px',
                overflow: 'hidden'
              }}>
              <iframe
                src={props?.loomEmbedUrl + '?autoplay=1&mute=1&loop=1&controls=0'}
                frameBorder="0"
                style={{
                  width: '100%',
                  height: '100%',
                  borderRadius: '4px'
                }}
              />
            </div>
          )} */}
        </Description>
        <LanguageContainer>
          {replitLink && (
            <div onClick={() => window.open(replitLink[0])} style={{ display: 'flex' }}>
              <CodingLabels
                key={0}
                border={`1px solid ${color.grayish.G06}`}
                LabelColor={color.grayish.G300}
                background={color.pureWhite}
                color={color}>
                <img
                  style={{ marginRight: '5px' }}
                  src={'/static/replit.png'}
                  alt={'replit_image'}
                  height={'15px'}
                  width={'15px'}
                />
                <EuiText className="LanguageText">Replit</EuiText>
              </CodingLabels>
            </div>
          )}
          {dataValue &&
            dataValue?.length > 0 &&
            dataValue?.map((lang: any, index) => {
              return (
                <CodingLabels
                  key={index}
                  border={props.isPaid ? `1px solid ${color.grayish.G06}` : lang?.border}
                  LabelColor={props.isPaid ? color.grayish.G300 : lang?.color}
                  background={props.isPaid ? color.grayish.G800 : lang?.background}
                  color={color}>
                  <EuiText className="LanguageText">{lang?.label}</EuiText>
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
  LabelColor?: string;
  color?: any;
}

interface bounty_description_props {}
interface replit_image_props {}

const BountyDescriptionContainer = styled.div<bounty_description_props>`
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
  .NameContainer {
    display: flex;
    flex-direction: column;
  }
`;

const Description = styled.div`
  display: flex;
  flex-direction: row;
  align-item: center;
  justify-content: space-between;
  .DescriptionContainer {
    display: flex;
    min-height: 64px;
    align-items: center;
  }
  .DescriptionTitle {
    font-size: 17px;
    line-height: 20px;
    font-weight: 500;
    display: flex;
    align-items: center;
  }
  .DescriptionImage {
    height: 77px;
    width: 130px;
    margin-left: 17px;
    margin-right: 18px;
    border-radius: 4px;
    overflow: hidden;
    margin-top: -13px;
    border: 1px solid #d0d5d8;
  }
`;

const LanguageContainer = styled.div`
  display: flex;
  flex-wrap: wrap;
  width: 80%;
  margin-top: 10px;
`;

const CodingLabels = styled.div<codingLangProps>`
  padding: 0px 8px;
  border: ${(p) => (p.border ? p?.border : `1px solid ${p.color.pureBlack}`)};
  color: ${(p) => (p.LabelColor ? p?.LabelColor : `${p.color.pureBlack}`)};
  background: ${(p) => (p.background ? p?.background : `${p.color.pureWhite}`)};
  border-radius: 4px;
  overflow: hidden;
  max-height: 22.75px;
  min-height: 22.75px;
  display: flex;
  flex-direction: row;
  align-items: center;
  margin-right: 4px;
  .LanguageText {
    font-size: 13px;
    fontweight: 500;
    text-align: center;
    font-family: Barlow;
    line-height: 16px;
  }
`;
