import { MagnifyingGlassIcon } from "@radix-ui/react-icons";
import { TextField } from "@radix-ui/themes";
import React, { ChangeEventHandler } from "react";

type Props = {
    placeholder: string
    onChange: ChangeEventHandler<HTMLInputElement>
};

export const SearchInput = ({ placeholder, onChange }: Props) => {
    return (
        <TextField.Root placeholder={placeholder} radius="large" m="2" mb="3" color="gray" onChange={onChange}>
            <TextField.Slot style={{ marginRight: '5px' }} >
                <MagnifyingGlassIcon height="16" width="16" />
            </TextField.Slot>
        </TextField.Root>
    )
}