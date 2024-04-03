import React, {
    FormEventHandler,
    useCallback,
    useEffect,
    useState,
} from "react";
import * as Form from "@radix-ui/react-form";
import { Button, Card, Container, Flex, Heading, RadioGroup } from "@radix-ui/themes";
import { useGetAccount, usePruneTokens, useSetAccount } from "./hooks";
import { useNavigate } from "react-router";
import { useAuth } from "../auth/useAuth";
import * as Toast from "@radix-ui/react-toast";

import s from "./styles.css";
import { MinusCircledIcon } from "@radix-ui/react-icons";

const useAccounts = () => {
    const getAccounts = useGetAccount();
    const [accounts, setAccounts] = useState([]);
    const [isLoading, setIsLoading] = useState(false);

    const load = useCallback(async () => {
        try {
            if (isLoading) return;
            setIsLoading(true);
            // TODO: Типизировать
            const response: any = await getAccounts(null);
            setAccounts(
                response.Accounts.map((a: any) => ({ id: a.Id, name: a.Name || a.Id }))
            );
        } catch (e) {
            // TODO: Показывать алерт
            console.log("19 SelectAccountForm", e);
            setAccounts([]);
        } finally {
            setIsLoading(false);
        }
    }, []);

    useEffect(() => {
        load();
    }, []);

    return { accounts, isLoading };
};

export const SelectAccountForm = () => {
    const { accounts, isLoading } = useAccounts();
    const { setShouldUpdateAuthInfo } = useAuth();
    const setAccount = useSetAccount();
    const navigate = useNavigate();
    const pruneTokens = usePruneTokens();
    const [alertOpen, setAlertOpen] = useState(false);
    const [alert, setAlert] = useState(null);

    const handleSubmit: FormEventHandler<HTMLFormElement> = useCallback(
        async (event) => {
            event.preventDefault();
            event.stopPropagation();
            try {
                setAlertOpen(false);
                setAlert(null);

                setShouldUpdateAuthInfo();
                const data = Object.fromEntries(new FormData(event.currentTarget));
                await setAccount({ id: data.account });
                navigate("/");
            } catch (e) {
                setAlertOpen(true);
                setAlert({
                    message: e.message || e,
                });

                console.log("24 SelectAccountForm", e);

                // TODO: Сетить serverErrorMessage
            }
        },
        []
    );

    const onLogout = useCallback(async () => {
      console.log('onLogout click');
      await pruneTokens({});
      navigate('/register');
    }, [])

    return (
        <Toast.Provider>
            <Container>
                <img src="/static/images/logo.svg" className={s.logo} />
                <Card size="3" variant="ghost" className={s.card}>
                    <Form.Root onSubmit={handleSubmit}>
                        <Flex direction="column" gap="3">
                            <RadioGroup.Root>
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
                                            ></Form.Control>
                                            <Form.Label>{account.name}</Form.Label>
                                        </Flex>
                                    ))}
                                </Form.Field>
                            </RadioGroup.Root>

                            <Form.Submit asChild>
                                <Button className={s.submitBtn}>Дальше</Button>
                            </Form.Submit>
                            <Button onClick={onLogout}>
                              Выйти
                            </Button>
                        </Flex>
                    </Form.Root>
                </Card>

                <Toast.Root
                    open={alertOpen}
                    onOpenChange={setAlertOpen}
                    className={s.ToastRoot}
                >
                    <Toast.Title>Упс! Возникла ошибка</Toast.Title>
                    <Toast.Description className={s.ToastDescription}>
                        {alert?.message}
                    </Toast.Description>
                    <Toast.Action className={s.ToastAction} altText="Бля бля" asChild>
                        <Button variant="surface" color="amber">
                            Бля
                        </Button>
                    </Toast.Action>
                </Toast.Root>
                <Toast.Viewport className={s.ToastViewport} />
            </Container>
        </Toast.Provider>
    );
};
