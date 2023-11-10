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
  IconWrap,
  NoBudgetText,
  NoBudgetWrap,
  OrgImg,
  OrgName,
  User,
  UserAction,
  UserDetails,
  UserImage,
  UserName,
  UserPubkey,
  UserWrap,
  UsersHeadWrap,
  UsersHeader,
  UsersList,
  ViewBudgetWrap,
  ViewBudgetTextWrap
} from './organization/style';

let interval;

const OrganizationDetails = (props: { close: () => void; org: Organization | undefined }) => {
  const [loading, setIsLoading] = useState<boolean>(false);

  const { main, ui } = useStores();
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const [isOpenRoles, setIsOpenRoles] = useState<boolean>(false);
  const [isOpenBudget, setIsOpenBudget] = useState<boolean>(false);
  const [isOpenWithdrawBudget, setIsOpenWithdrawBudget] = useState<boolean>(false);
  const [isOpenHistory, setIsOpenHistory] = useState<boolean>(false);
  const [showDeleteModal, setShowDeleteModal] = useState<boolean>(false);
  const [orgBudget, setOrgBudget] = useState<number>(0);
  const [paymentsHistory, setPaymentsHistory] = useState<PaymentHistory[]>([]);
  const [disableFormButtons, setDisableFormButtons] = useState(false);
  const [users, setUsers] = useState<Person[]>([]);
  const [user, setUser] = useState<Person>();
  const [userRoles, setUserRoles] = useState<any[]>([]);
  const [bountyRolesData, setBountyRolesData] = useState<BountyRoles[]>([]);
  const [toasts, setToasts]: any = useState([]);
  const [invoiceStatus, setInvoiceStatus] = useState(false);
  const { path, url } = useRouteMatch();

  const isOrganizationAdmin = props.org?.owner_pubkey === ui.meInfo?.owner_pubkey;

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
  const addWithdrawDisabled =
    !isOrganizationAdmin && !userHasRole(main.bountyRoles, userRoles, 'WITHDRAW BUDGET');

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
    const paymentHistories = await main.getPaymentHistories(uuid, 1, 20);
    setPaymentsHistory(paymentHistories);
  }, [main, uuid]);

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

  const closeWithdrawBudgetHandler = () => {
    setIsOpenWithdrawBudget(false);
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

  const successAction = () => {
    addToast('Budget was added successfully', 'success');
    closeBudgetHandler();

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
        getPaymentsHistory();

        const count = await main.organizationInvoiceCount(uuid);
        if (count === 0) {
          clearInterval(interval);
        }

        i++;
        if (i > 10) {
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
    getBountyRoles();
    getOrganizationBudget();
    getPaymentsHistory();
  }, [getOrganizationUsers, getBountyRoles, getOrganizationBudget, getPaymentsHistory]);

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
                <ViewBudgetTextWrap>
                  <Budget>
                    {orgBudget.toLocaleString()} <Grey>SATS</Grey>
                  </Budget>
                  <Budget className="budget-small">{satToUsd(orgBudget)} <Grey>USD</Grey></Budget>
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
          <UsersHeader>USERS</UsersHeader>
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
          <AddUserModal
            isOpen={isOpen}
            close={closeHandler}
            onSubmit={onSubmit}
            disableFormButtons={disableFormButtons}
            setDisableFormButtons={setDisableFormButtons}
            loading={loading}
          />
        )}
        {isOpenRoles && (
          <RolesModal
            userRoles={userRoles}
            bountyRolesData={bountyRolesData}
            uuid={uuid}
            user={user}
            addToast={addToast}
            close={closeRolesHandler}
            isOpen={isOpenRoles}
            roleChange={roleChange}
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
