import { useEffect } from "react";
import { useIpcInvoke } from "../../../ui/hooks";
import { useStrategy } from "./useStrategy";
import type { ISeriesPrimitive, ISeriesPrimitivePaneRenderer, ISeriesPrimitivePaneView, SeriesAttachedParameter, SeriesOptionsMap, Time } from "lightweight-charts";
import type { ChartApi } from "../chart/Chart";
import type { CanvasRenderingTarget2D } from "fancy-canvas";

const useStrategyActivitiesSub = () => useIpcInvoke("SUBSCRIBE_STRATEGY_ACTIVITIES");

type BaseActivityValue = {
    Text: string;
    DeleteFlag?: boolean;
}

type LevelActivity = {
    ID: string;
    Kind: "level";
    Value: BaseActivityValue & {
        Level: number;
    };
};

type PointActivity = {
    ID: string;
    Kind: "point";
    Value: BaseActivityValue & {
        X: string; // Datestring
        Y: number; // Price
    };
}

type StrategyActivity = LevelActivity | PointActivity;


export const useSubscribeStrategyActivities = () => {
    const [strategy] = useStrategy();

    const subscribeActivities = useStrategyActivitiesSub();

    const subscribe = async () => {
        try {
            await subscribeActivities({ strategy });
        } catch (e) {
            console.error('14 useStrategyActivities', e);
        }
    };

    return subscribe;

};

type UseMediaScope = {
    context: CanvasRenderingContext2D;
    mediaSize: { width: number; height: number; }
}

interface IDrawer {
    draw(activity: StrategyActivity): void;
}


const pointStyles = {
    dashWidth: 8,
    dashGap: 4,
    dashStroke: 2,
    textTopMargin: 14,
    color: 'cornflowerblue',
}

class PointRenderer implements ISeriesPrimitivePaneRenderer {
    private id: string;
    private text: string;


    private time: Time;
    private price: number;

    private timeAxisCoord = 0;
    private priceAxisCoord = 0;

    constructor(pointActivity: PointActivity, private chartApi: SeriesAttachedParameter<Time, keyof SeriesOptionsMap>) {
        this.id = pointActivity.ID;
        this.text = pointActivity.Value.Text;

        this.render = this.render.bind(this);
        this.update = this.update.bind(this);
    }

    draw(target: CanvasRenderingTarget2D): void {
        target.useMediaCoordinateSpace(this.render);
    }

    update(value: PointActivity['Value']) {
        if (value.Text != this.text) {
            this.chartApi.requestUpdate();
            return;
        }

        const ts = Math.floor(new Date(value.X).valueOf() / 1000) as any;
        const x = this.chartApi.chart.timeScale().timeToCoordinate(ts);
        this.time = ts;

        const y = this.chartApi.series.priceToCoordinate(value.Y);
        this.price = value.Y;

        if (x != this.timeAxisCoord || y != this.priceAxisCoord) {
            this.timeAxisCoord = x;
            this.priceAxisCoord = y;

            this.chartApi.requestUpdate();
        }
    }

    // eslint-disable-next-line @typescript-eslint/no-empty-function
    erase() { }


    private render(scope: UseMediaScope) {
        const ctx = scope.context;

        // save the current state of the context to the stack
        ctx.save();

        try {
            scope.context.beginPath();

            scope.context.moveTo(this.timeAxisCoord - pointStyles.dashWidth, this.priceAxisCoord);
            scope.context.lineTo(this.timeAxisCoord + pointStyles.dashWidth, this.priceAxisCoord);

            scope.context.lineWidth = pointStyles.dashStroke;
            scope.context.strokeStyle = pointStyles.color;

            scope.context.stroke();


            scope.context.font = 'bold 12px Arial';
            scope.context.fillStyle = pointStyles.color;
            scope.context.fillText(this.text, this.timeAxisCoord - 5, this.priceAxisCoord + pointStyles.textTopMargin);

            scope.context.closePath();
        } finally {
            // restore the saved context from the stack
            ctx.restore();
        }
    }
}


class PointsDrawer implements IDrawer, ISeriesPrimitive<Time> {
    private points = new Map<string, PointRenderer>();
    private latestPoints = new Map<string, PointActivity>();
    private _paneViews: ISeriesPrimitivePaneView[] = [];
    private attachParams: SeriesAttachedParameter<Time, keyof SeriesOptionsMap>;

    public hasAttached = false;

    constructor() {
        this.handleTimeRangeChange = this.handleTimeRangeChange.bind(this);
    }

    draw(activity: PointActivity): void {
        this.latestPoints.set(activity.ID, activity);

        if (this.points.has(activity.ID)) {
            const renderer = this.points.get(activity.ID);

            if (activity.Value.DeleteFlag) {
                renderer.erase();
                this.points.delete(activity.ID);
            } else {
                renderer.update(activity.Value);
            }

            return;
        }

        const drawer = new PointRenderer(activity, this.attachParams);
        this.points.set(activity.ID, drawer);

        this._paneViews.push({
            renderer() {
                return drawer;
            }
        });
        drawer.update(activity.Value);
    }
    paneViews(): ISeriesPrimitivePaneView[] {
        return this._paneViews;
    }
    attached(param: SeriesAttachedParameter<Time, keyof SeriesOptionsMap>) {
        this.hasAttached = true;
        this.attachParams = param;

        param.chart.timeScale().subscribeVisibleTimeRangeChange(this.handleTimeRangeChange);
        // TODO Еще одна подписка, где берем видимую дату / ширина канваса - паддинги-хуяддинги = ширина свечи. И от ширины свечи обновляем размеры для лейблов и точке
    }

    private handleTimeRangeChange() {
        this.points.forEach((pr, id) => {
            pr.update(this.latestPoints.get(id).Value);
        })
    }
}

class LevelsDrawer implements IDrawer {
    private erasers = new Map<string, () => void>();

    constructor(private drawLineBase: ChartApi['drawPriceLine']) { }

    draw(activity: LevelActivity) {
        if (activity.Value.DeleteFlag && this.erasers.has(activity.ID)) {
            this.erasers.get(activity.ID)();
            return;
        }
        const { Level, Text } = activity.Value;
        const eraser = this.drawLineBase({
            price: Level,
            title: Text,
            color: Text.includes('stop') ? 'sell' : 'buy',
        });

        this.erasers.set(activity.ID, eraser);
    }
}

export const useStrategyActivitiesSeries = ({ attachPrimitive, drawPriceLine }: ChartApi) => {
    const pointsDrawer = new PointsDrawer();
    const levelsDrawer = new LevelsDrawer(drawPriceLine);

    const handleActivity = (_: unknown, activity: StrategyActivity) => {
        if (activity.Kind === "level") {
            levelsDrawer.draw(activity);
            return;
        }
        if (activity.Kind === "point") {
            pointsDrawer.draw(activity);
        }
    };

    useEffect(() => {
        window.ipc.on('NEW_STRATEGY_ACTIVITY', handleActivity);
    }, []);

    useEffect(() => {
        if (attachPrimitive && !pointsDrawer.hasAttached) {
            attachPrimitive(pointsDrawer);
        }
    }, [attachPrimitive]);

    const destroy = () => {
        // TODO: Какая то хуйня на дестрой, не знаю
        window.ipc.removeListener("NEW_STRATEGY_ACTIVITY", handleActivity);
    };

    return destroy;
}