import React, { useCallback, useEffect, useRef, useState } from 'react';
import styled from 'styled-components';
import { useStores } from 'store';
import { OrgWrap, Wrap } from 'components/form/style';
import { EuiGlobalToastList } from '@elastic/eui';
import { InvoiceForm, InvoiceInput, InvoiceLabel } from 'people/utils/style';
import moment from 'moment';
import { SOCKET_MSG, createSocketInstance } from 'config/socket';
import { Button, IconButton } from 'components/common';
import { useIsMobile } from 'hooks/uiHooks';
import { Formik } from 'formik';
import { FormField, validator } from 'components/form/utils';
import { BountyRoles, BudgetHistory, Organization, PaymentHistory, Person } from 'store/main';
import MaterialIcon from '@material/react-material-icon';
import { Route, Router, Switch, useRouteMatch } from 'react-router-dom';
import { userHasRole } from 'helpers';
import { BountyModal } from 'people/main/bountyModal';
import history from '../../config/history';
import { Modal } from '../../components/common';
import { colors } from '../../config/colors';
import { nonWidgetConfigs } from '../utils/Constants';
import Invoice from '../widgetViews/summaries/wantedSummaries/Invoice';
import Input from '../../components/form/inputs';

const color = colors['light'];

const Container = styled.div`
  width: 100%;
  min-height: 100%;
  background: white;
  padding: 20px 0px;
  z-index: 100;
`;

const DetailsWrap = styled.div`
  width: 100%;
  min-height: 100%;
  margin-top: 17px;
  padding: 0px 20px;
`;

const OrgInfoWrap = styled.div`
  display: flex;
  align-items: center;
`;

const DataCount = styled.div`
  margin-bottom: 15px;
  display: flex;
  align-items: center;
  margin-right: 20px;
`;

const DataText = styled.h3`
  font-size: 1.3rem;
  padding: 0px;
  margin: 0px;
  margin-right: 10px;
`;

const ViewHistoryText = styled.p`
  padding: 0px;
  margin: 0px;
  margin-left: 10px;
  font-size: 0.9rem;
  cursor: pointer;
`;

const UsersTable = styled.div`
  display: flex;
  flex-direction: column;
  margin-top: 25px;
`;

const TableRow = styled.div`
  display: flex;
  flex-direction: row;
  padding: 10px;
`;

const TableHead = styled.div`
  display: flex;
  flex-direction: row;
  padding: 10px;
  background: #d3d3d3;
`;

const ModalTitle = styled.h3`
  font-size: 1.2rem;
`;

const Th = styled.div`
  font-size: 1.1rem;
  font-weight: bold;
  min-width: 25%;
`;

const ThKey = styled.div`
  font-size: 1.1rem;
  font-weight: bold;
  min-width: 50%;
`;

const Td = styled.div`
  font-size: 0.95rem;
  min-width: 25%;
  text-transform: capitalize;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
`;

const TdKey = styled.div`
  font-size: 0.95rem;
  min-width: 50%;
  text-transform: capitalize;
`;

const Actions = styled.div`
  font-size: 0.95rem;
  min-width: 25%;
`;

