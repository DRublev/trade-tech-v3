import { MutableRefObject, useSyncExternalStore } from "react";

const sizeStep = 50; // px

const sizeSelector = (size: number) => (size ? Math.floor(size / sizeStep) * sizeStep : 1)

function getRefWidthSnapshot(selector = (width: number) => width, ref: MutableRefObject<HTMLElement>) {
    if (!ref.current) return selector(800);

    const { width } = ref.current.getBoundingClientRect();
    return selector(width);
}
const getRefHeightSnapshot = (selector = (height: number) => height, ref: MutableRefObject<HTMLElement>) => {
    if (!ref.current) return selector(600);

    const { height } = ref.current.getBoundingClientRect();
    return selector(height);
};

const subscribeWindowResize = (onChange: () => void) => {
    window.addEventListener('resize', onChange);
    return () => {
        window.removeEventListener('resize', onChange);
    }
};

export const useChartDimensions = (ref: MutableRefObject<HTMLElement>) => {
    const width = useSyncExternalStore(subscribeWindowResize, () => getRefWidthSnapshot(sizeSelector, ref));
    const height = useSyncExternalStore(subscribeWindowResize, () => getRefHeightSnapshot(sizeSelector, ref));

    return { width, height };
}