import { useEffect } from "react";
import { useIpcInvoke } from "../hooks";
import { ipcEvents } from "../../ipcEvents";

const latestBounds: { height?: number; width?: number } = {};
export const useResizeBasedOnContent = () => {
    const resize = useIpcInvoke(ipcEvents.RESIZE);
    const bodyStyle = window.getComputedStyle(document.body) as any;
    const padding =
        parseInt(bodyStyle['margin-top'], 10) +
        parseInt(bodyStyle['margin-bottom'], 10) +
        parseInt(bodyStyle['padding-top'], 10) +
        parseInt(bodyStyle['padding-bottom'], 10);

    const observer = new ResizeObserver(() => {
        const height = Math.min(900, document.body.offsetHeight + padding);
        const width = Math.min(900, document.body.offsetHeight + padding);

        if (latestBounds.height === height && latestBounds.width === width) {
            return;
        }

        resize({ height, width });
        observer.disconnect();
    });

    useEffect(() => {
        observer.observe(document.body);
        return () => {
            observer.disconnect();
        }
    }, [window.location.href]);

}