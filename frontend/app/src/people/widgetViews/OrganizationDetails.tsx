import React, { useCallback, useEffect, useRef, useState } from 'react';
import styled from 'styled-components';
import { useStores } from 'store';
import { OrgWrap, Wrap } from 'components/form/style';
import { EuiGlobalToastList } from '@elastic/eui';
import { InvoiceForm, InvoiceInput, InvoiceLabel } from 'people/utils/style';
import moment from 'moment';
import { SOCKET_MSG, createSocketInstance } from 'config/socket';
import { Button } from 'components/common';
import { useIsMobile } from 'hooks/uiHooks';
import { Formik } from 'formik';
import { FormField, validator } from 'components/form/utils';
import { BountyRoles, BudgetHistory, Organization, PaymentHistory, Person } from 'store/main';
import MaterialIcon from '@material/react-material-icon';
import { Route, Router, Switch, useRouteMatch } from 'react-router-dom';
import { satToUsd, userHasRole } from 'helpers';
import { BountyModal } from 'people/main/bountyModal';
import history from '../../config/history';
import { Modal } from '../../components/common';
import { colors } from '../../config/colors';
import { nonWidgetConfigs } from '../utils/Constants';
import Invoice from '../widgetViews/summaries/wantedSummaries/Invoice';
import Input from '../../components/form/inputs';
import avatarIcon from '../../public/static/profile_avatar.svg';
import DeleteTicketModal from './DeleteModal';

const color = colors['light'];

const Container = styled.div`
  width: 100%;
  min-height: 100%;
  background: white;
  padding: 20px 0px;
  padding-top: 0px;
  z-index: 100;
`;

const HeadWrap = styled.div`
  display: flex;
  align-items: center;
  padding: 25px 40px;
  border-bottom: 1px solid #ebedef;
  @media only screen and (max-width: 800px) {
    padding: 15px 0px;
  }
  @media only screen and (max-width: 700px) {
    padding: 12px 0px;
  }
  @media only screen and (max-width: 500px) {
    padding: 0px;
    padding-bottom: 15px;
    flex-direction: column;
    align-items: start;
  }
`;

const HeadNameWrap = styled.div`
  display: flex;
  align-items: center;
  @media only screen and (max-width: 500px) {
    margin-bottom: 20px;
  }
`;

const OrgImg = styled.img`
  width: 48px;
  height: 48px;
  border-radius: 50%;
  margin-left: 20px;
  @media only screen and (max-width: 700px) {
    width: 42px;
    height: 42px;
  }
  @media only screen and (max-width: 500px) {
    width: 38px;
    height: 38px;
  }
  @media only screen and (max-width: 470px) {
    width: 35px;
    height: 35px;
  }
`;

const OrgName = styled.h3`
  padding: 0px;
  margin: 0px;
  font-size: 1.3rem;
  margin-left: 25px;
  font-weight: 700;
  margin-left: 20px;
  @media only screen and (max-width: 800px) {
    font-size: 1.05rem;
  }
  @media only screen and (max-width: 700px) {
    font-size: 1rem;
  }
  @media only screen and (max-width: 600px) {
    font-size: 0.9rem;
  }
  @media only screen and (max-width: 470px) {
    font-size: 0.8rem;
  }
`;

const HeadButtonWrap = styled.div<{ forSmallScreen: boolean }>`
  margin-left: auto;
  display: flex;
  flex-direction: row;
  gap: 15px;
  @media only screen and (max-width: 700px) {
    gap: 10px;
    margin-left: auto;
  }
  @media only screen and (max-width: 500px) {
    gap: 8px;
    margin-left: 0px;
    width: 100vw;
    margin-left: ${(p: any) => (p.forSmallScreen ? '50px' : '0px')};
  }
  @media only screen and (max-width: 470px) {
    gap: 6px;
  }
`;

const DetailsWrap = styled.div`
  width: 100%;
  min-height: 100%;
  margin-top: 17px;
  padding: 0px 20px;
`;

