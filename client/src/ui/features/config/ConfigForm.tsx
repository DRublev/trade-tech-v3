import * as Form from "@radix-ui/react-form";
import { Button, Container, Flex, Text, TextField } from '@radix-ui/themes';
import type { FC, FormEventHandler } from 'react';
import React from 'react';
import s from './ConfigForm.css';
import { ConfigScheme } from './types';

type Props = {
    scheme: ConfigScheme;
    defaultValues?: Record<string, any>;
    onSubmit: (values: Record<string, any>) => void | Promise<void>;
};

export const ConfigForm: FC<Props> = ({ scheme, defaultValues = {}, onSubmit }: Props) => {
    const handleSubmit: FormEventHandler<HTMLFormElement> = async (e) => {
        e.preventDefault();

        const data = Object.fromEntries(new FormData(e.currentTarget));
        await onSubmit(data);
    };

    return (
        <Container>
            <Form.Root onSubmit={handleSubmit}>
                <Flex gap="3" wrap="wrap">
                    {scheme.fields.map(f => (
                        <Form.Field name={f.name} key={f.name}>
                            <Form.Label>
                                <Text>{f.label}</Text>
                            </Form.Label>
                            <Form.Control required={f.required} type={f.htmlType} asChild>
                                <TextField.Root
                                    type={f.htmlType as any}
                                    placeholder={f.placeholder}
                                    min={f.min}
                                    max={f.max}
                                    step={f.step}
                                    radius="large"
                                    defaultValue={defaultValues[f.name]}
                                />
                            </Form.Control>
                            <Form.Message match="rangeOverflow" className={s.validationMessage}>
                                Разрешенный диапазон {f.min || 0} - {f.max}
                            </Form.Message>
                        </Form.Field>
                    ))}
                </Flex>
                <Form.Submit asChild className={s.submitButton}>
                    <Button>Сохранить</Button>
                </Form.Submit>
            </Form.Root>
        </Container>
    );
};
