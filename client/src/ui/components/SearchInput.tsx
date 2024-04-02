import { MagnifyingGlassIcon } from "@radix-ui/react-icons"
import { TextField } from "@radix-ui/themes"
import React, { ChangeEventHandler } from "react"

type Props = {
    placeholder: string
    onChange: ChangeEventHandler<HTMLInputElement>
};

export const SearchInput = ({ placeholder, onChange }: Props) => {
    return (
        <TextField.Root style={{ color: 'gray', marginBottom: '10px', padding: '5px' }}>
            <TextField.Slot style={{ marginRight: '5px' }} >
                <MagnifyingGlassIcon height="16" width="16" />
            </TextField.Slot>
            <TextField.Input onChange={onChange} radius="large" placeholder={placeholder} style={{ border: 'none', backgroundColor: 'inherit' }} />
        </TextField.Root>
    )
}