const CheckUl = styled.ul`
  list-style: none;
  padding: 0;
  margin-top: 20px;
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

const OrganizationDetails = (props: { close: () => void; org: Organization | undefined }) => {
  const [loading, setIsLoading] = useState<boolean>(false);
  const isMobile = useIsMobile();
  const { main, ui } = useStores();
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const [isOpenRoles, setIsOpenRoles] = useState<boolean>(false);
  const [isOpenBudget, setIsOpenBudget] = useState<boolean>(false);
  const [isOpenHistory, setIsOpenHistory] = useState<boolean>(false);
  const [isOpenBudgetHistory, setIsOpenBudgetHistory] = useState<boolean>(false);
  const [usersCount, setUsersCount] = useState<number>(0);
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

  const initValues = {
    owner_pubkey: ''
  };

  const uuid = props.org?.uuid || '';

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

  const getOrganizationUsersCount = useCallback(async () => {
    if (uuid) {
      const count = await main.getOrganizationUsersCount(uuid);
      setUsersCount(count);
    }
  }, [main, uuid]);

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
        await getOrganizationUsersCount();
      } else {
        addToast('Error: could not delete user', 'danger');
      }
    }
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
      await getOrganizationUsersCount();
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
    getOrganizationUsersCount();
    getBountyRoles();
    getOrganizationBudget();
    getPaymentsHistory();
    getBudgetHistory();
  }, [
    getOrganizationUsers,
    getOrganizationUsersCount,
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
  }, [onHandle]);

  return (
    <Container>
      <MaterialIcon
        onClick={() => props.close()}
        icon={'arrow_back'}
        style={{
          fontSize: 30,
          marginLeft: 15,
          cursor: 'pointer'
        }}
      />
      <DetailsWrap>
        <OrgInfoWrap>
          <DataCount>
            <DataText>
              User{usersCount > 1 && 's'} {usersCount}
            </DataText>
            {(isOrganizationAdmin || userHasRole(main.bountyRoles, userRoles, 'ADD USER')) && (
              <IconButton
                width={80}
                height={isMobile ? 36 : 40}
                text="Add"
                onClick={() => setIsOpen(true)}
              />
            )}
          </DataCount>
          <DataCount>
            <DataText>Budget {orgBudget} sats</DataText>
            {(isOrganizationAdmin || userHasRole(main.bountyRoles, userRoles, 'ADD BUDGET')) && (
              <IconButton
                width={80}
                height={isMobile ? 36 : 40}
                text="Add"
                onClick={() => setIsOpenBudget(true)}
              />
            )}
            {(isOrganizationAdmin || userHasRole(main.bountyRoles, userRoles, 'VIEW REPORT')) && (
              <>
                <ViewHistoryText onClick={() => setIsOpenBudgetHistory(true)}>
                  Budget history
                </ViewHistoryText>
                <ViewHistoryText onClick={() => setIsOpenHistory(true)}>
                  Payment history
                </ViewHistoryText>
              </>
            )}
          </DataCount>
        </OrgInfoWrap>

        <UsersTable>
          <TableHead>
            <Th>Unique name</Th>
            <ThKey>Public key</ThKey>
            <Th>User actions</Th>
          </TableHead>
          {users.map((user: Person, i: number) => (
            <TableRow key={i}>
              <Td>{user.unique_name}</Td>
              <TdKey>{user.owner_pubkey}</TdKey>
              <Td>
                <Actions>
                  {(isOrganizationAdmin ||
                    userHasRole(main.bountyRoles, userRoles, 'ADD ROLES')) && (
                      <MaterialIcon
                        onClick={() => handleSettingsClick(user)}
                        icon={'settings'}
                        style={{
                          fontSize: 20,
                          marginLeft: 10,
                          cursor: 'pointer',
                          color: 'green'
                        }}
                      />
                    )}
                  {(isOrganizationAdmin ||
                    userHasRole(main.bountyRoles, userRoles, 'DELETE USER')) && (
                      <MaterialIcon
                        onClick={() => {
                          deleteOrganizationUser(user);
                        }}
                        icon={'delete'}
                        style={{
                          fontSize: 20,
                          marginLeft: 10,
                          cursor: 'pointer',
                          color: 'red'
                        }}
                      />
                    )}
                </Actions>
              </Td>
            </TableRow>
          ))}
        </UsersTable>
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
            overlayClick={closeRolesHandler}
            bigCloseImage={closeRolesHandler}
            bigCloseImageStyle={{
              top: '-18px',
              right: '-18px',
              background: '#000',
              borderRadius: '50%'
            }}
          >
            <Wrap newDesign={true}>
              <ModalTitle>Add user roles</ModalTitle>
              <CheckUl>
                {bountyRolesData.map((role: any, i: number) => (
                  <CheckLi key={i}>
                    <Check
                      checked={role.status}
                      onChange={roleChange}
                      type="checkbox"
                      name={role.name}
                      value={role.name}
                    />
                    <CheckLabel>{role.name}</CheckLabel>
                  </CheckLi>
                ))}
              </CheckUl>
              <Button
                onClick={() => submitRoles()}
                style={{ width: '100%' }}
                color={'primary'}
                text={'Add roles'}
              />
            </Wrap>
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
