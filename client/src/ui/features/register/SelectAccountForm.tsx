import React, { FormEventHandler, useCallback, useEffect, useState } from "react";
import * as Form from '@radix-ui/react-form';
import { Button, Card, Flex, RadioGroup } from "@radix-ui/themes";
import { useGetAccount, useSetAccount } from "./hooks";
import { useNavigate } from "react-router";

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
            setAccounts(response.Accounts.map((a: any) => ({ id: a.Id, name: a.Name || a.Id })))
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
    }, [])

    return { accounts, isLoading }
};

export const SelectAccountForm = () => {
    const { accounts, isLoading } = useAccounts();
    const setAccount = useSetAccount();
    const navigate = useNavigate();

    const handleSubmit: FormEventHandler<HTMLFormElement> = useCallback(async (event) => {
        try {
            event.preventDefault();
            event.stopPropagation();

            const data = Object.fromEntries(new FormData(event.currentTarget));
            await setAccount({ id: data.account });
            navigate('/')

        } catch (e) {
            // TODO: Показывать алерт

            console.log("24 SelectAccountForm", e);

            // TODO: Сетить serverErrorMessage
        }
    }, []);

    return (
        <Card size="3">
            <Form.Root onSubmit={handleSubmit}>
                <Flex direction="column" gap="3">
                    <RadioGroup.Root>
                        <Form.Field name="account">
                            {accounts.map(account => (
                                <Flex align="center" gap="3">
                                    <Form.Control type="radio" value={account.id}></Form.Control>
                                    <Form.Label >{account.name}</Form.Label>
                                </Flex>
                            ))}
                        </Form.Field>
                    </RadioGroup.Root>

                    <Form.Submit asChild><Button>Дальше</Button></Form.Submit>
                </Flex>
            </Form.Root>
        </Card>
    )
}