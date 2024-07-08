import React, {
    FC,
    MutableRefObject,
    RefObject,
    useCallback,
    useEffect,
    useRef,
    useState,
} from "react";
import {
    ColorType,
    IChartApi,
    createChart,
    SeriesMarker,
    Time,
    MouseEventHandler,
    CreatePriceLineOptions,
    type ISeriesPrimitive,
} from "lightweight-charts";
import { useChartDimensions } from "./hooks";
import { OHLCData, OrderOperations, OrderState } from "../../../types";
import { useCandles, useOrders } from "../space/hooks";
import { useCurrentInstrument } from "../../utils/useCurrentInstrumentId";
import { Link } from "@radix-ui/themes";
import type { Share } from "../../../node/grpc/contracts/shares";
import { useStrategyActivitiesSeries } from "../strategy/useStrategyActivities";

type ChartProps = {
    containerRef: MutableRefObject<HTMLElement>;
};

type UseChart = (
    containerRef: MutableRefObject<HTMLElement>,
    instrument: Share,
    initialData?: OHLCData[]
) => [RefObject<HTMLDivElement>, ChartApi];

type DrawPriceLineParams = { price: number; title: string; color: 'buy' | 'sell' | string };

export type ChartApi = {
    setInitialPriceSeries: (initialData: OHLCData[]) => void;
    updatePriceSeries: (newItem: OHLCData) => void;
    updateMarkers: (newMarkers: SeriesMarker<Time>) => void;
    drawPriceLine: (params: DrawPriceLineParams) => () => void;
    attachPrimitive(primitive: ISeriesPrimitive<Time>): void
};

const chartTheme: Parameters<typeof createChart>[1] = {
    rightPriceScale: {
        scaleMargins: {
            top: 0.1,
            bottom: 0.1,
        },
        textColor: "#fff",
    },
    crosshair: {
        mode: 0,
    },
    layout: {
        textColor: "#fff",
        background: {
            type: ColorType.Solid,
            color: "transparent",
        },
    },
    grid: {
        vertLines: { color: "#272a2dc0" },
        horzLines: { color: "#272a2d" },
    },
};

const candleSeriesTheme = {
    upColor: '#b5d2c1c7',
    wickUpColor: "#b5d2c1c7",
    downColor: '#ff7f7fc7',
    wickDownColor: "#ff7f7fc7",

    borderVisible: true,
    borderUpColor: "#b5d2c1",
    borderDownColor: "#ff7f7f",
};

const buyLineColor = "#9ce0b8";
const sellLineColor = "#ef6060";

const legendStyle = `position: absolute; left: 12px; top: 40px; z-index: 1; font-size: 14px; font-family: sans-serif; line-height: 18px; font-weight: 300; z-index: 10;`;

