import React, { FormEventHandler, useCallback } from "react";
import * as Form from '@radix-ui/react-form';
import { Button, Card, Flex, Switch, TextField } from "@radix-ui/themes";
import { useNavigate } from "react-router-dom";
import { useRegister } from "./hooks";


export const RegisterForm = () => {
    const regitster = useRegister();
    const navigate = useNavigate();

    const handleSubmit: FormEventHandler<HTMLFormElement> = useCallback(async (event) => {
        try {
            event.preventDefault();
            event.stopPropagation();
            const data = Object.fromEntries(new FormData(event.currentTarget));
            await regitster(data);
            navigate('select-account');
        } catch (e) {
            // TODO: Выводить алерт
            console.log("22 RegisterForm", e);
        }
    }, []);

    return (
        <Card size="3">
            <Form.Root onSubmit={handleSubmit}>
                <Flex direction="column" gap="3">
                    <Form.Field name="token">
                        <Flex align="baseline" justify="between" gap="5">
                            <Form.Label>Токен</Form.Label>
                            <Form.Message match="valueMissing">Введите токен доступа</Form.Message>
                        </Flex>
                        <Form.Control required type="password" asChild>
                            <TextField.Input placeholder="Токен доступа" />
                        </Form.Control>
                    </Form.Field>

                    <Form.Field name="isSandbox">
                        <Flex align="center" gap="3">
                            <Form.Control asChild type="checkbox">
                                <Switch defaultChecked id="is-sandbox" role="checkbox" onChange={e => { e.preventDefault(); e.stopPropagation() }} />
                            </Form.Control>
                            <Form.Label htmlFor="is-sandbox">Песочница</Form.Label>
                        </Flex>
                    </Form.Field>

                    <Form.Submit asChild><Button>Запомнить</Button></Form.Submit>
                </Flex>
            </Form.Root>
        </Card>
    )
}