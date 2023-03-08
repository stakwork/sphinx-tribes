/* eslint-disable func-style */
import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { useIsMobile } from '../../hooks';
import { Title } from '../../components/common';
import { extractGithubIssue, extractGithubIssueFromUrl } from '../../helpers';
import { useStores } from '../../store';
import PaidBounty from '../utils/paidBounty';
import Bounties from '../utils/assigned_unassigned_bounties';
import { colors } from '../../config/colors';
import MobileView from "./wantedViews/mobileView";
import DesktopView from "./wantedViews/mobileView";

export default function WantedView(props: any) {
  const {
    one_sentence_summary,
    title,
    description,
    priceMin,
    priceMax,
    price,
    person,
    created,
    issue,
    ticketUrl,
    repo,
    type,
    codingLanguage,
    assignee,
    estimate_session_length,
    loomEmbedUrl,
    onPanelClick,
    key
  } = props;
  const titleString = title ?? one_sentence_summary;

  let { show, paid } = props;
  const isMobile = useIsMobile();
  const { ui, main } = useStores();
  const [saving, setSaving] = useState(false);
  const [labels, setLabels] = useState([]);
  const { peopleWanteds } = main;
  const color = colors['light'];

  const isMine = ui.meInfo?.owner_pubkey === person?.owner_pubkey;

  if ('show' in props) {
    // show has a value
  } else {
    // if no value default to true
    show = true;
  }

  if ('paid' in props) {
    // show has no value
  } else {
    // if no value default to false
    paid = false;
  }

  async function setExtrasPropertyAndSave(propertyName: string) {
    if (peopleWanteds) {
      setSaving(true);
      try {
        const targetProperty = props[propertyName];
        const [clonedEx, targetIndex] = await main.setExtrasPropertyAndSave(
          'wanted',
          propertyName,
          created,
          !targetProperty
        );

        // saved? ok update in wanted list if found
        const peopleWantedsClone: any = [...peopleWanteds];
        const indexFromPeopleWanted = peopleWantedsClone.findIndex((f) => {
          const val = f.body || {};
          return f.person.owner_pubkey === ui.meInfo?.owner_pubkey && val.created === created;
        });

        // if we found it in the wanted list, update in people wanted list
        if (indexFromPeopleWanted > -1) {
          // if it should be hidden now, remove it from the list
          if ('show' in clonedEx[targetIndex] && clonedEx[targetIndex].show === false) {
            peopleWantedsClone.splice(indexFromPeopleWanted, 1);
          } else {
            peopleWantedsClone[indexFromPeopleWanted] = {
              person: person,
              body: clonedEx[targetIndex]
            };
          }
          main.setPeopleWanteds(peopleWantedsClone);
        }
      } catch (e) {
        console.log('e', e);
      }

      setSaving(false);
    }
  }

  useEffect(() => {
    if (codingLanguage) {
      const values = codingLanguage.map((value) => ({ ...value }));
      setLabels(values);
    }
  }, [codingLanguage]);

  function renderTickets() {
    const { status } = ticketUrl
      ? extractGithubIssueFromUrl(person, ticketUrl)
      : extractGithubIssue(person, repo, issue);

    const isClosed = status === 'closed' || paid ? true : false;

    const isCodingTask =
      type === 'coding_task' || type === 'wanted_coding_task' || type === 'freelance_job_request';

    // mobile view
    if (isMobile) {
      return (
        <MobileView
          {...props} labels={labels}
          saving={saving}
          setExtrasPropertyAndSave={setExtrasPropertyAndSave}
          isClosed={isClosed}
          isCodingTask={isCodingTask}
          status={status}
          show={show}
          paid={paid}
          isMine={isMine}
          titleString={titleString}
        />
      )
    }

    if (props?.fromBountyPage) {
      return (
        <div key={key}>
          {paid ? (
            <BountyBox color={color}>
              <PaidBounty
                {...person}
                onPanelClick={onPanelClick}
                assignee={assignee}
                created={created}
                ticketUrl={ticketUrl}
                loomEmbedUrl={loomEmbedUrl}
                title={titleString}
                codingLanguage={labels}
                priceMin={priceMin}
                priceMax={priceMax}
                price={price}
                sessionLength={estimate_session_length}
                description={description}
              />
            </BountyBox>
          ) : (
            <BountyBox color={color}>
              <Bounties
                onPanelClick={onPanelClick}
                person={person}
                assignee={assignee}
                created={created}
                ticketUrl={ticketUrl}
                loomEmbedUrl={loomEmbedUrl}
                title={titleString}
                codingLanguage={labels}
                priceMin={priceMin}
                priceMax={priceMax}
                price={price}
                sessionLength={estimate_session_length}
                description={description}
              />
            </BountyBox>
          )}
        </div>
      );
    }

    return (
      <DesktopView
        {...props} labels={labels}
        saving={saving}
        setExtrasPropertyAndSave={setExtrasPropertyAndSave}
        isClosed={isClosed}
        isCodingTask={isCodingTask}
        status={status}
        show={show}
        paid={paid}
        isMine={isMine}
        titleString={titleString}
      />
    );
  }

  return renderTickets();
}

