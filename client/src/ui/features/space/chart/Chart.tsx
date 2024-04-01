import React, { FC, MutableRefObject, RefObject, useCallback, useEffect, useRef, useState } from 'react';
import { ColorType, IChartApi, createChart, SeriesMarker, Time, MouseEventHandler, CreatePriceLineOptions } from 'lightweight-charts';
import { useChartDimensions } from "./hooks";
import { OHLCData, OrderState } from "../../../../types";
import { useCandles, useCurrentInstrumentId, useOrders } from '../hooks';

type ChartProps = {
    containerRef: MutableRefObject<HTMLElement>;

};

type UseChartProps = {
    containerRef: MutableRefObject<HTMLElement>;
    initialData?: OHLCData[];
    instrument: string;
}

type DrawPriceLineParams = { price: number, title: string, direction: 1 | 2 };

type ChartApi = {
    setInitialPriceSeries: (initialData: OHLCData[]) => void;
    updatePriceSeries: (newItem: OHLCData) => void;
    updateMarkers: (newMarkers: SeriesMarker<Time>) => void;
    drawPriceLine: (params: DrawPriceLineParams) => () =>void;
}

const chartTheme: Parameters<typeof createChart>[1] = {
    rightPriceScale: {
        scaleMargins: {
            top: 0.1,
            bottom: 0.1,
        },
        textColor: '#fff'
    },
    crosshair: {
        mode: 0,
    },
    layout: {
        textColor: '#fff',
        background: {
            type: ColorType.Solid,
            color: 'transparent',
        },
    },
    grid: {
        vertLines: { color: '#272a2d' },
        horzLines: { color: '#272a2d' },
    },
};

const candleSeriesTheme = {
    upColor: '#b5d2c1',
    wickUpColor: '#b5d2c1',
    downColor: '#ff7f7f',
    wickDownColor: '#ff7f7f',
    borderVisible: false,
}

const buyLineColor = '#9ce0b8'
const sellLineColor = '#ef6060';

const useChart = ({ containerRef, instrument }: UseChartProps): [RefObject<HTMLDivElement>, ChartApi] => {
    const chartSize = useChartDimensions(containerRef);
    const chartRef = useRef();
    const chartApiRef = useRef<IChartApi>();
    const candlesApiRef = useRef<ReturnType<IChartApi['addCandlestickSeries']>>();
    const [markers, setMarkers] = useState<SeriesMarker<Time>[]>([]);
    
    const legend = document.createElement('div');
    legend.style = `position: absolute; left: 12px; top: 12px; z-index: 1; font-size: 14px; font-family: sans-serif; line-height: 18px; font-weight: 300;`;
    const firstRow = document.createElement('div');
    firstRow.innerHTML = instrument;
    firstRow.style.color = 'white';

    const updateMarkers = useCallback((newMarkers: SeriesMarker<Time>) => {
        setMarkers(markers => [...markers, newMarkers]);
    }, [])

    const updatePriceSeries = useCallback((newItem: OHLCData) => {
        if (!candlesApiRef.current) return;
        candlesApiRef.current.update(newItem);
    }, [candlesApiRef.current]);

    const setInitialPriceSeries = useCallback((initialData: OHLCData[]) => {
        if (!initialData || !candlesApiRef.current) return;

        candlesApiRef.current.setData(initialData);
    }, [candlesApiRef.current]);

    const updateLegend: MouseEventHandler<Time> = (param) => {
        let priceFormatted = '';
        if (param.time) {
            const data: OHLCData = param.seriesData.get(candlesApiRef.current) as any;
            const price = data.close;
            priceFormatted = price.toFixed(2);
        }
        firstRow.innerHTML = `${instrument} <strong>${priceFormatted}</strong>`;
    };

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
        const container = document.querySelector('#chart-container');
        if (container) {
            legend.appendChild(firstRow);

            container.appendChild(legend);
        }

        chartApiRef.current.subscribeCrosshairMove(updateLegend)

        // TODO: Add volume series
    }, []);

    const drawPriceLine = ({ price, title, direction }: DrawPriceLineParams) => {
        const line: CreatePriceLineOptions = {
            price,
            title,
            color: direction == 1 ? buyLineColor : sellLineColor,
            lineWidth: 2,
            lineStyle: 2, // LineStyle.Dashed
            axisLabelVisible: true,
        }
        const createdLine = candlesApiRef.current.createPriceLine(line)
        return () => {
            candlesApiRef.current.removePriceLine(createdLine)
        }
    };

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

    return [chartRef, { updatePriceSeries, setInitialPriceSeries, updateMarkers, drawPriceLine }];
};


function orderToMarkerMapper(order: OrderState): SeriesMarker<Time> {
    return {
        time: order.time,
        position: order.operationType === 1 ? 'belowBar' : 'aboveBar',
        shape: 'circle',
        color: order.operationType === 1 ? 'green' : 'red',
        text: `${order.lotsExecuted} x ${order.price}`,
        size: 2,
    }
}



const Chart: FC<ChartProps> = ({ containerRef }) => {
    const [instrument] = useCurrentInstrumentId();
    const [ref, api] = useChart({ containerRef, instrument })
    const { initialData, isLoading } = useCandles(api.updatePriceSeries, instrument);
    const [removeLinesMap, setRemoveLinesMap] = useState<Record<string, () => void>>({});

    const filterOrdersForMarkers = useCallback((order: OrderState) => {
        if (order.lotsExecuted !== order.lotsRequested) return;
        const marker = orderToMarkerMapper(order);
        api.updateMarkers(marker)
    }, []);
    const drawWaitingPositions = useCallback((order: OrderState) => {
        if (removeLinesMap[order.id]) {
            removeLinesMap[order.id]();
            setRemoveLinesMap({
                ...removeLinesMap,
                [order.id]: undefined,
            });
            return;
        }
        const removeLineFunc = api.drawPriceLine({
            price: order.price,
            title: order.price.toString(),
            direction: order.operationType == 1 ? 1 : 2,
        })
        setRemoveLinesMap({
            ...removeLinesMap,
            [order.id]: removeLineFunc,
        });
    }, []);
    useOrders([filterOrdersForMarkers, drawWaitingPositions], instrument);

    useEffect(() => {
        api.setInitialPriceSeries(initialData)
    }, [initialData])

    // TODO: Запилить лоадер
    return <div id="chart-container" ref={ref} />
}

export default Chart;