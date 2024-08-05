import type { StrategyKey } from "../strategy/types";
import { ConfigFieldTypes, type ConfigScheme } from "./types";



const spreadScheme: ConfigScheme = {
    fields: [
        {
            name: 'Balance',
            required: true,
            label: 'Доступный для торговли баланс',
            placeholder: 'рублей',
            type: ConfigFieldTypes.number,
            min: 0,
            max: 5000,
            htmlType: 'number',
        },
        {
            name: 'MaxSharesToHold',
            required: true,
            label: 'Максимально лотов',
            placeholder: 'штук',
            type: ConfigFieldTypes.number,
            min: 1,
            step: 1,
            htmlType: 'number',
        },
        {
            name: 'MinProfit',
            required: true,
            label: 'Минимальный профит со сделки',
            placeholder: '',
            type: ConfigFieldTypes.number,
            min: 0,
            step: 0.01,
            htmlType: 'number',
        },
        {
            name: 'StopLossAfter',
            label: 'Стоп-лосс',
            placeholder: 'цена покупки - стоп-лосс = цена продажи',
            type: ConfigFieldTypes.number,
            min: 0,
            step: 0.01,
            htmlType: 'number',
        },
    ]
};

const rosshookScheme: ConfigScheme = {
    fields: [
        {
            name: 'Balance',
            required: true,
            label: 'Доступный для торговли баланс',
            placeholder: 'рублей',
            type: ConfigFieldTypes.number,
            min: 0,
            max: 5000,
            htmlType: 'number',
        },
        {
            name: 'MaxSharesToHold',
            required: true,
            label: 'Максимально лотов',
            placeholder: 'штук',
            type: ConfigFieldTypes.number,
            min: 1,
            step: 1,
            htmlType: 'number',
        },
        {
            name: 'StopLoss',
            label: 'Стоп-лосс',
            placeholder: 'Второй минимум - этот параметр = цена выставления стоп-лосс',
            type: ConfigFieldTypes.number,
            min: 0,
            step: 0.0001, // TODO: тут стоит плясать от минимального шага цены инструмента
            htmlType: 'number',
        },
        {
            name: 'SaveProfit',
            label: 'Макс просадка у Тейк профит',
            placeholder: 'Просадка от максимума',
            type: ConfigFieldTypes.number,
            min: 0,
            step: 0.0001, // TODO: тут стоит плясать от минимального шага цены инструмента
            htmlType: 'number',
        },
    ]
}

const macdScheme: ConfigScheme = {
    fields: [
        {
            name: 'Balance',
            required: true,
            label: 'Доступный для торговли баланс',
            placeholder: 'рублей',
            type: ConfigFieldTypes.number,
            min: 0,
            max: 5000,
            htmlType: 'number',
        },
        {
            name: 'MaxSharesToHold',
            required: true,
            label: 'Максимально лотов',
            placeholder: 'штук',
            type: ConfigFieldTypes.number,
            min: 1,
            step: 1,
            htmlType: 'number',
        },
        {
            name: 'StopLossAfter',
            label: 'Стоп-лосс',
            placeholder: 'цена покупки - этот параметр = цена выставления стоп-лосс',
            type: ConfigFieldTypes.number,
            min: 0,
            step: 0.01, // TODO: тут стоит плясать от минимального шага цены инструмента
            htmlType: 'number',
        }
    ]
}

export const DEFAULT_SCHEME: ConfigScheme = {
    fields: []
}

export const schemes: Partial<Record<StrategyKey, ConfigScheme>> = {
    spread_v0: spreadScheme,
    macd: macdScheme,
    rosshook: rosshookScheme,
}