
import { Container, Theme } from "@radix-ui/themes"
import React, { FC } from "react"
import OfflineBanner from './OfflineBanner/OfflineBanner';


type Props = {
    children?: React.ReactNode
};

export const Layout: FC<Props> = ({ children }) => {
    return <>
        <Theme appearance='dark'>

            <OfflineBanner />
            <Container flexShrink={"1"} flexGrow={"1"}>

                {children}
            </Container>
        </Theme>
    </>
}