import * as Form from "@radix-ui/react-form";
import * as Toast from "@radix-ui/react-toast";
import { Button, Card, Container, Flex, Heading, RadioGroup } from "@radix-ui/themes";
import React, {
    FormEventHandler,
    useCallback,
    useEffect,
    useState,
} from "react";
import { Navigate, useNavigate } from "react-router";
import { useAuth } from "../auth/useAuth";
import { useGetAccounts, usePruneTokens, useSetAccount } from "./hooks";

import { useAppDispatch, useAppSelector } from '../../../store';
import { RawAccount, setAccounts } from '../accounts/accountsSlice';
import { useLogger } from "../../hooks";
import s from "./styles.css";

const useAccounts = () => {
    const dispatch = useAppDispatch();
    const accounts = useAppSelector(state => state.accounts.accounts)
    const getAccounts = useGetAccounts();

    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState(null);
    const logger = useLogger({ component: 'useAccounts' });

    const load = useCallback(async () => {
        try {
            if (isLoading) return;
            setError(null);
            logger.trace('Loading accounts');
            setIsLoading(true);
            // TODO: Типизировать
            const response = await getAccounts(null);
            logger.trace('Accounts loaded', { count: response.Accounts.length });
            dispatch(setAccounts(
                response.Accounts.map((a: RawAccount) => ({ id: a.Id, name: a.Name || a.Id }))
            ));
        } catch (e) {
            setError(e);
            logger.error('Failed to load accounts ' + e);
            dispatch(setAccounts([]));
        } finally {
            setIsLoading(false);
        }
    }, []);

    useEffect(() => {
        load();
    }, []);

    return { accounts, isLoading, error };
};

export const SelectAccountForm = () => {
    const { isAuthorized, account } = useAppSelector(state => state.auth);
    const { accounts, isLoading, error } = useAccounts();
    const { selectAccount } = useAuth();
    const navigate = useNavigate();
    const pruneTokens = usePruneTokens();
    const [alertOpen, setAlertOpen] = useState(false);
    const [alert, setAlert] = useState(null);
    const logger = useLogger({ component: 'SelectAccountForm' });
    const [selectedAccount, setSelectedAccount] = useState(null)

    const handleSubmit: FormEventHandler<HTMLFormElement> = useCallback(
        async (event) => {
            event.preventDefault();
            event.stopPropagation();
            try {
                logger.info('Setting account');
                setAlertOpen(false);
                setAlert(null);


                await selectAccount(selectedAccount);
            } catch (e) {
                setAlertOpen(true);
                setAlert({
                    message: e.message || e,
                });

                logger.error('Failed to set account ' + e);
                // TODO: Сетить serverErrorMessage
            }
        },
        [selectedAccount]
    );


    useEffect(() => {
        setAlertOpen(!!error);
        setAlert({ message: error?.message || error });
    }, [error]);

    const onLogout = useCallback(async () => {
        logger.info('Logout clicked')
        await pruneTokens({});
        navigate('/register');
    }, []);

    const preventSubmit: React.MouseEventHandler<HTMLButtonElement> = (e) => {
        e.preventDefault();
        setSelectedAccount(e.currentTarget.value)
    }

    if (isAuthorized && account) {
        return <Navigate to="/" />;
    }

    return (
        <Toast.Provider>
            <Container>
                <img src="/static/images/logo.svg" className={s.logo} />
                <Card size="3" variant="ghost" className={s.card}>
                    <RadioGroup.Root value={selectedAccount}>
                        <Form.Root onSubmit={handleSubmit}>
                            <Flex direction="column" gap="3">
                                <Form.Field name="account">
                                    <Form.Label>
                                        <Heading className={s.heading}>
                                            Выберите аккаунт
                                        </Heading>
                                    </Form.Label>


                                    {accounts.map((account) => (
                                        <Flex align="center" gap="3" key={account.id}>

                                            <Form.Control
                                                required
                                                type="radio"
                                                value={account.id}
                                                asChild
                                            >
                                                <RadioGroup.Item value={account.id} onClick={preventSubmit}>
                                                    {account.name}
                                                </RadioGroup.Item>
                                            </Form.Control>
                                        </Flex>
                                    ))}
                                </Form.Field>

                                <Form.Submit asChild>
                                    <Button className={s.submitBtn}>Дальше</Button>
                                </Form.Submit>

                                <Button onClick={onLogout} color="crimson">
                                    Выйти
                                </Button>
                            </Flex>
                        </Form.Root>
                    </RadioGroup.Root>
                </Card>

                <Toast.Root
                    open={alertOpen}
                    onOpenChange={setAlertOpen}
                    className={s.ToastRoot}
                >
                    <Toast.Title>Упс! Возникла ошибка</Toast.Title>
                    <Toast.Description className={s.ToastDescription}>
                        Мы получили о ней сведения и примем меры
                        <br />
                        {alert?.message}
                    </Toast.Description>
                    <Toast.Action className={s.ToastAction} altText="Ок" asChild>
                        <Button variant="surface" color="amber">
                            Не ок
                        </Button>
                    </Toast.Action>
                </Toast.Root>
                <Toast.Viewport className={s.ToastViewport} />
            </Container>
        </Toast.Provider>
    );
};
