import { Container, Theme } from "@radix-ui/themes"
import React from "react"

export const Layout = ({ children }) => {
    return <>
        <Theme appearance='dark'>
            <Container shrink="1" grow="1">
                {children}
            </Container>
        </Theme>
    </>
}