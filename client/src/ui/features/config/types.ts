
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

export type ConfigScheme = {
    fields: ConfigSchemeField[];
}