export const retry = <T extends (...args: any) => any>(action: T, { attempts = 3, delay = 100 }) => (...args: Parameters<T>) => {
    let attempt = 1;

    const r: any = async () => {
        try {
            const res = await action(args);
            return res;
        } catch (e) {
            if (attempt > attempts) {
                throw e;
            }
            await new Promise((resolve) => setTimeout(resolve, delay));
            attempt += 1;
            return r();
        }
    }

    return r();
}