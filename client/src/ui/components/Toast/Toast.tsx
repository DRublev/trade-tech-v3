import type { FC } from "react";
import React from "react";
import * as ToastRdx from "@radix-ui/react-toast";
import { Button } from "@radix-ui/themes";

import s from "./styles.css";

type Props = {
    title?: string;
    description?: string | React.ReactNode;
    actionButtonText?: string;
    open: boolean;
    setOpen: (open: boolean) => void;

    /**
     * Определяет стиль уведомления: как ворнинг или как успешное действие
     * @default 'warn'
     */
    type?: 'warn' | 'ok';
};

const DEFAULT_WARN_PROPS = {
    title: "Упс! Возникла ошибка",
    description: "Мы получили о ней сведения и примем меры",
    actionButtonText: "Не ок",
};

const DEFAULT_OK_PROPS = {
    title: "Что то сделали",
    description: "Все прошло ок",
    actionButtonText: "Ну ок",
};

export const Toast: FC<Props> = ({
    open,
    setOpen,
    type = 'warn',
    ...props
}) => {
    const { title, description, actionButtonText } = {
        ...type == 'warn' ? DEFAULT_WARN_PROPS : DEFAULT_OK_PROPS,
        ...props,
    };
    return (
        
        <ToastRdx.Root open={open} onOpenChange={setOpen} className={s.ToastRoot} duration={type == 'ok' && 1000}>
            <ToastRdx.Title>{title}</ToastRdx.Title>
            <ToastRdx.Description className={s.ToastDescription}>
                {description}
            </ToastRdx.Description>
            <ToastRdx.Action
                className={s.ToastAction}
                altText={actionButtonText}
                asChild
            >
                <Button variant="surface" color={type == 'warn' ? "amber" : undefined}>
                    {actionButtonText}
                </Button>
            </ToastRdx.Action>
        </ToastRdx.Root>
    );
};
