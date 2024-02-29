import React from "react"
import * as Popover from '@radix-ui/react-popover';

type Props = {
    children?: React.ReactNode
    triger?: React.ReactNode
};

export const PopoverWindow = ({ children, triger }) => {
    return (
        <Popover.Root  >
            <Popover.Trigger style={{ border: 'none', backgroundColor: 'inherit', color: 'white' }}>
                {triger}
            </Popover.Trigger>
            <Popover.Portal>
                <Popover.Content style={{ backgroundColor: '#18191B', color: 'white' }}>
                    {children}
                </Popover.Content>
            </Popover.Portal>
        </Popover.Root>
    );
};