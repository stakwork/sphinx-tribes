import React, { useState, useEffect, useRef, useCallback } from 'react';
import styled from 'styled-components';
import PageLoadSpinner from 'people/utils/PageLoadSpinner';
import NoResults from 'people/utils/OrgNoResults';
import { useStores } from 'store';
import { Organization } from 'store/main';
import { Wrap } from 'components/form/style';
import { EuiGlobalToastList } from '@elastic/eui';
import { Button } from 'components/common';
import { useIsMobile } from 'hooks/uiHooks';
import { Formik } from 'formik';
import { FormField, validator } from 'components/form/utils';
import { DollarConverter, satToUsd } from 'helpers';
import { Modal } from '../../components/common';
import avatarIcon from '../../public/static/profile_avatar.svg';
import { colors } from '../../config/colors';
import { widgetConfigs } from '../utils/Constants';
import Input from '../../components/form/inputs';
import { Person } from '../../store/main';
import OrganizationDetails from './OrganizationDetails';
import ManageButton from './ManageOrgButton';

const color = colors['light'];

const Container = styled.div`
  display: flex;
  flex-flow: column wrap;
  min-width: 100%;
  min-height: 100%;
  flex: 1 1 100%;
`;

const OrganizationWrap = styled.div`
  display: flex;
  flex-direction: row;
  width: 100%;
  background: white;
  padding: 25px 30px;
  border-radius: 6px;
  cursor: pointer;
  @media only screen and (max-width: 800px) {
    padding: 15px 0px;
  }
  @media only screen and (max-width: 700px) {
    padding: 12px 0px;
    margin-bottom: 10px;
  }
  @media only screen and (max-width: 500px) {
    padding: 0px;
  }
`;

const OrganizationData = styled.div`
  display: flex;
  align-items: center;
  flex-direction: row;
  width: 100%;
  @media only screen and (max-width: 470px) {
    flex-direction: column;
    justify-content: center;
    border: 1px solid #ccc;
    border-radius: 10px;
    padding: 15px 0px;
  }
`;

const OrganizationImg = styled.img`
  width: 65px;
  height: 65px;
  @media only screen and (max-width: 700px) {
    width: 55px;
    height: 55px;
  }
  @media only screen and (max-width: 500px) {
    width: 48px;
    height: 48px;
  }
  @media only screen and (max-width: 470px) {
    width: 60px;
    height: 60px;
  }
`;

const OrganizationTextWrap = styled.div`
  margin-left: 20px;
  display: flex;
  flex-direction: column;
  @media only screen and (max-width: 470px) {
    margin-left: 0px;
    margin-top: 15px;
    justify-content: center;
  }
`;

const OrganizationText = styled.p`
  font-size: 1rem;
  font-weight: bold;
  @media only screen and (max-width: 700px) {
    font-size: 0.85rem;
  }
  @media only screen and (max-width: 500px) {
    font-size: 0.79rem;
  }
  @media only screen and (max-width: 470px) {
    font-size: 0.85rem;
    text-align: center;
  }
`;

const OrganizationBudgetText = styled.small`
  margin-top: auto;
  font-size: 0.9rem;
  @media only screen and (max-width: 700px) {
    font-size: 0.8rem;
  }
  @media only screen and (max-width: 500px) {
    font-size: 0.75rem;
  }
`;

const OrganizationContainer = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
  cursor: pointer;
  gap: 15px;
`;

const OrgHeadWrap = styled.div`
  display: flex;
  align-items: center;
  margin-top: 5px;
  margin-bottom: 20px;
`;

const OrgText = styled.div`
  font-size: 1.4rem;
  font-weight: bold;
  @media only screen and (max-width: 700px) {
    font-size: 1.1rem;
  }
  @media only screen and (max-width: 700px) {
    font-size: 0.95rem;
  }
`;
const OrganizationActionWrap = styled.div`
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 15px;
  @media only screen and (max-width: 470px) {
    margin-left: 0;
    margin-top: 20px;
  }
`;

const SatsGap = styled.span`
  margin: 0px 10px;
  @media only screen and (max-width: 700px) {
    margin: 0px 5px;
  }
`;

const AddOrgButton = styled(Button)`
  width: 100%;
  borderradius: 10;
  height: 45;
  margintop: 15;
