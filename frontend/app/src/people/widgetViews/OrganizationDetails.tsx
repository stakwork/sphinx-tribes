import React, { useCallback, useEffect, useRef, useState } from 'react';
import styled from 'styled-components';
import { useStores } from 'store';
import { Wrap } from 'components/form/style';
import { Button, IconButton } from 'components/common';
import { useIsMobile } from 'hooks/uiHooks';
import { Formik } from 'formik';
import { FormField, validator } from 'components/form/utils';
import { Organization, Person } from 'store/main';
import MaterialIcon from '@material/react-material-icon';
import { Modal } from '../../components/common';
import { colors } from '../../config/colors';
import { widgetConfigs } from '../utils/Constants';
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

const OrganizationDetails = (props: { close: () => void, org: Organization | undefined }) => {
    const [loading, setIsLoading] = useState<boolean>(false);
    const isMobile = useIsMobile();
    const { main, ui } = useStores();
    const [isOpen, setIsOpen] = useState<boolean>(false);
    const [usersCount, setUsersCount] = useState(0);
    const [disableFormButtons, setDisableFormButtons] = useState(false);
    const [users, setUsers] = useState<Person[]>([]);
    const config = widgetConfigs['organizationusers'];

    const formRef = useRef(null);
    const isOrganizationAdmin = props.org?.owner_pubkey === ui.meInfo?.owner_pubkey;
    const schema = [...config.schema];

    const initValues = {
        owner_pubkey: '',
    };

    const getOrganizationUsersCount = useCallback(async () => {
        if (props.org?.uuid) {
            const count = await main.getOrganizationUsersCount(props.org?.uuid);
            setUsersCount(count);
        }
    }, [main, props.org?.uuid]);

    const getOrganizationUsers = useCallback(async () => {
        if (props.org?.uuid) {
            const users = await main.getOrganizationUsers(props.org?.uuid);
            setUsers(users);
        }
    }, [main, props.org?.uuid]);

    const closeHandler = () => {
        setIsOpen(false)
    };

    const onSubmit = async (body: any) => {
        setIsLoading(true);

        body.organization = props.org?.uuid;

        await main.addOrganizationUser(body);
        await getOrganizationUsers();
        await getOrganizationUsersCount();

        setIsLoading(false);
        closeHandler();
    };

    useEffect(() => {
        getOrganizationUsers();
        getOrganizationUsersCount();
    }, [getOrganizationUsers, getOrganizationUsersCount]);

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

                {isOrganizationAdmin && (
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
                                {isOrganizationAdmin && (
                                    <Actions>
                                        <MaterialIcon
                                            onClick={() => props.close()}
                                            icon={'settings'}
                                            style={{
                                                fontSize: 20,
                                                marginLeft: 10,
                                                cursor: 'pointer',
                                                color: 'green',
                                            }}
                                        />
                                        <MaterialIcon
                                            onClick={() => props.close()}
                                            icon={'delete'}
                                            style={{
                                                fontSize: 20,
                                                marginLeft: 10,
                                                cursor: 'pointer',
                                                color: 'red',
                                            }}
                                        />
                                    </Actions>
                                )}
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
                                        <h5>Add new user</h5>
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
            </DetailsWrap>
        </Container>
    );
};

export default OrganizationDetails;
