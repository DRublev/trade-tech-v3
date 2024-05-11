
import { Container, Theme } from "@radix-ui/themes"
import React, { FC } from "react"
import OfflineBanner from './OfflineBanner/OfflineBanner';
import { useResizeBasedOnContent } from "../utils/useResizeBasedOnContent";
import * as ToastRdx from "@radix-ui/react-toast";
import s from './Layout.css';

type Props = {
    children?: React.ReactNode
};

export const Layout: FC<Props> = ({ children }) => {
    useResizeBasedOnContent();

    return <>
        <Theme appearance='dark' radius="large">
            <ToastRdx.Provider>
                <OfflineBanner />
                <Container flexShrink="1" flexGrow="1">
                    {children}

                    <ToastRdx.Viewport className={s.ToastViewport} />
                </Container>
            </ToastRdx.Provider>
        </Theme>
    </>
}