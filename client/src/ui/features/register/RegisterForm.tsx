import React, { FormEventHandler, useCallback, useState } from "react";
import * as Form from "@radix-ui/react-form";
import {
    Button,
    Callout,
    Card,
    Container,
    Flex,
    Heading,
    Link,
    TextField,
} from "@radix-ui/themes";
import { useNavigate } from "react-router-dom";
import { useRegister } from "./hooks";
import { InfoCircledIcon } from "@radix-ui/react-icons";
// TODO: Эту херь нужно вынести в отдельный компонент с хуком и рефом
import * as Toast from '@radix-ui/react-toast';

import s from "./styles.css";

export const RegisterForm = () => {
    const register = useRegister();
    const navigate = useNavigate();
    const [alertOpen, setAlertOpen] = useState(false);
    const [alert, setAlert] = useState(null);


    const handleSubmit: FormEventHandler<HTMLFormElement> = useCallback(
        async (event) => {
            event.preventDefault();
            event.stopPropagation();
            try {
                setAlertOpen(false);
                setAlert(null);

                const data = Object.fromEntries(new FormData(event.currentTarget));
                await register(data);

                navigate("select-account");
            } catch (e) {
                setAlertOpen(true);
                setAlert({
                    message: e.message || e
                });
            }
        },
        []
    );

    return (
        <Toast.Provider>
            <Container>
                <img src="/static/images/logo.svg" className={s.logo} />
                <Card size="3" variant="ghost" className={s.card}>
                    <Form.Root onSubmit={handleSubmit}>
                        <Flex direction="column" gap="3">
                            <Form.Field name="token">
                                <Flex align="baseline" justify="between" gap="5">
                                    <Form.Label>
                                        <Heading className={s.heading}>
                                            Нам нужен токен, чтобы начать работу
                                        </Heading>
                                    </Form.Label>
                                </Flex>
                                <Form.Message match="valueMissing">
                                    Введите токен доступа
                                </Form.Message>
                                <Form.Control required type="password" asChild>
                                    <TextField.Input placeholder="Токен доступа" />
                                </Form.Control>
                            </Form.Field>

                            <Callout.Root className={s.hintbox}>
                                <Callout.Icon>
                                    <InfoCircledIcon />
                                </Callout.Icon>
                                <Callout.Text>
                                    Вы&nbsp;можете взять его в&nbsp;
                                    <Link href="https://www.tinkoff.ru/invest/settings/">
                                        настройках Тинькофф Инвестиций
                                    </Link>
                                    .<br />
                                    Нужен токен с&nbsp;полным доступом, к&nbsp;отдельному счету,
                                    только для этого бота
                                </Callout.Text>
                            </Callout.Root>

                            <Form.Submit asChild>
                                <Button className={s.submitBtn}>Запомнить</Button>
                            </Form.Submit>
                        </Flex>
                    </Form.Root>
                </Card>

                <Toast.Root open={alertOpen} onOpenChange={setAlertOpen} className={s.ToastRoot}>
                    <Toast.Title>Упс! Возникла ошибка</Toast.Title>
                    <Toast.Description className={s.ToastDescription}>{alert?.message}</Toast.Description>
                    <Toast.Action className={s.ToastAction} altText="Бля бля" asChild>
                        <Button variant="surface" color="amber">Бля</Button>
                    </Toast.Action>
                </Toast.Root>
                <Toast.Viewport className={s.ToastViewport} />
            </Container>
        </Toast.Provider>
    );
};
