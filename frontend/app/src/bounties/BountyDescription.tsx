import { EuiText } from '@elastic/eui';
import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { isString } from 'lodash';
import { Link } from 'react-router-dom';
import { OrganizationText, OrganizationWrap } from '../people/utils/style';
import { colors } from '../config/colors';
import { LanguageObject } from '../people/utils/languageLabelStyle';
import NameTag from '../people/utils/NameTag';
import { BountiesDescriptionProps } from './interfaces';

interface codingLangProps {
  background?: string;
  border?: string;
  LabelColor?: string;
  color?: any;
}

interface bounty_description_props {
  isPaid?: any;
  color?: any;
}

const BountyDescriptionContainer = styled.div<bounty_description_props>`
  display: flex;
  flex-direction: column;
  height: 100%;
  min-width: 519px;
  max-width: 519px;
  padding-left: 17px;
  padding-right: 17px;
`;

const Header = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  height: 32px;
  margin-top: 16px;
  .NameContainer {
    display: flex;
    flex-direction: column;
  }
`;

const Description = styled.div<bounty_description_props>`
  display: flex;
  flex-direction: row;
  align-items: center;
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
    border: 1px solid ${(p: any) => p?.color && p.color.grayish.G500};
    opacity: ${(p: any) => (p.isPaid ? 0.3 : 1)};
    filter: ${(p: any) => p.isPaid && 'grayscale(100%)'};
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
  border: ${(p: any) => (p.border ? p?.border : `1px solid ${p.color.pureBlack}`)};
  color: ${(p: any) => (p.LabelColor ? p?.LabelColor : `${p.color.pureBlack}`)};
  background: ${(p: any) => (p.background ? p?.background : `${p.color.pureWhite}`)};
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
    font-weight: 500;
    text-align: center;
    font-family: 'Barlow';
    line-height: 16px;
  }
`;

const Img = styled.div<{
  readonly src: string;
}>`
  background-image: url('${(p: any) => p.src}');
  background-position: center;
  background-size: cover;
  width: 20px;
  height: 20px;
  border-radius: 50%;
`;

const BountyDescription = (props: BountiesDescriptionProps) => {
  const color = colors['light'];
  const [dataValue, setDataValue] = useState([]);
  const [replitLink, setReplitLink] = useState('');
  const [descriptionImage, setDescriptionImage] = useState('');

  useEffect(() => {
    if (props.description) {
      const found = props?.description.match(/(https?:\/\/.*\.(?:png|jpg|jpeg|gif))(?![^`]*`)/);
      setReplitLink(
        props?.description.match(
          /https?:\/\/(?:www\.)?(?:replit\.[a-zA-Z0-9()]{1,256}|replit\.it)\b([-a-zA-Z0-9()@:%_+.~#?&//=]*)/
        )
      );
      setDescriptionImage(found && found.length > 0 && found[0]);
    }
  }, [props]);

  useEffect(() => {
    let res;
    if (props.codingLanguage.length > 0) {
      res = LanguageObject?.filter((value: any) =>
        !isString(props.codingLanguage)
          ? props.codingLanguage?.find((val: any) => val.label === value.label)
          : props.codingLanguage
      );
    }
    setDataValue(res);
  }, [props.codingLanguage]);

  return (
    <>
      <BountyDescriptionContainer style={{ ...props.style }}>
        <Header>
          <div className="NameContainer">
            <NameTag
              {...props}
              iconSize={32}
              owner_pubkey={props.owner_pubkey}
              img={props.img}
              id={props.id}
              widget={props.widget}
              owner_alias={props.owner_alias}
              isPaid={props?.isPaid}
              org_img={props.org_img}
              org_name={props.name}
              org_uuid={props.uuid}
            />
          </div>
          {props.org_uuid && props.name && (
            <Link
              onClick={(e: any) => {
                e.stopPropagation();
              }}
              to={`/org/bounties/${props.org_uuid}`}
              target="_blank"
            >
              <OrganizationWrap>
                <Img
                  title={`${props.name} logo`}
                  src={props.org_img || '/static/person_placeholder.png'}
                />
                <OrganizationText>{props.name}</OrganizationText>
                <img
                  className="buttonImage"
                  src={'/static/github_ticket.svg'}
                  alt={'github_ticket'}
                  height={'10px'}
                  width={'10px'}
                  style={{ transform: 'translateY(1px)' }}
                />
              </OrganizationWrap>
            </Link>
          )}
        </Header>
        <Description isPaid={props?.isPaid} color={color}>
          <div
            className="DescriptionContainer"
            style={{
              width: descriptionImage ? '334px' : '481px'
            }}
          >
            <EuiText
              className="DescriptionTitle"
              style={{
                color: props.isPaid ? color.grayish.G50 : color.grayish.G10
              }}
            >
              {props.title?.slice(0, descriptionImage ? 80 : 120)}
              {props.title?.length > 80 ? '...' : ''}
            </EuiText>
          </div>
        </Description>
        <LanguageContainer>
          {replitLink && (
            <div onClick={() => window.open(replitLink[0])} style={{ display: 'flex' }}>
              <CodingLabels
                key={0}
                border={`1px solid ${color.grayish.G06}`}
                LabelColor={color.grayish.G300}
                background={color.pureWhite}
                color={color}
              >
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
            dataValue?.map((lang: any, index: number) => (
              <CodingLabels
                key={index}
                border={props.isPaid ? `1px solid ${color.grayish.G06}` : lang?.border}
                LabelColor={props.isPaid ? color.grayish.G300 : lang?.color}
                background={props.isPaid ? color.grayish.G800 : lang?.background}
                color={color}
              >
                <EuiText className="LanguageText">{lang?.label}</EuiText>
              </CodingLabels>
            ))}
        </LanguageContainer>
      </BountyDescriptionContainer>
    </>
  );
};

export default BountyDescription;
