import React, { useCallback, useEffect, useState } from 'react';
import { useStores } from 'store';
import { EuiGlobalToastList } from '@elastic/eui';
import { Button } from 'components/common';
import { BountyRoles, Organization, PaymentHistory, Person } from 'store/main';
import MaterialIcon from '@material/react-material-icon';
import { Route, Router, Switch, useRouteMatch } from 'react-router-dom';
import { satToUsd, userHasRole } from 'helpers';
import { BountyModal } from 'people/main/bountyModal';
import history from '../../config/history';
import avatarIcon from '../../public/static/profile_avatar.svg';
import DeleteTicketModal from './DeleteModal';
import RolesModal from './organization/RolesModal';
import HistoryModal from './organization/HistoryModal';
import AddUserModal from './organization/AddUserModal';
import AddBudgetModal from './organization/AddBudgetModal';
import WithdrawBudgetModal from './organization/WithdrawBudgetModal';
import EditOrgModal from './organization/EditOrgModal';
import Users from './organization/UsersList';

import {
  ActionWrap,
  Budget,
  BudgetSmallHead,
  BudgetWrap,
  Container,
  DetailsWrap,
  Grey,
  HeadButton,
  HeadButtonWrap,
  HeadNameWrap,
  HeadWrap,
  NoBudgetText,
  NoBudgetWrap,
  OrgImg,
  OrgName,
  UserWrap,
  UsersHeadWrap,
  UsersHeader,
  ViewBudgetWrap,
  ViewBudgetTextWrap
} from './organization/style';
import AssignUserRoles from './organization/AssignUserRole';

let interval;

