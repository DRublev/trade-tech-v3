import React, { ComponentClass, FC, MutableRefObject, memo, useCallback, useMemo, useSyncExternalStore } from "react";
import { ChartCanvas, Chart as RFChart, CrossHairCursor, lastVisibleItemBasedZoomAnchor, CandlestickSeries, YAxis, XAxis, discontinuousTimeScaleProviderBuilder, withDeviceRatio, ZoomButtons, OHLCTooltip, WithRatioProps, EdgeIndicator, MouseCoordinateY, BarSeries } from "react-financial-charts";
import s from './styles.css';
import { useChartDimensions } from "./hooks";


// TODO: Тоже вынести в хук
const openCloseColor = (d: OHLCData) => d.close > d.open ? "rgba(181, 210, 193)" : "rgba(255, 127, 127)";
const volumeColor = (d: OHLCData) => d.close > d.open ? "rgba(181, 210, 193, 0.6)" : "rgba(255, 127, 127, 0.6)";

const candlesAppearance = {
    wickStroke: "#e1e1e1",
    fill: openCloseColor,
    stroke: "#e1e1e1",
    candleStrokeWidth: 1,
    widthRatio: 0.8,
    opacity: 1,
};

type Props = WithRatioProps & { parentRef: MutableRefObject<HTMLElement>; data: OHLCData[] }

const Chart: FC<Props> = ({ parentRef, ratio, data }) => {
    const chartSize = useChartDimensions(parentRef);
    const xScaleProvider = useMemo(() => discontinuousTimeScaleProviderBuilder().inputDateAccessor(
        (d: OHLCData) => d.date,
    ), []);

    const { data: scaledData, xScale, xAccessor, displayXAccessor } = useMemo(() => xScaleProvider(data || debugData), [])

    const yExtents = useCallback((d: OHLCData) => [d.high, d.low], []);
    const volumeSeries = useCallback((d: OHLCData) => d.volume, []);

    const barChartHeight = useMemo(() => chartSize.height / 4, [chartSize.height]);
    const chartHeight = useMemo(() => chartSize.height, [chartSize.height, barChartHeight]);

    const barChartOrigin = useCallback((_: number, __: number) => [0, chartHeight - barChartHeight], [chartSize]);

    return (
        <>
            <ChartCanvas
                ratio={ratio}
                height={chartSize.height}
                width={chartSize.width}
                data={scaledData}
                displayXAccessor={displayXAccessor}
                seriesName="Data"
                xScale={xScale}
                xAccessor={xAccessor}
                zoomAnchor={lastVisibleItemBasedZoomAnchor}
            >
                <RFChart id={2} height={barChartHeight} origin={barChartOrigin} yExtents={volumeSeries}>
                    <BarSeries fillStyle={volumeColor} yAccessor={volumeSeries} />
                </RFChart>
                <RFChart id={3} height={chartHeight} yExtents={yExtents} >
                    {/* TODO: Вынести в тему (сделать хук useChatTheme или useChartConfig) */}
                    <XAxis showGridLines gridLinesStrokeStyle="#5d5d5d" strokeStyle="#fff" showTicks={false} showTickLabel={false} />
                    <YAxis showGridLines gridLinesStrokeStyle="#5d5d5d" strokeStyle="#fff" tickStrokeStyle="#fff" tickLabelFill="#fff" />

                    <MouseCoordinateY displayFormat={d => d.toFixed(2)} />
                    <CandlestickSeries {...candlesAppearance} />

                    <OHLCTooltip className={s.ohlTooltipText} origin={[8, 16]} />
                    <EdgeIndicator
                        itemType="last"
                        rectWidth={40}
                        fill={openCloseColor}
                        lineStroke={openCloseColor}
                        yAccessor={d => d.close}
                    />

                    <ZoomButtons />
                </RFChart>
                <CrossHairCursor strokeStyle="#e1e1e1" />
            </ChartCanvas>
        </>
    );
};

export default memo(withDeviceRatio()(Chart as unknown as ComponentClass<Props, any>));

const debugData: OHLCData[] = [
    { date: new Date("2010-01-04"), open: 25.436282332605284, high: 25.835021381744056, low: 25.411360259406774, close: 25.710416, volume: 38409100 },
    { date: new Date("2010-01-05"), open: 25.627344939513726, high: 25.83502196495549, low: 25.452895407434543, close: 25.718722, volume: 49749600 },
    { date: new Date("2010-01-06"), open: 25.65226505944465, high: 25.81840750861228, low: 25.353210976925574, close: 25.560888, volume: 58182400 },
    { date: new Date("2010-01-07"), open: 25.444587793771767, high: 25.502739021094353, low: 25.079077898061875, close: 25.295062, volume: 50559700 },
    { date: new Date("2010-01-08"), open: 25.153841756996414, high: 25.6522649488092, low: 25.120612602739726, close: 25.46951, volume: 51197400 },
    { date: new Date("2010-01-11"), open: 25.511044730573705, high: 25.55258096597291, low: 25.02092861663475, close: 25.145534, volume: 68754700 },
    { date: new Date("2010-01-12"), open: 25.045848646491518, high: 25.253525666777517, low: 24.84647870701696, close: 24.979392, volume: 65912100 },
    { date: new Date("2010-01-13"), open: 25.13722727051071, high: 25.353211377924218, low: 24.929550244151567, close: 25.211991, volume: 51863500 },
    { date: new Date("2010-01-14"), open: 25.178761733851413, high: 25.83502196495549, low: 25.137227159471163, close: 25.718722, volume: 63228100 },
    { date: new Date("2010-01-15"), open: 25.818406945612217, high: 25.95132023748152, low: 25.51104412745638, close: 25.635652, volume: 79913200 },
    { date: new Date("2010-01-19"), open: 25.544274163987136, high: 25.95132113440514, low: 25.486124596784563, close: 25.835022, volume: 46575700 },
    { date: new Date("2010-01-20"), open: 25.59411494568944, high: 25.702108656795026, low: 25.17876090842236, close: 25.41136, volume: 54849500 },
    { date: new Date("2010-01-21"), open: 25.427975689088637, high: 25.51935191837554, low: 24.92124291902699, close: 24.92955, volume: 73086700 },
]


interface OHLCData {
    readonly close: number;
    readonly date: Date;
    readonly high: number;
    readonly low: number;
    readonly open: number;
    readonly volume: number;
}