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
import { InfoCircledIcon } from "@radix-ui/react-icons";

import { useRegistration } from "./hooks";
import { useLogger } from "../../hooks";
import { Toast } from "../../components/Toast/Toast";

import s from "./styles.css";

export const RegisterForm = () => {
    const [register] = useRegistration();
    const navigate = useNavigate();
    const [alertOpen, setAlertOpen] = useState(false);
    const [alert, setAlert] = useState(null);
    const logger = useLogger({ component: "RegisterForm" });

    const handleSubmit: FormEventHandler<HTMLFormElement> = useCallback(
        async (event) => {
            event.preventDefault();
            event.stopPropagation();
            try {
                logger.info("Submitting form");
                setAlertOpen(false);
                setAlert(null);

                const data = Object.fromEntries(new FormData(event.currentTarget));
                await register(data as Record<string, string>);

                navigate("select-account");
            } catch (e) {
                setAlertOpen(true);
                setAlert({
                    message: e.message || e,
                });
                logger.error("Error while submitting form " + e);
            }
        },
        []
    );

    return (
        <Container>
            <img src="trademedia://static/images/logo.svg" className={s.logo} />
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
                                <TextField.Root placeholder="Токен доступа" />
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

            <Toast
                open={alertOpen}
                setOpen={setAlertOpen}
                description={
                    <>
                        Мы получили о ней сведения и примем меры
                        <br />
                        {alert?.message}
                    </>
                }
            />
        </Container>
    );
};
