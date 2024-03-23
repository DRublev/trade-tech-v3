import React, { FC, MutableRefObject, RefObject, useCallback, useEffect, useRef, useState } from 'react';
import { ColorType, IChartApi, createChart, SeriesMarker, Time } from 'lightweight-charts';
import { useChartDimensions } from "./hooks";
import { OHLCData } from "../../../../types";
import { useCandles, useOrders } from '../hooks';

type ChartProps = {
    containerRef: MutableRefObject<HTMLElement>;

};

type UseChartProps = {
    containerRef: MutableRefObject<HTMLElement>;
    initialData?: OHLCData[];
}

type ChartApi = {
    setInitialPriceSeries: (initialData: OHLCData[]) => void;
    updatePriceSeries: (newItem: OHLCData) => void;
    updateMarkers: (newMarkers: SeriesMarker<Time>) => void;
}

const chartTheme = {
    rightPriceScale: {
        scaleMargins: {
            top: 0.1,
            bottom: 0.1,
        },
        textColor: '#fff'
    },
    crosshair: {
        mode: 1,
    },
    layout: {
        textColor: '#fff',
        background: {
            type: ColorType.Solid,
            color: 'transparent',
        },
    },
    grid: {
        vertLines: { color: '#5a6169' },
        horzLines: { color: '#5a6169' },
    },
};

const candleSeriesTheme = {
    upColor: '#b5d2c1',
    wickUpColor: '#b5d2c1',
    downColor: '#ff7f7f',
    wickDownColor: '#ff7f7f',
    borderVisible: false,
}

const useChart = ({ containerRef, }: UseChartProps): [RefObject<HTMLDivElement>, ChartApi] => {
    const chartSize = useChartDimensions(containerRef);
    const chartRef = useRef();
    const chartApiRef = useRef<IChartApi>();
    const candlesApiRef = useRef<ReturnType<IChartApi['addCandlestickSeries']>>();
    const [markers, setMarkers] = useState<SeriesMarker<Time>[]>([]);

    useEffect(() => {
        chartApiRef.current = createChart("chart-container", {
            width: chartSize.width,
            height: chartSize.height,
            ...chartTheme,
        });
        chartApiRef.current.timeScale().fitContent();
        candlesApiRef.current = chartApiRef.current.addCandlestickSeries(candleSeriesTheme);
        candlesApiRef.current.priceScale().applyOptions({
            autoScale: true,
            scaleMargins: {
                top: 0.1,
                bottom: 0.2,
            },
        });

        // TODO: Add volume series
    }, []);

    const updateMarkers = useCallback((newMarkers: SeriesMarker<Time>) => {
        setMarkers(markers => [...markers, newMarkers]);
    }, [])

    const updatePriceSeries = useCallback((newItem: OHLCData) => {
        console.log('43 Chart', 'new candle!', newItem);

        if (!candlesApiRef.current) return;
        candlesApiRef.current.update(newItem);
    }, [candlesApiRef.current]);

    const setInitialPriceSeries = useCallback((initialData: OHLCData[]) => {
        if (!initialData || !candlesApiRef.current) return;

        candlesApiRef.current.setData(initialData);
    }, [candlesApiRef.current]);

    useEffect(() => {
        if (chartApiRef.current) {
            chartApiRef.current.applyOptions({
                width: chartSize.width,
                height: chartSize.height,
            });
        }
    }, [chartSize.width, chartApiRef.current]);

    useEffect(() => {
        if (markers.length && candlesApiRef.current) {
            candlesApiRef.current.setMarkers(markers);
        }
    }, [markers]);

    return [chartRef, { updatePriceSeries, setInitialPriceSeries, updateMarkers }];
};


const Chart: FC<ChartProps> = ({ containerRef }) => {
    const [ref, api] = useChart({ containerRef });
    // TODO: Прокидывать id выбранного инструмента
    const { initialData, isLoading } = useCandles(api.updatePriceSeries);
    useOrders(api.updateMarkers);

    useEffect(() => {
        api.setInitialPriceSeries(initialData)
    }, [initialData])

    // TODO: Запилить лоадер
    return <div id="chart-container" ref={ref} />
}

export default Chart;