const OrganizationDetails = (props: {
  close: () => void;
  org: Organization | undefined;
  resetOrg: (Organization) => void;
  getOrganizations: () => Promise<void>;
}) => {
  const { main, ui } = useStores();

  const [loading, setIsLoading] = useState<boolean>(false);
  const [isOpenAddUser, setIsOpenAddUser] = useState<boolean>(false);
  const [isOpenRoles, setIsOpenRoles] = useState<boolean>(false);
  const [isOpenBudget, setIsOpenBudget] = useState<boolean>(false);
  const [isOpenWithdrawBudget, setIsOpenWithdrawBudget] = useState<boolean>(false);
  const [isOpenHistory, setIsOpenHistory] = useState<boolean>(false);
  const [isOpenEditOrg, setIsOpenEditOrg] = useState<boolean>(false);
  const [showDeleteModal, setShowDeleteModal] = useState<boolean>(false);
  const [orgBudget, setOrgBudget] = useState<number>(0);
  const [paymentsHistory, setPaymentsHistory] = useState<PaymentHistory[]>([]);
  const [disableFormButtons, setDisableFormButtons] = useState(false);
  const [users, setUsers] = useState<Person[]>([]);
  const [user, setUser] = useState<Person>();
  const [userRoles, setUserRoles] = useState<any[]>([]);
  const [toasts, setToasts]: any = useState([]);
  const [invoiceStatus, setInvoiceStatus] = useState(false);
  const { path, url } = useRouteMatch();
  const [isOpenAssignRoles, setIsOpenAssignRoles] = useState<boolean>(false);

  const isOrganizationAdmin = props.org?.owner_pubkey === ui.meInfo?.owner_pubkey;

  const editOrgDisabled =
    !isOrganizationAdmin && !userHasRole(main.bountyRoles, userRoles, 'EDIT ORGANIZATION');
  const viewReportDisabled =
    !isOrganizationAdmin && !userHasRole(main.bountyRoles, userRoles, 'VIEW REPORT');
  const addBudgetDisabled =
    !isOrganizationAdmin && !userHasRole(main.bountyRoles, userRoles, 'ADD BUDGET');
  const addWithdrawDisabled =
    !isOrganizationAdmin && !userHasRole(main.bountyRoles, userRoles, 'WITHDRAW BUDGET');

  const { org, close, getOrganizations } = props;
  const uuid = org?.uuid || '';

  function addToast(title: string, color: 'danger' | 'success') {
    setToasts([
      {
        id: `${Math.random()}`,
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
      return users;
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

  const getUserRoles = useCallback(
    async (user: any) => {
      const pubkey = user.owner_pubkey;
      if (uuid && pubkey) {
        const userRoles = await await main.getUserRoles(uuid, pubkey);
        setUserRoles(userRoles);
      }
    },
    [uuid, main]
  );

  const getOrganizationBudget = useCallback(async () => {
    if (!viewReportDisabled) {
      const organizationBudget = await main.getOrganizationBudget(uuid);
      setOrgBudget(organizationBudget.total_budget);
    }
  }, [main, uuid, viewReportDisabled]);

  const getPaymentsHistory = useCallback(async () => {
    if (!viewReportDisabled) {
      const paymentHistories = await main.getPaymentHistories(uuid, 1, 2000);
      if (Array.isArray(paymentHistories)) {
        const payments = paymentHistories.map((history: PaymentHistory) => {
          if (!history.payment_type) {
            history.payment_type = 'payment';
          }
          return history;
        });
        setPaymentsHistory(payments);
      }
    }
  }, [main, uuid, viewReportDisabled]);

  const handleSettingsClick = async (user: any) => {
    setUser(user);
    setIsOpenRoles(true);
    getUserRoles(user);
  };

  const handleDeleteClick = async (user: any) => {
    setUser(user);
    setShowDeleteModal(true);
  };

  const closeAddUserHandler = () => {
    setIsOpenAddUser(false);
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

  const closeWithdrawBudgetHandler = () => {
    setIsOpenWithdrawBudget(false);
  };

  const closeAssignRolesHandler = () => {
    setIsOpenAssignRoles(false);
  };

  const onSubmitUser = async (body: any) => {
    setIsLoading(true);

    body.org_uuid = uuid;

    const res = await main.addOrganizationUser(body);
    if (res.status === 200) {
      addToast('User added to organization successfully', 'success');
      const recentUsers = await getOrganizationUsers();
      const user = recentUsers?.filter((user: Person) => user.owner_pubkey === body.owner_pubkey);
      if (user?.length === 1) {
        setUser(user[0]);
        setIsOpenAssignRoles(true);
      }
    } else {
      addToast('Error: could not add user', 'danger');
    }
    closeAddUserHandler();
    setIsLoading(false);
  };

  const onDeleteOrg = async () => {
    const res = await main.organizationDelete(uuid);
    if (res.status === 200) {
      addToast('Deleted organization', 'success');
      if (ui.meInfo) {
        getOrganizations();
        close();
      }
    } else {
      addToast('Error: could not delete organization', 'danger');
    }
  };

  const submitRoles = async (bountyRoles: BountyRoles[]) => {
    const roleData = bountyRoles
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

  const successAction = () => {
    setInvoiceStatus(true);
    main.setBudgetInvoice('');

    // get new organization budget
    getOrganizationBudget();
    getPaymentsHistory();
  };

  const pollInvoices = useCallback(async () => {
    let i = 0;
    interval = setInterval(async () => {
      try {
        await main.pollOrgBudgetInvoices(uuid);
        getOrganizationBudget();

        const count = await main.organizationInvoiceCount(uuid);
        if (count === 0) {
          clearInterval(interval);
        }

        i++;
        if (i > 5) {
          if (interval) clearInterval(interval);
        }
      } catch (e) {
        console.warn('Poll invoices error', e);
      }
    }, 6000);
  }, []);

  useEffect(() => {
    pollInvoices();

    return () => {
      clearInterval(interval);
    };
  }, [pollInvoices]);

  async function startPolling(paymentRequest: string) {
    let i = 0;
    interval = setInterval(async () => {
      try {
        const invoiceData = await main.pollInvoice(paymentRequest);
        if (invoiceData) {
          if (invoiceData.success && invoiceData.response.settled) {
            successAction();
            clearInterval(interval);
          }
        }

        i++;
        if (i > 22) {
          if (interval) clearInterval(interval);
        }
      } catch (e) {
        console.warn('AddBudget Modal Invoice Polling Error', e);
      }
    }, 5000);
  }

  useEffect(() => {
    getOrganizationUsers();
    getOrganizationBudget();
    getPaymentsHistory();
    if (uuid && ui.meInfo) {
      getUserRoles(ui.meInfo);
    }
  }, [getOrganizationUsers, getOrganizationBudget, getPaymentsHistory, getUserRoles]);

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
          <HeadButton
            text="Edit"
            color="white"
            disabled={editOrgDisabled}
            onClick={() => setIsOpenEditOrg(true)}
            style={{ borderRadius: '5px' }}
          />
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
              <ViewBudgetTextWrap>
                <Budget>
                  {orgBudget ? orgBudget.toLocaleString() : 0} <Grey>SATS</Grey>
                </Budget>
                <Budget className="budget-small">
                  {satToUsd(orgBudget)} <Grey>USD</Grey>
                </Budget>
              </ViewBudgetTextWrap>
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
          <Button
            disabled={addWithdrawDisabled}
            text="Withdraw"
            color="withdraw"
            style={{ borderRadius: '5px' }}
            onClick={() => setIsOpenWithdrawBudget(true)}
          />
          <Button
            disabled={addBudgetDisabled}
            text="Deposit"
            color="success"
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
              disabled={editOrgDisabled}
              text="Add User"
              color="white"
              style={{
                borderRadius: '5px'
              }}
              onClick={() => setIsOpenAddUser(true)}
            />
          </HeadButtonWrap>
        </UsersHeadWrap>
        <Users
          org={org}
          handleDeleteClick={handleDeleteClick}
          handleSettingsClick={handleSettingsClick}
          userRoles={userRoles}
          users={users}
        />
      </UserWrap>
      <DetailsWrap>
        {isOpenEditOrg && (
          <EditOrgModal
            isOpen={isOpenEditOrg}
            close={() => setIsOpenEditOrg(false)}
            onDelete={onDeleteOrg}
            org={org}
            resetOrg={props.resetOrg}
            addToast={addToast}
          />
        )}
        {showDeleteModal && (
          <DeleteTicketModal
            closeModal={closeDeleteModal}
            confirmDelete={confirmDelete}
            text={'User'}
            imgUrl={user?.img}
            userDelete={true}
          />
        )}
        {isOpenAddUser && (
          <AddUserModal
            isOpen={isOpenAddUser}
            close={closeAddUserHandler}
            onSubmit={onSubmitUser}
            disableFormButtons={disableFormButtons}
            setDisableFormButtons={setDisableFormButtons}
            loading={loading}
          />
        )}
        {isOpenAssignRoles && (
          <AssignUserRoles
            close={closeAssignRolesHandler}
            isOpen={isOpenAssignRoles}
            loading={loading}
            onSubmit={submitRoles}
            user={user}
            setLoading={setIsLoading}
            addToast={addToast}
          />
        )}
        {isOpenRoles && (
          <RolesModal
            uuid={uuid}
            user={user}
            addToast={addToast}
            close={closeRolesHandler}
            isOpen={isOpenRoles}
            submitRoles={submitRoles}
          />
        )}
        {isOpenBudget && (
          <AddBudgetModal
            isOpen={isOpenBudget}
            close={closeBudgetHandler}
            uuid={uuid}
            invoiceStatus={invoiceStatus}
            startPolling={startPolling}
            setInvoiceStatus={setInvoiceStatus}
          />
        )}
        {isOpenHistory && (
          <HistoryModal
            url={url}
            paymentsHistory={paymentsHistory}
            close={closeHistoryHandler}
            isOpen={isOpenHistory}
          />
        )}
        {isOpenWithdrawBudget && (
          <WithdrawBudgetModal
            uuid={uuid}
            isOpen={isOpenWithdrawBudget}
            close={closeWithdrawBudgetHandler}
            getOrganizationBudget={getOrganizationBudget}
          />
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