const useChart: UseChart = (containerRef, instrument) => {
    const chartSize = useChartDimensions(containerRef);
    const chartRef = useRef();
    const chartApiRef = useRef<IChartApi>();
    const candlesApiRef = useRef<ReturnType<IChartApi["addCandlestickSeries"]>>();
    const [markers, setMarkers] = useState<SeriesMarker<Time>[]>([]);

    const legendRef = useRef<HTMLDivElement>();

    const updateMarkers = useCallback((newMarkers: SeriesMarker<Time>) => {
        setMarkers((markers) => [...markers, newMarkers]);
    }, []);

    const updatePriceSeries = useCallback(
        (newItem: OHLCData) => {
            if (!candlesApiRef.current) return;

            candlesApiRef.current.update(newItem);
        },
        [candlesApiRef.current]
    );

    // TODO: Вынести в useCandlesSeries
    const setInitialPriceSeries = useCallback(
        (initialData: OHLCData[]) => {
            if (!initialData || !candlesApiRef.current) return;

            candlesApiRef.current.setData(initialData);
        },
        [candlesApiRef.current]
    );

    // TODO: Вынести в useLegend
    const updateLegend: MouseEventHandler<Time> = useCallback(
        (param) => {
            let priceFormatted = "";
            if (param.time) {
                const data: OHLCData = param.seriesData.get(
                    candlesApiRef.current
                ) as any;
                const price = data.close;
                priceFormatted = price.toFixed(2);
            }

            legendRef.current.innerHTML = `${instrument?.name} (${instrument?.ticker}) <strong>${priceFormatted}</strong>`;
        },
        [instrument, candlesApiRef.current]
    );

    useEffect(() => {
        chartApiRef.current = createChart("chart-container", {
            width: chartSize.width,
            height: chartSize.height,
            ...chartTheme,
        });
        chartApiRef.current.timeScale().applyOptions({ timeVisible: true, });
        chartApiRef.current.timeScale().fitContent();
        candlesApiRef.current = chartApiRef.current.addCandlestickSeries(candleSeriesTheme);
        candlesApiRef.current.priceScale().applyOptions({
            autoScale: true,
            scaleMargins: {
                top: 0.1,
                bottom: 0.2,
            },
        });
        legendRef.current = document.querySelector("#chart-container #legend");
        legendRef.current.style.cssText = legendStyle;

        // TODO: Add volume series
    }, []);

    useEffect(() => {
        legendRef.current.innerHTML = `${instrument?.name} (${instrument?.ticker})`;

        chartApiRef.current.subscribeCrosshairMove(updateLegend);
        return () => {
            chartApiRef.current.unsubscribeCrosshairMove(updateLegend);
        };
    }, [instrument]);

    const drawPriceLine = ({ price, title, color }: DrawPriceLineParams) => {
        const _color = color == 'buy' ? buyLineColor : color == 'sell' ? sellLineColor : color
        const line: CreatePriceLineOptions = {
            price,
            title,
            color: _color,
            lineWidth: 1,
            lineStyle: 2, // LineStyle.Dashed
            axisLabelVisible: true,
        };
        const createdLine = candlesApiRef.current.createPriceLine(line);
        return () => {
            candlesApiRef.current.removePriceLine(createdLine);
        };
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

    return [
        chartRef,
        {
            updatePriceSeries,
            setInitialPriceSeries,
            updateMarkers,
            drawPriceLine,
            attachPrimitive: candlesApiRef.current?.attachPrimitive.bind(candlesApiRef.current),
        },
    ];
};

function orderToMarkerMapper(order: OrderState): SeriesMarker<Time> {
    return {
        time: order.time,
        position:
            order.operationType === OrderOperations.Buy ? "belowBar" : "aboveBar",
        shape:
            order.operationType === OrderOperations.Buy ? "arrowUp" : "arrowDown",
        color: order.operationType === OrderOperations.Buy ? "#2196F3" : "#e91e63",
        text: `${order.lotsExecuted} x ${order.price}`,
    };
}



const cacledOrders: Record<string, boolean> = {};
const Chart: FC<ChartProps> = ({ containerRef }) => {
    const [instrument, _, instrumentInfo] = useCurrentInstrument();
    const [ref, api] = useChart(containerRef, instrumentInfo);
    const { initialData, isLoading } = useCandles(api.updatePriceSeries, instrument);

    const destroyStrategiesView = useStrategyActivitiesSeries(api);

    // TODO: Вынести в фичу ордеров
    const filterOrdersForMarkers = useCallback(
        (order: OrderState) => {
            // if (order.lotsExecuted !== order.lotsRequested) return;
            if (cacledOrders[order.id]) return;
            cacledOrders[order.id] = true;

            const marker = orderToMarkerMapper(order);
            api.updateMarkers(marker);
        },
        [api]
    );

    // TODO: Вынести в фичу ордеров
    useOrders([filterOrdersForMarkers], instrument);

    useEffect(() => {
        api.setInitialPriceSeries(initialData);
    }, [initialData]);

    // TODO: Запилить лоадер
    return (
        <div id="chart-container" data-testid="chart-container" ref={ref}>
            {/* НЕ УДАЛЯТЬ!!! Требования либы графиков */}
            <div className="lw-attribution">
                <Link
                    href="https://tradingview.github.io/lightweight-charts/"
                    target="_blank"
                >
                    Powered by Lightweight Charts™
                </Link>
            </div>
            <div id="legend" />
        </div>
    );
};

export default Chart;
