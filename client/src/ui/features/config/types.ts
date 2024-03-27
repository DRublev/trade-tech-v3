
type NumericField = {
    type: 'number';
    max: number;
    min: number;
}

type MoneyFiled = NumericField & {
    type: 'money';
}

type BaseFieldProps = {
    name: string;
    placeholder: string;
    alt?: string;
    icon?: string;
}

type ConfigSchemeField = BaseFieldProps & (NumericField | MoneyFiled);

export type ConfigScheme = {
    fields: ConfigSchemeField[];
}