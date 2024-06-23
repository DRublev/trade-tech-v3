
type NumericField = {
    type: 'number';
    max?: number;
    min?: number;
    step?: number;
}

type MoneyFiled = Omit<NumericField, 'type'> & {
    type: 'money';
}

type BaseFieldProps = {
    name: string;
    label: string;
    placeholder: string;
    htmlType: React.HTMLInputTypeAttribute;
    info?: string;
    icon?: string;
    required?: boolean;

    /**
     * Дефолтное (рекомендуемое) значение, при отсутствии оного в сохраненном конфиге
     * Скорее всего не понадобится
     */
    defaultValue?: number | string | (number | string)[];
}

type ConfigSchemeField = BaseFieldProps & (NumericField | MoneyFiled);

type FieldTypes = ConfigScheme['fields'][number]['type'];
export const ConfigFieldTypes: Record<FieldTypes, FieldTypes> = {
    number: 'number',
    money: 'money',
}

export type ConfigScheme = {
    fields: ConfigSchemeField[];
}