`;

const Organizations = (props: { person: Person }) => {
  const [loading, setIsLoading] = useState<boolean>(false);
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const [detailsOpen, setDetailsOpen] = useState<boolean>(false);
  const [organization, setOrganization] = useState<Organization>();
  const [disableFormButtons, setDisableFormButtons] = useState(false);
  const [toasts, setToasts]: any = useState([]);
  const [user, setUser] = useState<Person>();
  const { main, ui } = useStores();
  const isMobile = useIsMobile();
  const config = widgetConfigs['organizations'];
  const formRef = useRef(null);
  const isMyProfile = ui?.meInfo?.pubkey === props?.person?.owner_pubkey;

  const schema = [...config.schema];

  const initValues = {
    name: '',
    img: '',
    show: false
  };

  function addToast(title: string) {
    setToasts([
      {
        id: '1',
        title,
        color: 'danger'
      }
    ]);
  }

  function removeToast() {
    setToasts([]);
  }

  const getUserOrganizations = useCallback(async () => {
    setIsLoading(true);
    if (ui.selectedPerson !== 0) {
      await main.getUserOrganizations(ui.selectedPerson);
      const user = await main.getPersonById(ui.selectedPerson);
      setUser(user);
    }
    setIsLoading(false);
  }, [main, ui.selectedPerson]);

  useEffect(() => {
    getUserOrganizations();
  }, [getUserOrganizations]);

  const closeHandler = () => {
    setIsOpen(false);
  };

  const closeDetails = () => {
    setDetailsOpen(false);
  };

  const onSubmit = async (body: any) => {
    setIsLoading(true);
    body.owner_pubkey = ui.meInfo?.owner_pubkey;
    const res = await main.addOrganization(body);
    if (res.status === 200) {
      await getUserOrganizations();
    } else {
      addToast('Error: could not create organization');
    }
    closeHandler();
    setIsLoading(false);
  };

  const orgUi = (org: any, key: number) => {
    const btnDisabled = (!org.bounty_count && org.bount_count !== 0) || !org.uuid;
    return (
      <OrganizationWrap key={key}>
        <OrganizationData>
          <OrganizationImg src={org.img || avatarIcon} />
          <OrganizationTextWrap>
            <OrganizationText>{org.name}</OrganizationText>
            <OrganizationBudgetText>
              {DollarConverter(org.budget ?? 0)}
              <SatsGap>/</SatsGap>
              {satToUsd(org.budget ?? 0)} USD
            </OrganizationBudgetText>
          </OrganizationTextWrap>
          <OrganizationActionWrap>
            <ManageButton
              org={org}
              user_pubkey={user?.owner_pubkey ?? ''}
              action={() => {
                setOrganization(org);
                setDetailsOpen(true);
              }}
            />
            <Button
              disabled={btnDisabled}
              color={!btnDisabled ? 'white' : 'grey'}
              text="View Bounties"
              endingIcon="open_in_new"
              onClick={() => window.open(`/org/bounties/${org.uuid}`, '_target')}
              style={{
                height: 40,
                color: '#000000',
                borderRadius: 10
              }}
            />
          </OrganizationActionWrap>
        </OrganizationData>
      </OrganizationWrap>
    );
  };

  const renderOrganizations = () => {
    if (main.organizations.length) {
      return (
        <>
          <OrgHeadWrap>
            <OrgText>Organizations</OrgText>
            {isMyProfile && (
              <Button
                leadingIcon={'add'}
                height={isMobile ? 40 : 45}
                text="Add Organization"
                onClick={() => setIsOpen(true)}
                style={{ marginLeft: 'auto', borderRadius: 10 }}
              />
            )}
          </OrgHeadWrap>
          <OrganizationContainer>
            {main.organizations.map((org: Organization, i: number) => orgUi(org, i))}
          </OrganizationContainer>
        </>
      );
    } else {
      return (
        <Container>
          <NoResults showAction={isMyProfile} action={() => setIsOpen(true)} />
        </Container>
      );
    }
  };

  return (
    <Container>
      <PageLoadSpinner show={loading} />
      {detailsOpen && <OrganizationDetails close={closeDetails} org={organization} />}
      {!detailsOpen && (
        <>
          {renderOrganizations()}
          {isOpen && (
            <Modal
              visible={isOpen}
              style={{
                height: '100%',
                flexDirection: 'column'
              }}
              envStyle={{
                marginTop: isMobile ? 64 : 0,
                background: color.pureWhite,
                zIndex: 20,
                ...(config?.modalStyle ?? {}),
                maxHeight: '100%',
                borderRadius: '10px'
              }}
              overlayClick={closeHandler}
              bigCloseImage={closeHandler}
              bigCloseImageStyle={{
                top: '-18px',
                right: '-18px',
                background: '#000',
                borderRadius: '50%'
              }}
            >
              <Formik
                initialValues={initValues || {}}
                onSubmit={onSubmit}
                innerRef={formRef}
                validationSchema={validator(schema)}
              >
                {({
                  setFieldTouched,
                  handleSubmit,
                  values,
                  setFieldValue,
                  errors,
                  initialValues
                }: any) => (
                  <Wrap newDesign={true}>
                    <h5>Add new organization</h5>
                    <div className="SchemaInnerContainer">
                      {schema.length &&
                        schema.map((item: FormField) => (
                          <Input
                            {...item}
                            key={item.name}
                            values={values}
                            errors={errors}
                            value={values[item.name]}
                            error={errors[item.name]}
                            initialValues={initialValues}
                            deleteErrors={() => {
                              if (errors[item.name]) delete errors[item.name];
                            }}
                            handleChange={(e: any) => {
                              setFieldValue(item.name, e);
                            }}
                            setFieldValue={(e: any, f: any) => {
                              setFieldValue(e, f);
                            }}
                            setFieldTouched={setFieldTouched}
                            handleBlur={() => setFieldTouched(item.name, false)}
                            handleFocus={() => setFieldTouched(item.name, true)}
                            setDisableFormButtons={setDisableFormButtons}
                            borderType={'bottom'}
                            imageIcon={true}
                            style={
                              item.name === 'github_description' && !values.ticket_url
                                ? {
                                    display: 'none'
                                  }
                                : undefined
                            }
                          />
                        ))}

                      <AddOrgButton
                        disabled={disableFormButtons || loading}
                        onClick={() => {
                          handleSubmit();
                        }}
                        loading={loading}
                        color={'primary'}
                        text={'Add Organization'}
                      />
                    </div>
                  </Wrap>
                )}
              </Formik>
            </Modal>
          )}
        </>
      )}
      <EuiGlobalToastList toasts={toasts} dismissToast={removeToast} toastLifeTimeMs={5000} />
    </Container>
  );
};

export default Organizations;
