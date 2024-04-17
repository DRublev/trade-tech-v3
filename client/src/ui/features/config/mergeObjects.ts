import { ConfigFieldTypes } from "./hooks";
import type { ConfigScheme } from "./types";

export const mergeObjects = (a: Record<string, any>, b: Record<string, any>, scheme: ConfigScheme): [Record<string, any>, number] => {
    let changedFields = 0;
    const values: Record<string, any> = {};
    for (const fieldKey in a) {
        if (!(fieldKey in b) || a[fieldKey] != b[fieldKey]) {
            changedFields++;
        }
        const field = scheme.fields.find((f) => f.name === fieldKey);
        if (!field) continue;
        if (field.type === ConfigFieldTypes.number || field.type === ConfigFieldTypes.money) {
            values[fieldKey] = Number(a[fieldKey]);
            continue;
        }
        values[fieldKey] = a[fieldKey];
    }
    return [values, changedFields];
}