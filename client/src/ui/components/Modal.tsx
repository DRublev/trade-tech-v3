import { Button, Dialog, Flex } from "@radix-ui/themes";
import type { FC } from "react";
import React, { useEffect, useState } from "react";

type Props = {
    title: string | React.ReactNode;
    description?: string | React.ReactNode;
    trigger?: React.ReactNode;
    children?: React.ReactNode;
    actions?: React.ReactNode[];
    close?: boolean;
};

export const Modal: FC<Props> = ({
    title,
    description,
    children,
    trigger,
    close,
    actions = [
        <Button variant="soft" color="gray">
            Закрыть
        </Button>,
    ],
}: Props) => {
    const [open, setOpen] = useState(false);

    useEffect(() => {
        if (open && close) {
            setOpen(false);
        }
    }, [open, close])

    return (
        <Dialog.Root open={open} onOpenChange={setOpen}>
            <Dialog.Trigger>{trigger}</Dialog.Trigger>
            <Dialog.Content>
                <Dialog.Title>{title}</Dialog.Title>
                {description && <Dialog.Description>{description}</Dialog.Description>}
                {children}
                <Flex gap="3" mt="4" justify="end">
                    {actions.map((action) => (
                        <Dialog.Close>{action}</Dialog.Close>
                    ))}
                </Flex>
            </Dialog.Content>
        </Dialog.Root>
    );
};
