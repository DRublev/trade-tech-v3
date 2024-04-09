import React, { FC } from "react"
import * as Popover from '@radix-ui/react-popover';

type Props = {
    children?: React.ReactNode;
    trigger?: React.ReactNode;
};

export const PopoverWindow: FC<Props> = ({ children, trigger }) => {
    const container = (document.querySelector(".radix-themes") || document.body) as HTMLElement;
    return (
        <Popover.Root>
            <Popover.Trigger style={{ border: 'none', backgroundColor: 'inherit' }}>
                {trigger}
            </Popover.Trigger>
            <Popover.Portal container={container}>
                <Popover.Content style={{ backgroundColor: '#18191B', color: 'white' }}>
                    {children}
                </Popover.Content>
            </Popover.Portal>
        </Popover.Root>
    );
};