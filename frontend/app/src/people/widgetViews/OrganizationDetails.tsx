import React, { useCallback, useEffect, useRef, useState } from 'react';
import styled from 'styled-components';
import { useStores } from 'store';
import { Wrap } from 'components/form/style';
import { EuiGlobalToastList } from '@elastic/eui';
import { Button, IconButton } from 'components/common';
import { useIsMobile } from 'hooks/uiHooks';
import { Formik } from 'formik';
import { FormField, validator } from 'components/form/utils';
import { BountyRoles, Organization, Person } from 'store/main';
import MaterialIcon from '@material/react-material-icon';
import { userHasRole } from 'helpers';
import { Modal } from '../../components/common';
import { colors } from '../../config/colors';
import { nonWidgetConfigs } from '../utils/Constants';
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

const UsersCount = styled.h3`
    font-size: 1.3rem;
    margin-bottom: 15px;
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
`

const TableHead = styled.div`
  display: flex;
  flex-direction: row;
  padding: 10px;
  background: #D3D3D3;
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

const OrganizationDetails = (props: { close: () => void, org: Organization | undefined }) => {
    const [loading, setIsLoading] = useState<boolean>(false);
    const isMobile = useIsMobile();
    const { main, ui } = useStores();
    const [isOpen, setIsOpen] = useState<boolean>(false);
    const [isOpenRoles, setIsOpenRoles] = useState<boolean>(false);
    const [usersCount, setUsersCount] = useState(0);
    const [disableFormButtons, setDisableFormButtons] = useState(false);
    const [users, setUsers] = useState<Person[]>([]);
    const [user, setUser] = useState<Person>();
    const [userRoles, setUserRoles] = useState<any[]>([]);
    const [bountyRoles, setBountyRoles] = useState<any[]>([]);
    const [bountyRolesData, setBountyRolesData] = useState<BountyRoles[]>([]);
    const [toasts, setToasts]: any = useState([]);

    const config = nonWidgetConfigs['organizationusers'];

    const formRef = useRef(null);
    const isOrganizationAdmin = props.org?.owner_pubkey === ui.meInfo?.owner_pubkey;
    const schema = [...config.schema];

    const initValues = {
        owner_pubkey: '',
    };

    const uuid = props.org?.uuid;

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
                addToast('Error: could not delete user');
            }
        }
    };

    const getBountyRoles = useCallback(async () => {
        const roles = await main.getRoles();
        setBountyRoles(roles);

        const bountyRolesData = roles.map((role: any) => ({
            name: role.name,
            status: false
        }));
        setBountyRolesData(bountyRolesData);
    }, [main])

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

    const handleSettingsClick = async (user: any) => {
        setUser(user);
        setIsOpenRoles(true);
        getUserRoles(user);
    };

    const closeHandler = () => {
        setIsOpen(false)
    };

    const closeRolesHandler = () => {
        setIsOpenRoles(false)
    };

    const onSubmit = async (body: any) => {
        setIsLoading(true);

        body.organization = uuid;

        const res = await main.addOrganizationUser(body);
        if (res.status === 200) {
            await getOrganizationUsers();
            await getOrganizationUsersCount();
        } else {
            addToast('Error: could not add user');
        }
        closeHandler();
        setIsLoading(false);
    };

    const roleChange = (e: any) => {
        const rolesData = bountyRolesData.map((role: any) => {
            if (role.name === e.target.value) {
                role.status = !role.status
            }
            return role;
        });

        setBountyRolesData(rolesData);
    };

    const submitRoles = async () => {
        const roleData = bountyRolesData.filter((r: any) => r.status).map((role: any) => (
            {
                owner_pubkey: user?.owner_pubkey,
                organization: uuid,
                role: role.name
            }
        ));

        if (uuid && user?.owner_pubkey) {
            const res = await main.addUserRoles(roleData, uuid, user.owner_pubkey);
            if (res.status === 200) {
                await main.getUserRoles(uuid, user.owner_pubkey);
            } else {
                addToast('Error: could not add user roles');
            }
            setIsOpenRoles(false);
        }
    };

    useEffect(() => {
        getOrganizationUsers();
        getOrganizationUsersCount();
        getBountyRoles();
    }, [getOrganizationUsers, getOrganizationUsersCount, getBountyRoles]);

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
                <UsersCount>{usersCount} User{usersCount > 1 && 's'}</UsersCount>

                {(isOrganizationAdmin || userHasRole(bountyRoles, userRoles, 'ADD USER')) && (
                    <IconButton
                        width={150}
                        height={isMobile ? 36 : 48}
                        text="Add User"
                        onClick={() => setIsOpen(true)}
                    />)
                }

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
                                    {(isOrganizationAdmin || userHasRole(bountyRoles, userRoles, 'ADD ROLES')) && (
                                        <MaterialIcon
                                            onClick={() => handleSettingsClick(user)}
                                            icon={'settings'}
                                            style={{
                                                fontSize: 20,
                                                marginLeft: 10,
                                                cursor: 'pointer',
                                                color: 'green',
                                            }}
                                        />
                                    )}
                                    {(isOrganizationAdmin || userHasRole(bountyRoles, userRoles, 'DELETE USER')) && (
                                        <MaterialIcon
                                            onClick={() => {
                                                deleteOrganizationUser(user)
                                            }}
                                            icon={'delete'}
                                            style={{
                                                fontSize: 20,
                                                marginLeft: 10,
                                                cursor: 'pointer',
                                                color: 'red',
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
                            {({ setFieldTouched, handleSubmit, values, setFieldValue, errors, initialValues }: any) => {
                                return (
                                    <Wrap
                                        newDesign={true}
                                    >
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
                                )
                            }}
                        </Formik>
                    </Modal>
                )}
                {
                    isOpenRoles && (
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
                            <Wrap
                                newDesign={true}
                            >
                                <ModalTitle>Add user roles</ModalTitle>
                                <CheckUl>
                                    {

                                        bountyRolesData.map((role: any, i: number) => (
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
                                        ))
                                    }
                                </CheckUl>
                                <Button
                                    onClick={() => submitRoles()}
                                    style={{ width: '100%' }}
                                    color={'primary'}
                                    text={'Add roles'}
                                />
                            </Wrap>
                        </Modal>
                    )
                }
            </DetailsWrap>
            <EuiGlobalToastList toasts={toasts} dismissToast={removeToast} toastLifeTimeMs={5000} />
        </Container>
    );
};

export default OrganizationDetails;
