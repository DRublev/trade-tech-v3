import React, { FormEventHandler, useCallback } from "react";
import * as Form from '@radix-ui/react-form';
import { Button, Card, Flex, Switch, TextField } from "@radix-ui/themes";
import { ValidChannel } from "../../../types";


const useIpcInoke = (channel: ValidChannel) => {
    const invoke = useCallback((payload: unknown) => window.ipc ? window.ipc.invoke(channel, payload) : Promise.reject, []);

    return invoke;
};

export const RegisterForm = () => {
    const regitster = useIpcInoke("REGISTER");

    const handleSubmit: FormEventHandler<HTMLFormElement> = useCallback(async (event) => {
        try {
            const data = Object.fromEntries(new FormData(event.currentTarget));
            console.log("9 RegisterForm", data);
            await regitster(data);
        } catch (e) {
            console.log("22 RegisterForm", e);
            alert('Error')
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
                                <Switch defaultChecked id="is-sandbox" role="checkbox" />
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