const ActionWrap = styled.div`
  display: flex;
  align-items: center;
  padding: 25px 40px;
  border-bottom: 1px solid #ebedef;
  @media only screen and (max-width: 700px) {
    padding: 25px 0px;
  }
  @media only screen and (max-width: 500px) {
    flex-direction: column;
    width: 100%;
    padding: 25px 0px;
  }
`;

const BudgetWrap = styled.div`
  padding: 25px 40px;
  width: 55%;
  display: flex;
  flex-direction: column;
  @media only screen and (max-width: 700px) {
    width: 100%;
    padding: 22px 0px;
  }
  @media only screen and (max-width: 500px) {
    width: 100%;
    padding: 20px 0px;
  }
`;

const NoBudgetWrap = styled.div`
  display: flex;
  flex-direction: row;
  align-items: center;
  width: 100%;
  border: 1px solid #ebedef;
`;

const ViewBudgetWrap = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
`;

const BudgetSmall = styled.h6`
  padding: 0px;
  font-size: 0.8rem;
  color: #8e969c;
  @media only screen and (max-width: 500px) {
    font-size: 0.75rem;
  }
`;

const BudgetSmallHead = styled.h6`
  padding: 0px;
  font-size: 0.7rem;
  color: #8e969c;
`;

const Budget = styled.h4`
  color: #3c3f41;
  font-size: 1.15rem;
  @media only screen and (max-width: 500px) {
    font-size: 1rem;
  }
`;

const Grey = styled.span`
  color: #8e969c;
`;

const NoBudgetText = styled.p`
  font-size: 0.85rem;
  padding: 0px;
  margin: 0px;
  color: #8e969c;
  width: 90%;
  margin-left: auto;
`;

const UserWrap = styled.div`
  display: flex;
  flex-direction: column;
  padding: 25px 40px;
  @media only screen and (max-width: 700px) {
    width: 100%;
    padding: 20px 0px;
  }
  @media only screen and (max-width: 500px) {
    padding: 20px 0px;
  }
`;

const UsersHeadWrap = styled.div`
  display: flex;
  align-items: center;
  width: 100%;
  border-bottom: 1px solid #ebedef;
  padding-top: 5px;
  padding-bottom: 20px;
  @media only screen and (max-width: 500px) {
    width: 100%;
  }
`;

const UsersHeader = styled.h4`
  font-size: 0.9rem;
  font-weight: 600;
  padding: 0;
  margin: 0;
  @media only screen and (max-width: 500px) {
    font-size: 0.8rem;
    margin-right: 55%;
  }
`;

const UsersList = styled.div`
  @media only screen and (max-width: 500px) {
    width: 100%;
  }
`;

const UserImage = styled.img`
  width: 60px;
  height: 60px;
  border-radius: 50%;
  align-self: center;
`;

const User = styled.div`
  padding: 15px 0px;
  border-bottom: 1px solid #ebedef;
  display: flex;
  align-items: center;
  @media only screen and (max-width: 500px) {
    padding: 10px 0px;
    width: 100%;
  }
`;

const UserDetails = styled.div`
  display: flex;
  flex-gap: 12px;
  flex-direction: column;
  margin-left: 2%;
  width: 30%;
  @media only screen and (max-width: 500px) {
    width: 60%;
    margin-left: 5%;
  }
`;

const UserName = styled.p`
  padding: 0px;
  margin: 0px;
  font-size: 0.9rem;
  text-transform: capitalize;
  font-weight: bold;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  width: 100%;
`;

const UserPubkey = styled.p`
  padding: 0px;
  margin: 0px;
  font-size: 0.75rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  width: 100%;
  color: #5f6368;
`;

const UserAction = styled.div`
  display: flex;
  align-items: center;
  margin-left: auto;
`;

const IconWrap = styled.div`
  :first-child {
    margin-right: 40px;
    @media only screen and (max-width: 700px) {
      margin-right: 20px;
    }
    @media only screen and (max-width: 500px) {
      margin-right: 8px;
    }
  }
`;

const ModalTitle = styled.h3`
  font-weight: 800;
  font-size: 26
  color: ##3C3F41;
`;

const CheckUl = styled.ul`
  list-style: none;
  padding: 0;
  margin-top: 20px;
  justifyContent: 'left',
  width: 100%
