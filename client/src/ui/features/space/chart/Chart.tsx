import React, { FC, MutableRefObject, RefObject, useCallback, useEffect, useRef } from "react";
import { IChartApi, createChart } from 'lightweight-charts';
import { useChartDimensions } from "./hooks";
import { OHLCData } from "../../../../types";
import { useCandles } from "../hooks";

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
}
const useChart = ({ containerRef, }: UseChartProps): [RefObject<HTMLDivElement>, ChartApi] => {
    const chartSize = useChartDimensions(containerRef);
    const chartRef = useRef();
    let chart: IChartApi;
    let candlesApi: ReturnType<IChartApi['addCandlestickSeries']>;

    useEffect(() => {
        chart = createChart("chart-container", {
            width: chartSize.width,
            height: chartSize.height,
        });
        candlesApi = chart.addCandlestickSeries();
    // TODO: Add volume series

        chart.timeScale().fitContent();

        return () => {
            chart.remove();
        }
    }, []);

    const updatePriceSeries = useCallback((newItem: OHLCData) => {
        console.log('43 Chart', 'new candle!', newItem);

        if (!candlesApi) return;
        candlesApi.update(newItem);
    }, []);

    const setInitialPriceSeries = useCallback((initialData: OHLCData[]) => {
        if (!initialData || !candlesApi) return;

        candlesApi.setData(initialData);
    }, []);

    useEffect(() => {
        if (chart) {
            chart.applyOptions({
                width: chartSize.width,
                height: chartSize.height,
            });
        }
    }, [chartSize]);

    return [chartRef, { updatePriceSeries, setInitialPriceSeries }];
};


const Chart: FC<ChartProps> = ({ containerRef }) => {
    const [ref, api] = useChart({ containerRef });
    // TODO: Прокидывать id выбранного инструмента
    const { initialData, isLoading } = useCandles(api.updatePriceSeries);

    useEffect(() => {
        api.setInitialPriceSeries(initialData)
    }, [initialData])

    // TODO: Запилить лоадер
    return <div id="chart-container" ref={ref} />
}

export default Chart;