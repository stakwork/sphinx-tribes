/* eslint-disable func-style */
import React, { useCallback, useEffect, useState } from 'react';
import { useLocation } from 'react-router-dom';
import { observer } from 'mobx-react-lite';
import api from '../../../api';
import { colors } from '../../../config/colors';
import Form from '../../../components/form/bounty';
import { sendBadgeSchema } from '../../../components/form/schema';
import { useIsMobile } from '../../../hooks';
import { Button } from '../../../components/common';
import { useStores } from '../../../store';
import { LanguageObject, awards } from '../../utils/languageLabelStyle';
import NameTag from '../../utils/NameTag';
import { sendToRedirect } from '../../../helpers';
import { CodingLanguageLabel, WantedSummaryProps, LocalPaymeentState } from '../../interfaces';
import CodingBounty from './wantedSummaries/CodingBounty';
import CodingDesktop from './wantedSummaries/CodingDesktop';
import { ButtonRow, Img, Assignee } from './wantedSummaries/style';
import { paidString, unpaidString } from './constants';

function useQuery() {
  const { search } = useLocation();
  return React.useMemo(() => new URLSearchParams(search), [search]);
}

function WantedSummary(props: WantedSummaryProps) {
  const {
    description,
    priceMin,
    ticket_url,
    person,
    created,
    repo,
    issue,
    price,
    type,
    tribe,
    paid,
    badgeRecipient,
    loomEmbedUrl,
    coding_languages,
    estimated_session_length,
    assignee,
    fromBountyPage,
    wanted_type,
    one_sentence_summary,
    github_description,
    show,
    setIsModalSideButton,
    setIsExtraStyle,
    formSubmit,
    title,
    org_uuid,
    id
  } = props;
  const titleString = one_sentence_summary || title || '';
  const bountyPath = `/bounty/${id}`;

  const isMobile = useIsMobile();
  const { main, ui } = useStores();
  const { peopleBounties } = main;
  const color = colors['light'];

  const [assigneeInfo, setAssigneeInfo]: any = useState(null);
  const [saving, setSaving]: any = useState('');
  const [isCopied, setIsCopied] = useState(false);
  const [owner_idURL, setOwnerIdURL] = useState('');
  const [createdURL, setCreatedURL] = useState('');
  const [dataValue, setDataValue] = useState([]);
  const [peopleList, setPeopleList] = useState<any>();
  const [isAssigned, setIsAssigned] = useState<boolean>(false);
  const [assignedPerson, setAssignedPerson] = useState<any>();
  const [replitLink, setReplitLink] = useState('');
  const [creatorStep, setCreatorStep] = useState<number>(0);
  const [bountyPrice, setBountyPrice] = useState<any>(price ?? priceMin ?? 0);
  const [selectedAward, setSelectedAward] = useState('');
  const [isPaidStatusPopOver, setIsPaidStatusPopOver] = useState<boolean>(false);
  const [isPaidStatusBadgeInfo, setIsPaidStatusBadgeInfo] = useState<boolean>(false);
  const [isMarkPaidSaved, setIsMarkPaidSaved] = useState<boolean>(false);
  const [awardDetails, setAwardDetails] = useState<any>({
    name: '',
    image: ''
  });

  useEffect(() => {
    if (description) {
      setReplitLink(
        description.match(
          /https?:\/\/(?:www\.)?(?:replit\.[a-zA-Z0-9()]{1,256}|replit\.it)\b([-a-zA-Z0-9()@:%_+.~#?&//=]*)/
        )
      );
    }
  }, [description]);

  useEffect(() => {
    const timer = setTimeout(() => {
      setIsPaidStatusPopOver(false);
    }, 7000);

    return () => {
      clearTimeout(timer);
    };
  }, [isPaidStatusPopOver]);

  const handleAwards = (optionId: any) => {
    setSelectedAward(optionId);
  };

  const [showBadgeAwardDialog, setShowBadgeAwardDialog] = useState(false);

  const isMine = ui.meInfo?.owner_pubkey === person?.owner_pubkey;

  const [labels, setLabels] = useState<Array<CodingLanguageLabel>>([]);
  const [assigneeValue, setAssigneeValue] = useState(false);
  const [localPaid, setLocalPaid] = useState<LocalPaymeentState>('UNKNOWN');

  const assigneeHandlerOpen = () => setAssigneeValue((assigneeValue: any) => !assigneeValue);

  useEffect(() => {
    if (assignee?.owner_alias) {
      setIsAssigned(true);
    }
  }, [assignee]);

  useEffect(() => {
    (async () => {
      try {
        const response = await api.get(`people?page=1&search=&sortBy=last_login&limit=100`);
        setPeopleList(response);
      } catch (error) {
        console.log(error);
      }
    })();
  }, []);

  const handleAssigneeDetails = useCallback(
    (value: any) => {
      setIsAssigned(true);
      setAssignedPerson(value);
      assigneeHandlerOpen();

      const newValue = {
        id,
        title: titleString,
        wanted_type: wanted_type,
        one_sentence_summary: one_sentence_summary,
        ticketUrl: ticket_url,
        github_description: github_description,
        description: description,
        price: price,
        assignee: value?.owner_pubkey || '',
        coding_language: coding_languages?.map((x: string) => ({ label: x, value: x })),
        estimated_session_length: estimated_session_length,
        show: show,
        type: type,
        created: created,
        org_uuid
      };

      formSubmit && formSubmit(newValue, true);
    },
    [
      coding_languages,
      created,
      description,
      estimated_session_length,
      formSubmit,
      github_description,
      one_sentence_summary,
      price,
      show,
      ticket_url,
      titleString,
      type,
      wanted_type
    ]
  );

  const changeAssignedPerson = useCallback(() => {
    setIsAssigned(false);

    setAssignedPerson(null);
    const newValue = {
      id,
      title: titleString,
      wanted_type: wanted_type,
      one_sentence_summary: one_sentence_summary,
      ticketUrl: ticket_url,
      github_description: github_description,
      description: description,
      price: price,
      assignee: '',
      coding_language: coding_languages?.map((x: string) => ({ label: x, value: x })),
      estimated_session_length: estimated_session_length,
      show: show,
      type: type,
      created: created,
      org_uuid
    };
    formSubmit && formSubmit(newValue, true);
  }, [
    coding_languages,
    created,
    description,
    estimated_session_length,
    formSubmit,
    github_description,
    one_sentence_summary,
    price,
    show,
    ticket_url,
    titleString,
    type,
    wanted_type
  ]);

  useEffect(() => {
    (async () => {
      if (props.assignee) {
        try {
          setAssigneeInfo(props.assignee);
        } catch (e) {
          console.log('e', e);
        }
      }
    })();
  }, [main, props.assignee, tribe]);

  useEffect(() => {
    let res;
    if (coding_languages?.length > 0) {
      res = LanguageObject?.filter(
        (value: any) => coding_languages?.find((val: string) => val === value.label)
      );
    }
    setDataValue(res);
  }, [coding_languages]);

  const searchParams = useQuery();

  useEffect(() => {
    const owner_id = searchParams.get('owner_id');
    const created = searchParams.get('created');
    setOwnerIdURL(owner_id ?? '');
    setCreatedURL(created ?? '');
  }, [owner_idURL, createdURL, searchParams]);

  useEffect(() => {
    if (coding_languages) {
      const values = coding_languages.map((value: string) => ({ label: value, value: value }));
      setLabels(values);
    }
  }, [coding_languages]);

  async function setExtrasPropertyAndSave(propertyName: string, value: any) {
    if (peopleBounties) {
      setSaving(propertyName);
      try {
        const [clonedEx, targetIndex] = await main.setExtrasPropertyAndSave(
          'wanted',
          propertyName,
          created ?? 0,
          value
        );

        // saved? ok update in wanted list if found
        const peopleBountiesClone: any = [...peopleBounties];
        const indexFromPeopleBounty = peopleBountiesClone.findIndex((f: any) => {
          const val = f.body || {};
          return f.person.owner_pubkey === ui.meInfo?.owner_pubkey && val.created === created;
        });

        // if we found it in the wanted list, update in people wanted list
        if (indexFromPeopleBounty > -1) {
          // if it should be hidden now, remove it from the list
          if ('show' in clonedEx[targetIndex] && clonedEx[targetIndex].show === false) {
            peopleBountiesClone.splice(indexFromPeopleBounty, 1);
          } else {
            // gotta update person extras! this is what is used for summary viewer
            const personClone: any = person;
            personClone.extras['wanted'][targetIndex] = clonedEx[targetIndex];

            peopleBountiesClone[indexFromPeopleBounty] = {
              person: personClone,
              body: clonedEx[targetIndex]
            };
          }
          main.setPeopleBounties(peopleBountiesClone);
        }
      } catch (e) {
        console.log('e', e);
      }

      setSaving('');
    }
  }

  async function setExtrasPropertyAndSaveMultiple(propertyName: any, dataObject: any) {
    if (peopleBounties) {
      setIsMarkPaidSaved(true);
      try {
        const [clonedEx, targetIndex] = await main.setExtrasMultipleProperty(
          dataObject,
          'wanted',
          created ?? 0
        );

        // saved? ok update in wanted list if found
        const peopleBountiesClone: any = [...peopleBounties];
        const indexFromPeopleBounty = peopleBountiesClone.findIndex((f: any) => {
          const val = f.body || {};
          return f.person.owner_pubkey === ui.meInfo?.owner_pubkey && val.created === created;
        });

        // if we found it in the wanted list, update in people wanted list
        if (indexFromPeopleBounty > -1) {
          // if it should be hidden now, remove it from the list
          if ('show' in clonedEx[targetIndex] && clonedEx[targetIndex].show === false) {
            peopleBountiesClone.splice(indexFromPeopleBounty, 1);
          } else {
            // gotta update person extras! this is what is used for summary viewer
            const personClone: any = person;
            personClone.extras['wanted'][targetIndex] = clonedEx[targetIndex];

            peopleBountiesClone[indexFromPeopleBounty] = {
              person: personClone,
              body: clonedEx[targetIndex]
            };
          }

          main.setPeopleBounties(peopleBountiesClone);
        }
      } catch (e) {
        console.log('e', e);
      }

      setIsMarkPaidSaved(false);
      setLocalPaid('PAID');
    }
  }

  const handleCopyUrl = () => {
    navigator.clipboard.writeText(`${window.location.origin}${bountyPath}`);
    setIsCopied(true);

    setTimeout(() => {
      // UI Enhancement: Show "Copied" for 3 seconds, then reset
      setIsCopied(false);
    }, 3000);
  };

  async function sendBadge(body: any) {
    const { recipient, badge } = body;

    setSaving('badgeRecipient');
    try {
      if (badge?.amount < 1) {
        alert("You don't have any of the selected badge");
        throw new Error("You don't have any of the selected badge");
      }

      // first get the user's liquid address
      const recipientDetails = await main.getPersonByPubkey(recipient.owner_pubkey);

      const liquidAddress =
        recipientDetails?.extras?.liquid && recipientDetails?.extras?.liquid[0]?.value;

      if (!liquidAddress) {
        alert('This user has not provided an L-BTC address');
        throw new Error('This user has not provided an L-BTC address');
      }

      const pack = {
        asset: badge.id,
        to: liquidAddress,
        amount: 1,
        memo: props.ticket_url
      };

      const r = await main.sendBadgeOnLiquid(pack);

      if (r.ok) {
        await setExtrasPropertyAndSave('badgeRecipient', recipient.owner_pubkey);
        setShowBadgeAwardDialog(false);
      } else {
        alert(r.statusText || 'Operation failed! Contact support.');
        throw new Error(r.statusText);
      }
    } catch (e) {
      console.log(e);
    }

    setSaving('');
  }

  //  if my own, show this option to show/hide
  const markPaidButton = (
    <Button
      color={'primary'}
      iconSize={14}
      style={{ fontSize: 14, height: 48, width: '100%', marginBottom: 20 }}
      endingIcon={'paid'}
      text={paid ? unpaidString : paidString}
      loading={saving === 'paid'}
      onClick={(e: any) => {
        e.stopPropagation();
        setExtrasPropertyAndSave('paid', !paid);
      }}
    />
  );

  const awardBadgeButton = !badgeRecipient && (
    <Button
      color={'primary'}
      iconSize={14}
      endingIcon={'offline_bolt'}
      style={{ fontSize: 14, height: 48, width: '100%', marginBottom: 20 }}
      text={badgeRecipient ? 'Badge Awarded' : 'Award Badge'}
      disabled={badgeRecipient ? true : false}
      loading={saving === 'badgeRecipient'}
      onClick={(e: any) => {
        e.stopPropagation();
        if (!badgeRecipient) {
          setShowBadgeAwardDialog(true);
        }
      }}
    />
  );

  const actionButtons = isMine && (
    <ButtonRow>
      {showBadgeAwardDialog ? (
        <>
          <Form
            loading={saving === 'badgeRecipient'}
            smallForm
            buttonsOnBottom
            wrapStyle={{ padding: 0, margin: 0, maxWidth: '100%' }}
            close={() => setShowBadgeAwardDialog(false)}
            onSubmit={(e: any) => {
              sendBadge(e);
            }}
            submitText={'Send Badge'}
            schema={sendBadgeSchema}
          />
          <div style={{ height: 100 }} />
        </>
      ) : (
        <>
          {markPaidButton}
          {awardBadgeButton}
        </>
      )}
    </ButtonRow>
  );

  const nametag = (
    <NameTag
      iconSize={24}
      textSize={13}
      style={{ marginBottom: 10 }}
      {...person}
      created={created}
      widget={'wanted'}
    />
  );

  function renderCodingTask() {
    // const { status } = ticketUrl
    //   ? extractGithubIssueFromUrl(person, ticketUrl)
    //   : extractGithubIssue(person, repo, issue);

    let assigneeLabel: any = null;

    if (assigneeInfo) {
      if (!isMobile) {
        assigneeLabel = (
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              fontSize: 12,
              color: color.grayish.G100,
              marginTop: isMobile ? 20 : 0,
              marginLeft: '-16px'
            }}
          >
            <Img
              src={assigneeInfo.img || '/static/person_placeholder.png'}
              style={{ borderRadius: 30 }}
            />

            <Assignee
              color={color}
              onClick={() => {
                const profileUrl = `https://community.sphinx.chat/p/${assigneeInfo.owner_pubkey}`;
                sendToRedirect(profileUrl);
              }}
              style={{ marginLeft: 3, fontWeight: 500, cursor: 'pointer' }}
            >
              {assigneeInfo.owner_alias}
            </Assignee>
          </div>
        );
      } else {
        assigneeLabel = (
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              fontSize: 12,
              color: color.grayish.G100,
              marginLeft: '16px'
            }}
          >
            <Img
              src={assigneeInfo.img || '/static/person_placeholder.png'}
              style={{ borderRadius: 30 }}
            />

            <Assignee
              color={color}
              onClick={() => {
                const profileUrl = `https://community.sphinx.chat/p/${assigneeInfo.owner_pubkey}`;
                sendToRedirect(profileUrl);
              }}
              style={{ marginLeft: 3, fontWeight: 500, cursor: 'pointer' }}
            >
              {assigneeInfo.owner_alias}
            </Assignee>
          </div>
        );
      }
    }

    // desktop view
    if (fromBountyPage) {
      return (
        <CodingBounty
          {...props}
          labels={labels}
          person={person}
          awardDetails={awardDetails}
          setAwardDetails={setAwardDetails}
          isAssigned={isAssigned}
          dataValue={dataValue}
          assigneeValue={assigneeValue}
          assignedPerson={assignedPerson}
          changeAssignedPerson={changeAssignedPerson}
          sendToRedirect={sendToRedirect}
          handleCopyUrl={handleCopyUrl}
          isCopied={isCopied}
          setExtrasPropertyAndSave={setExtrasPropertyAndSave}
          setExtrasPropertyAndSaveMultiple={setExtrasPropertyAndSaveMultiple}
          setIsModalSideButton={setIsModalSideButton}
          replitLink={replitLink}
          assigneeHandlerOpen={assigneeHandlerOpen}
          setCreatorStep={setCreatorStep}
          setIsExtraStyle={setIsExtraStyle}
          awards={awards}
          handleAssigneeDetails={handleAssigneeDetails}
          peopleList={peopleList}
          setIsPaidStatusBadgeInfo={setIsPaidStatusBadgeInfo}
          bountyPrice={bountyPrice}
          selectedAward={selectedAward}
          handleAwards={handleAwards}
          repo={repo}
          issue={issue}
          isMarkPaidSaved={isMarkPaidSaved}
          setBountyPrice={setBountyPrice}
          owner_idURL={owner_idURL}
          createdURL={createdURL}
          creatorStep={creatorStep}
          isPaidStatusBadgeInfo={isPaidStatusBadgeInfo}
          isPaidStatusPopOver={isPaidStatusPopOver}
          titleString={titleString}
          nametag={nametag}
          org_uuid={org_uuid}
          id={id}
          localPaid={localPaid}
          setLocalPaid={setLocalPaid}
          isMobile={isMobile}
          actionButtons={actionButtons}
          assigneeLabel={assigneeLabel}
        />
      );
    }

    return (
      <div>
        <CodingDesktop
          {...props}
          labels={labels}
          nametag={nametag}
          actionButtons={actionButtons}
          assigneeLabel={assigneeLabel}
          assignee={assignee}
          loomEmbedUrl={loomEmbedUrl}
          titleString={titleString}
          isCopied={isCopied}
        />
      </div>
    );
  }

  if (type === 'coding_task' || type === 'wanted_coding_task' || type === 'freelance_job_request') {
    return renderCodingTask();
  }
  return <div />;
}
export default observer(WantedSummary);
