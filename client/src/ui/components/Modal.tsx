import { Button, Dialog, Flex } from "@radix-ui/themes";
import type { FC } from "react";
import React from "react";

type Props = {
    children?: React.ReactNode;
    trigger?: React.ReactNode;
    title: string | React.ReactNode;
    description?: string | React.ReactNode;
    actions?: React.ReactNode[];
};

export const Modal: FC<Props> = ({
    title,
    description,
    children,
    trigger,
    actions = [
        <Button variant="soft" color="gray">
            Закрыть
        </Button>,
    ],
}: Props) => {
    return (
        <Dialog.Root>
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
