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
};

export const Toast: FC<Props> = ({
    title = "Упс! Возникла ошибка",
    description = "Мы получили о ней сведения и примем меры",
    actionButtonText = "Не ок",
    open,
    setOpen,
}) => {
    return (
        <ToastRdx.Provider>
            <ToastRdx.Root open={open} onOpenChange={setOpen} className={s.ToastRoot}>
                <ToastRdx.Title>{title}</ToastRdx.Title>
                <ToastRdx.Description className={s.ToastDescription}>
                    {description}
                </ToastRdx.Description>
                <ToastRdx.Action
                    className={s.ToastAction}
                    altText={actionButtonText}
                    asChild
                >
                    <Button variant="surface" color="amber">
                        {actionButtonText}
                    </Button>
                </ToastRdx.Action>
            </ToastRdx.Root>
            <ToastRdx.Viewport className={s.ToastViewport} />
        </ToastRdx.Provider>
    );
};