`;

const CheckLi = styled.li`
  display: flex;
  flex-direction: row;
  align-items: center;
  padding: 0px;
  margin-bottom: 10px;
`;

const Check = styled.input`
  width: 20px;
  height: 20px;
  border-radius: 5px;
  padding: 0px;
  margin-right: 10px;
`;

const CheckLabel = styled.label`
  padding: 0px;
  margin: 0px;
`;

const ViewBounty = styled.p`
  padding: 0px;
  margin: 0px;
  cursor: pointer;
  font-size: 0.9rem;
  color: green;
  font-size: bold;
`;

const HeadButton = styled(Button)`
  border-radius: 5px;
`;

const UserRolesWrap = styled(Wrap)`
  width: 100%;
`;

const UserRolesHeader = styled.div`
  display: flex;
  flex-direction: row;
`;

const UserRolesName = styled.p`
  color: #8E969C;
  margin: 5px 0px;
`;

// the view for all details about an organization
const OrganizationDetails = (props: { close: () => void; org: Organization | undefined }) => {
  const [loading, setIsLoading] = useState<boolean>(false);
  const isMobile = useIsMobile();
  const { main, ui } = useStores();
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const [isOpenRoles, setIsOpenRoles] = useState<boolean>(false);
  const [isOpenBudget, setIsOpenBudget] = useState<boolean>(false);
  const [isOpenHistory, setIsOpenHistory] = useState<boolean>(false);
  const [isOpenBudgetHistory, setIsOpenBudgetHistory] = useState<boolean>(false);
  const [showDeleteModal, setShowDeleteModal] = useState<boolean>(false);
  const [orgBudget, setOrgBudget] = useState<number>(0);
  const [paymentsHistory, setPaymentsHistory] = useState<PaymentHistory[]>([]);
  const [budgetsHistory, setBudgetsHistory] = useState<BudgetHistory[]>([]);
  const [disableFormButtons, setDisableFormButtons] = useState(false);
  const [users, setUsers] = useState<Person[]>([]);
  const [user, setUser] = useState<Person>();
  const [userRoles, setUserRoles] = useState<any[]>([]);
  const [bountyRolesData, setBountyRolesData] = useState<BountyRoles[]>([]);
  const [toasts, setToasts]: any = useState([]);
  const [lnInvoice, setLnInvoice] = useState('');
  const [invoiceStatus, setInvoiceStatus] = useState(false);
  const [amount, setAmount] = useState(1);
  const { path, url } = useRouteMatch();

  const pollMinutes = 2;

  const config = nonWidgetConfigs['organizationusers'];

  const formRef = useRef(null);
  const isOrganizationAdmin = props.org?.owner_pubkey === ui.meInfo?.owner_pubkey;
  const schema = [...config.schema];

  const addUserDisabled =
    !isOrganizationAdmin && !userHasRole(main.bountyRoles, userRoles, 'ADD USER');
  const viewReportDisabled =
    !isOrganizationAdmin && !userHasRole(main.bountyRoles, userRoles, 'VIEW REPORT');
  const addBudgetDisabled =
    !isOrganizationAdmin && !userHasRole(main.bountyRoles, userRoles, 'ADD BUDGET');
  const deleteUserDisabled =
    !isOrganizationAdmin && !userHasRole(main.bountyRoles, userRoles, 'DELETE USER');
  const addRolesDisabled =
    !isOrganizationAdmin && !userHasRole(main.bountyRoles, userRoles, 'ADD ROLES');

  const initValues = {
    owner_pubkey: ''
  };

  const { org } = props;
  const uuid = org?.uuid || '';

  function addToast(title: string, color: 'danger' | 'success') {
    setToasts([
      {
        id: '1',
        title,
        color
      }
    ]);
  }

  function removeToast() {
    setToasts([]);
  }

  const getOrganizationUsers = useCallback(async () => {
    if (uuid) {
      const users = await main.getOrganizationUsers(uuid);
      setUsers(users);
    }
  }, [main, uuid]);

  const deleteOrganizationUser = async (user: any) => {
    if (uuid) {
      const res = await main.deleteOrganizationUser(user, uuid);

      if (res.status === 200) {
        await getOrganizationUsers();
      } else {
        addToast('Error: could not delete user', 'danger');
      }
    }
  };

  const closeDeleteModal = () => setShowDeleteModal(false);

  const confirmDelete = async () => {
    try {
      if (user) {
        await deleteOrganizationUser(user);
      }
    } catch (error) {
      console.log(error);
    }
    closeDeleteModal();
  };

  const getBountyRoles = useCallback(async () => {
    const bountyRolesData = main.bountyRoles.map((role: any) => ({
      name: role.name,
      status: false
    }));
    setBountyRolesData(bountyRolesData);
  }, [main.bountyRoles]);

  const getUserRoles = async (user: any) => {
    if (uuid && user.owner_pubkey) {
      const userRoles = await main.getUserRoles(uuid, user.owner_pubkey);
      setUserRoles(userRoles);

      // set all values to false, so every user data will be fresh
      const rolesData = bountyRolesData.map((data: any) => ({ name: data.name, status: false }));

      userRoles.forEach((userRole: any) => {
        const index = rolesData.findIndex((role: any) => role.name === userRole.role);
        rolesData[index]['status'] = true;
      });

      setBountyRolesData(rolesData);
    }
  };

  const getOrganizationBudget = useCallback(async () => {
    const organizationBudget = await main.getOrganizationBudget(uuid);
    setOrgBudget(organizationBudget.total_budget);
  }, [main, uuid]);

  const getPaymentsHistory = useCallback(async () => {
    const paymentHistories = await main.getPaymentHistories(uuid);
    setPaymentsHistory(paymentHistories);
  }, [main, uuid]);

  const getBudgetHistory = useCallback(async () => {
    const budgetHistories = await main.getBudgettHistories(uuid);
    setBudgetsHistory(budgetHistories);
  }, [main, uuid]);

  const generateInvoice = async () => {
    const token = ui.meInfo?.websocketToken;
    if (token) {
      const data = await main.getBudgetInvoice({
        amount: amount,
        sender_pubkey: ui.meInfo?.owner_pubkey ?? '',
        org_uuid: uuid,
        websocket_token: token,
        payment_type: 'deposit'
      });

      setLnInvoice(data.response.invoice);
    }
  };

  const handleSettingsClick = async (user: any) => {
    setUser(user);
    setIsOpenRoles(true);
    getUserRoles(user);
  };

  const handleDeleteClick = async (user: any) => {
    setUser(user);
    setShowDeleteModal(true);
  };

  const closeHandler = () => {
    setIsOpen(false);
  };

  const closeRolesHandler = () => {
    setIsOpenRoles(false);
  };

  const closeBudgetHandler = () => {
    setIsOpenBudget(false);
  };

  const closeHistoryHandler = () => {
    setIsOpenHistory(false);
  };

  const closeBudgetHistoryHandler = () => {
    setIsOpenBudgetHistory(false);
  };

  const onSubmit = async (body: any) => {
    setIsLoading(true);

    body.org_uuid = uuid;

    const res = await main.addOrganizationUser(body);
    if (res.status === 200) {
      await getOrganizationUsers();
    } else {
      addToast('Error: could not add user', 'danger');
    }
    closeHandler();
    setIsLoading(false);
  };

  const roleChange = (e: any) => {
    const rolesData = bountyRolesData.map((role: any) => {
      if (role.name === e.target.value) {
        role.status = !role.status;
      }
      return role;
    });

    setBountyRolesData(rolesData);
  };

  const submitRoles = async () => {
    const roleData = bountyRolesData
      .filter((r: any) => r.status)
      .map((role: any) => ({
        owner_pubkey: user?.owner_pubkey,
        org_uuid: uuid,
        role: role.name
      }));

    if (uuid && user?.owner_pubkey) {
      const res = await main.addUserRoles(roleData, uuid, user.owner_pubkey);
      if (res.status === 200) {
        await main.getUserRoles(uuid, user.owner_pubkey);
      } else {
        addToast('Error: could not add user roles', 'danger');
      }
      setIsOpenRoles(false);
    }
  };

  const onHandle = (event: any) => {
    const res = JSON.parse(event.data);
    if (res.msg === SOCKET_MSG.user_connect) {
      const user = ui.meInfo;
      if (user) {
        user.websocketToken = res.body;
        ui.setMeInfo(user);
      }
    } else if (res.msg === SOCKET_MSG.budget_success && res.invoice === main.lnInvoice) {
      addToast('Budget was added successfully', 'success');
      setLnInvoice('');
      setInvoiceStatus(true);
      main.setLnInvoice('');

      // get new organization budget
      getOrganizationBudget();
      getBudgetHistory();
      main.getUserOrganizations(ui.selectedPerson);
      closeBudgetHandler();
    }
  };

  const viewBounty = async (bountyId: number) => {
    ui.setBountyPerson(ui.meInfo?.id);

    history.push({
      pathname: `${url}/${bountyId}/${0}`
    });
  };

  useEffect(() => {
    getOrganizationUsers();
    getBountyRoles();
    getOrganizationBudget();
    getPaymentsHistory();
    getBudgetHistory();
  }, [
    getOrganizationUsers,
    getBountyRoles,
    getOrganizationBudget,
    getPaymentsHistory,
    getBudgetHistory
  ]);

  useEffect(() => {
    const socket: WebSocket = createSocketInstance();
    socket.onopen = () => {
      console.log('Socket connected');
    };

    socket.onmessage = (event: MessageEvent) => {
      onHandle(event);
    };

    socket.onclose = () => {
      console.log('Socket disconnected');
    };
  }, []);

  return (
    <Container>
      <HeadWrap>
        <HeadNameWrap>
          <MaterialIcon
            onClick={() => props.close()}
            icon={'arrow_back'}
            style={{
              fontSize: 25,
              cursor: 'pointer'
            }}
          />
          <OrgImg src={org?.img || avatarIcon} />
          <OrgName>{org?.name}</OrgName>
        </HeadNameWrap>
        <HeadButtonWrap forSmallScreen={false}>
          <HeadButton text="Edit" disabled={true} color="white" style={{ borderRadius: '5px' }} />
          <Button
            disabled={!org?.bounty_count}
            text="View Bounties"
            color="white"
            style={{ borderRadius: '5px' }}
            endingIcon="open_in_new"
            onClick={() => window.open(`/org/bounties/${uuid}`, '_target')}
          />
        </HeadButtonWrap>
      </HeadWrap>
      <ActionWrap>
        <BudgetWrap>
          {viewReportDisabled ? (
            <NoBudgetWrap>
              <MaterialIcon
                icon={'lock'}
                style={{
                  fontSize: 30,
                  cursor: 'pointer',
                  color: '#ccc'
                }}
              />
              <NoBudgetText>
                You have restricted permissions and are unable to view the budget. Reach out to the
                organization admin to get them updated.
              </NoBudgetText>
            </NoBudgetWrap>
          ) : (
            <ViewBudgetWrap>
              <BudgetSmallHead>YOUR BALANCE</BudgetSmallHead>
              <Budget>
                {orgBudget.toLocaleString()} <Grey>SATS</Grey>
              </Budget>
              <BudgetSmall>{satToUsd(orgBudget)} USD</BudgetSmall>
            </ViewBudgetWrap>
          )}
        </BudgetWrap>
        <HeadButtonWrap forSmallScreen={true}>
          <Button
            disabled={viewReportDisabled}
            text="History"
            color="white"
            style={{ borderRadius: '5px' }}
            onClick={() => setIsOpenHistory(true)}
          />
          <Button disabled={true} text="Withdraw" color="white" style={{ borderRadius: '5px' }} />
          <Button
            disabled={addBudgetDisabled}
            text="Deposit"
            color="white"
            style={{ borderRadius: '5px' }}
            onClick={() => setIsOpenBudget(true)}
          />
        </HeadButtonWrap>
      </ActionWrap>
      <UserWrap>
        <UsersHeadWrap>
          <UsersHeader>Users</UsersHeader>
          <HeadButtonWrap forSmallScreen={false}>
            <Button
              disabled={addUserDisabled}
              text="Add User"
              color="white"
              style={{
                borderRadius: '5px'
              }}
              onClick={() => setIsOpen(true)}
            />
          </HeadButtonWrap>
        </UsersHeadWrap>
        <UsersList>
          {users.map((user: Person, i: number) => (
            <User key={i}>
              <UserImage src={user.img || avatarIcon} />
              <UserDetails>
                <UserName>{user.unique_name}</UserName>
                <UserPubkey>{user.owner_pubkey}</UserPubkey>
              </UserDetails>
              <UserAction>
                <IconWrap>
                  <MaterialIcon
                    disabled={addRolesDisabled}
                    icon={'settings'}
                    style={{
                      fontSize: 24,
                      cursor: 'pointer',
                      color: '#ccc'
                    }}
                    onClick={() => handleSettingsClick(user)}
                  />
                </IconWrap>
                <IconWrap>
                  <MaterialIcon
                    icon={'delete'}
                    disabled={deleteUserDisabled}
                    style={{
                      fontSize: 24,
                      cursor: 'pointer',
                      color: '#ccc'
                    }}
                    onClick={() => {
                      setUser(user);
                      handleDeleteClick(user);
                    }}
                  />
                </IconWrap>
              </UserAction>
            </User>
          ))}
        </UsersList>
      </UserWrap>
      <DetailsWrap>
        {showDeleteModal && (
          <DeleteTicketModal
            closeModal={closeDeleteModal}
            confirmDelete={confirmDelete}
            text={'User'}
            imgUrl={user?.img}
            userDelete={true}
          />
        )}
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
                  <ModalTitle>Add new user</ModalTitle>
                  <div className="SchemaInnerContainer">
                    {schema.map((item: FormField) => (
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
                      text={'Add user'}
                    />
                  </div>
                </Wrap>
              )}
            </Formik>
          </Modal>
        )}
        {isOpenRoles && (
          <Modal
            visible={isOpenRoles}
            style={{
              height: '100%',
              flexDirection: 'column',
            }}
            envStyle={{
              marginTop: isMobile ? 64 : 0,
              background: color.pureWhite,
              zIndex: 20,
              ...(config?.modalStyle ?? {}),
              maxHeight: '100%',
              borderRadius: '10px',
              padding: '20px 60px 10px 60px'
            }}
            overlayClick={closeRolesHandler}
            bigCloseImage={closeRolesHandler}
            bigCloseImageStyle={{
              top: '-18px',
              right: '-18px',
              background: '#000',
              borderRadius: '50%'
            }}
          >
            <UserRolesWrap newDesign={true}>
              <UserRolesHeader>
                <div>
                  <UserRolesName>{user?.unique_name}</UserRolesName>
                  <ModalTitle>User Roles</ModalTitle>
                </div>
                <UserImage style={{height: '80px', width:'auto', marginLeft: 'auto'}} src={user?.img || avatarIcon} />
              </UserRolesHeader>
              <CheckUl>
                {bountyRolesData.map((role: any, i: number) => {
                  const capitalizeWords:string = role.name.charAt(0).toUpperCase() + role.name.slice(1).toLocaleLowerCase();
                  
                  return (<CheckLi key={i}>
                    <Check
                      checked={role.status}
                      onChange={roleChange}
                      type="checkbox"
                      name={role.name}
                      value={role.name}
                    />
                    <CheckLabel>{capitalizeWords}</CheckLabel>
                  </CheckLi>
                )})}
              </CheckUl>
              <Button
                onClick={() => submitRoles()}
                style={{ width: '150px', height: '50px', borderRadius: '5px', alignSelf: 'center' }}
                color={'primary'}
                text={'Update roles'}
              />
            </UserRolesWrap>
          </Modal>
        )}
        {isOpenBudget && (
          <Modal
            visible={isOpenBudget}
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
            overlayClick={closeBudgetHandler}
            bigCloseImage={closeBudgetHandler}
            bigCloseImageStyle={{
              top: '-18px',
              right: '-18px',
              background: '#000',
              borderRadius: '50%'
            }}
          >
            <Wrap newDesign={true}>
              <ModalTitle>Add budget</ModalTitle>
              {lnInvoice && ui.meInfo?.owner_pubkey && (
                <>
                  <Invoice
                    startDate={new Date(moment().add(pollMinutes, 'minutes').format().toString())}
                    invoiceStatus={invoiceStatus}
                    lnInvoice={lnInvoice}
                    invoiceTime={pollMinutes}
                  />
                </>
              )}
              {!lnInvoice && ui.meInfo?.owner_pubkey && (
                <>
                  <InvoiceForm>
                    <InvoiceLabel
                      style={{
                        display: 'block'
                      }}
                    >
                      Amount (in sats)
                    </InvoiceLabel>
                    <InvoiceInput
                      type="number"
                      style={{
                        width: '100%'
                      }}
                      value={amount}
                      onChange={(e: any) => setAmount(Number(e.target.value))}
                    />
                  </InvoiceForm>
                  <Button
                    text={'Generate Invoice'}
                    color={'primary'}
                    style={{ paddingLeft: 25, margin: '12px 0 10px' }}
                    img={'sphinx_white.png'}
                    imgSize={27}
                    height={48}
                    width={'100%'}
                    onClick={generateInvoice}
                  />
                </>
              )}
            </Wrap>
          </Modal>
        )}
        {isOpenHistory && (
          <Modal
            visible={isOpenHistory}
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
            overlayClick={closeHistoryHandler}
            bigCloseImage={closeHistoryHandler}
            bigCloseImageStyle={{
              top: '-18px',
              right: '-18px',
              background: '#000',
              borderRadius: '50%'
            }}
          >
            <OrgWrap style={{ width: '300px' }}>
              <ModalTitle>Payment history</ModalTitle>
              <table>
                <thead>
                  <tr>
                    <th>Sender</th>
                    <th>Recipient</th>
                    <th>Amount</th>
                    <th>Date</th>
                    <th />
                  </tr>
                </thead>
                <tbody>
                  {paymentsHistory.map((pay: PaymentHistory, i: number) => (
                    <tr key={i}>
                      <td className="ellipsis">{pay.sender_name}</td>
                      <td className="ellipsis">{pay.receiver_name}</td>
                      <td>{pay.amount} sats</td>
                      <td>{moment(pay.created).format('DD/MM/YY')}</td>
                      <td>
                        <ViewBounty onClick={() => viewBounty(pay.bounty_id)}>
                          View bounty
                        </ViewBounty>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </OrgWrap>
          </Modal>
        )}
        {isOpenBudgetHistory && (
          <Modal
            visible={isOpenBudgetHistory}
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
            overlayClick={closeBudgetHistoryHandler}
            bigCloseImage={closeBudgetHistoryHandler}
            bigCloseImageStyle={{
              top: '-18px',
              right: '-18px',
              background: '#000',
              borderRadius: '50%'
            }}
          >
            <OrgWrap>
              <ModalTitle>Budget history</ModalTitle>
              <table>
                <thead>
                  <tr>
                    <th>Sender</th>
                    <th>Amount</th>
                    <th>Type</th>
                    <th>Status</th>
                    <th>Date</th>
                  </tr>
                </thead>
                <tbody>
                  {budgetsHistory.map((b: BudgetHistory, i: number) => (
                    <tr key={i}>
                      <td className="ellipsis">{b.sender_name}</td>
                      <td>{b.amount} sats</td>
                      <td>{b.payment_type}</td>
                      <td>{b.status ? 'settled' : 'peending'}</td>
                      <td>{moment(b.created).format('DD/MM/YY')}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </OrgWrap>
          </Modal>
        )}
      </DetailsWrap>
      <EuiGlobalToastList toasts={toasts} dismissToast={removeToast} toastLifeTimeMs={5000} />
      <Router history={history}>
        <Switch>
          <Route path={`${path}/:wantedId/:wantedIndex`}>
            <BountyModal basePath={url} />
          </Route>
        </Switch>
      </Router>
    </Container>
  );
};

export default OrganizationDetails;
