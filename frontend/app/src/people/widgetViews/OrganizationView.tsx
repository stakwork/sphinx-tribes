import React, { useState, useEffect, useRef, useCallback } from 'react';
import styled from 'styled-components';
import PageLoadSpinner from 'people/utils/PageLoadSpinner';
import NoResults from 'people/utils/OrgNoResults';
import { useStores } from 'store';
import { Organization } from 'store/main';
import { Wrap } from 'components/form/style';
import { EuiGlobalToastList } from '@elastic/eui';
import { Link } from 'react-router-dom';
import { Button, IconButton } from 'components/common';
import { useIsMobile } from 'hooks/uiHooks';
import { Formik } from 'formik';
import { FormField, validator } from 'components/form/utils';
import { Modal } from '../../components/common';
import avatarIcon from '../../public/static/profile_avatar.svg';
import { colors } from '../../config/colors';
import { widgetConfigs } from '../utils/Constants';
import Input from '../../components/form/inputs';
import { Person } from '../../store/main';
import OrganizationDetails from './OrganizationDetails';

const color = colors['light'];

const Container = styled.div`
  display: flex;
  flex-flow: column wrap;
  gap: 1rem;
  min-width: 77vw;
  flex: 1 1 100%;
`;

const OrganizationText = styled.p`
  font-size: 1rem;
  text-transform: capitalize;
  font-weight: bold;
  margin-top: 15px;
`;

const OrganizationImg = styled.img`
  width: 60px;
  height: 60px;
`;

const OrganizationWrap = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  width: calc(19% - 40px);
  margin-left: 0.5%;
  margin-right: 0.5%;
  margin; 10px;
  background: white;
  padding: 20px;
  border-radius: 2px;
  cursor: pointer;
`;

const OrganizationData = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  width: 100%;
`;

const OrganizationContainer = styled.div`
  display: flex;
  flex-direction: row;
  width: 100%;
  cursor: pointer;
`;

const Organizations = (props: { person: Person }) => {
  const [loading, setIsLoading] = useState<boolean>(false);
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const [detailsOpen, setDetailsOpen] = useState<boolean>(false);
  const [organization, setOrganization] = useState<Organization>();
  const [disableFormButtons, setDisableFormButtons] = useState(false);
  const [toasts, setToasts]: any = useState([]);
  const { main, ui } = useStores();
  const isMobile = useIsMobile();
  const config = widgetConfigs['organizations'];
  const formRef = useRef(null);
  const isMyProfile = ui?.meInfo?.pubkey == props?.person?.owner_pubkey;

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
    await main.getUserOrganizations();
    setIsLoading(false);
  }, [main]);

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

  const renderOrganizations = () => {
    if (main.organizations.length) {
      return main.organizations.map((org: Organization, i: number) => (
        <OrganizationWrap key={i}>
          <OrganizationData
            onClick={() => {
              setOrganization(org);
              setDetailsOpen(true);
            }}
          >
            <OrganizationImg src={org.img || avatarIcon} />
            <OrganizationText>{org.name}</OrganizationText>
          </OrganizationData>
          <Link to={`/org/tickets/${org.uuid}`} target="_blank">
            Bounties
          </Link>
        </OrganizationWrap>
      ));
    } else {
      return <NoResults />;
    }
  };

  return (
    <Container>
      <PageLoadSpinner show={loading} />
      {detailsOpen && <OrganizationDetails close={closeDetails} org={organization} />}
      {!detailsOpen && (
        <>
          {isMyProfile && (
            <IconButton
              width={150}
              height={isMobile ? 36 : 48}
              text="Add Organization"
              onClick={() => setIsOpen(true)}
              style={{
                marginLeft: '10px'
              }}
            />
          )}
          <OrganizationContainer>{renderOrganizations()}</OrganizationContainer>
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

                      <Button
                        disabled={disableFormButtons || loading}
                        onClick={() => {
                          handleSubmit();
                        }}
                        loading={loading}
                        style={{ width: '100%' }}
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