interface WrapProps {
  isClosed?: boolean;
  color?: any;
}

interface styledProps {
  color?: any;
}

const BountyBox = styled.div<styledProps>`
  min-height: 160px;
  max-height: 160px;
  width: 1100px;
  box-shadow: 0px 1px 6px ${(p) => p?.color && p?.color.black100};
  border: none;
`;

const DWrap = styled.div<WrapProps>`
  display: flex;
  flex: 1;
  height: 100%;
  min-height: 510px;
  flex-direction: column;
  width: 100%;
  min-width: 100%;
  max-height: 510px;
  font-style: normal;
  font-weight: 500;
  font-size: 17px;
  line-height: 23px;
  color: ${(p) => p?.color && p?.color.grayish.G10} !important;
  letter-spacing: 0px;
  justify-content: space-between;
  opacity: ${(p) => (p.isClosed ? '0.5' : '1')};
  filter: ${(p) => (p.isClosed ? 'grayscale(1)' : 'grayscale(0)')};
`;

const Wrap = styled.div<WrapProps>`
  display: flex;
  justify-content: flex-start;
  opacity: ${(p) => (p.isClosed ? '0.5' : '1')};
  filter: ${(p) => (p.isClosed ? 'grayscale(1)' : 'grayscale(0)')};
`;

const B = styled.span<styledProps>`
  font-size: 14px;
  font-weight: bold;
  color: ${(p) => p?.color && p?.color.grayish.G10};
`;
const P = styled.div<styledProps>`
  font-weight: regular;
  font-size: 14px;
  color: ${(p) => p?.color && p?.color.grayish.G100};
`;

const Body = styled.div<styledProps>`
  font-size: 15px;
  line-height: 20px;
  padding: 10px;
  display: flex;
  flex-direction: column;
  justify-content: space-around;
  color: ${(p) => p?.color && p?.color.grayish.G05};
  overflow: hidden;
  min-height: 132px;
`;

const Pad = styled.div`
  display: flex;
  flex-direction: column;
`;

const DescriptionCodeTask = styled.div<styledProps>`
  margin-bottom: 10px;

  font-family: Roboto;
  font-style: normal;
  font-weight: normal;
  font-size: 13px;
  line-height: 20px;
  color: ${(p) => p?.color && p?.color.grayish.G50};
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 6;
  -webkit-box-orient: vertical;
  height: 120px;
  max-height: 120px;
`;
const DT = styled(Title)`
  margin-bottom: 9px;
  max-height: 52px;
  min-height: 43.5px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  /* Primary Text 1 */

  font-family: 'Roboto';
  font-style: normal;
  font-weight: 500;
  font-size: 17px;
  line-height: 23px;
`;

interface ImageProps {
  readonly src?: string;
}
const Img = styled.div<ImageProps>`
  background-image: url('${(p) => p.src}');
  background-position: center;
  background-size: cover;
  position: relative;
  width: 22px;
  height: 22px;
`;

const EyeDeleteTextContainerMobile = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
`;

const EyeDeleteContainerMobile = styled.div`
  margin-top: 10px;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
`;
