/*
Создать класс, а-ля скринер, который бы сканировал все инструменты на точки входа по конкретной стратегии
Далее этот скринер создает экземпляр стратегии по инструменту, где нашел точку входа
Как только стратегия отрабатывает сигнал, она закрывается

Таких скринеров можно создать кучу: на боковик для сеточника, на высокий спред для спред стратегии, для хуков росса и тп
Скринеры можно переиспользовать, например скринер на направление тренда может юзаться разными стратегиями

*/


//@version=5
strategy(title="CM_MacD_Ult_MTF_Strategy kotet912", shorttitle="CM_Ult_MacD_MTF_Strategy kotet912", overlay=false)

// User inputs
useCurrentRes = input(true, title="Use Current Chart Resolution?")
resCustom = input.timeframe(title="Use Different Timeframe? Uncheck Box Above", defval="60")
smd = input(true, title="Show MacD & Signal Line? Also Turn Off Dots Below")
sd = input(true, title="Show Dots When MacD Crosses Signal Line?")
sh = input(true, title="Show Histogram?")
macd_colorChange = input(true, title="Change MacD Line Color-Signal Line Cross?")
hist_colorChange = input(true, title="MacD Histogram 4 Colors?")

// Time resolution setting
res = useCurrentRes ? timeframe.period : resCustom

// MACD parameters
fastLength = input.int(12, minval=1, title="Fast Length")
slowLength = input.int(26, minval=1, title="Slow Length")
signalLength = input.int(9, minval=1, title="Signal Length")

source = close

// MACD calculations
fastMA = ta.ema(source, fastLength)
slowMA = ta.ema(source, slowLength)

macd = fastMA - slowMA
signal = ta.sma(macd, signalLength)
hist = macd - signal

outMacD = request.security(syminfo.tickerid, res, macd)
outSignal = request.security(syminfo.tickerid, res, signal)
outHist = request.security(syminfo.tickerid, res, hist)

// Histogram color conditions
histA_IsUp = outHist > outHist[1] and outHist > 0
histA_IsDown = outHist < outHist[1] and outHist > 0
histB_IsDown = outHist < outHist[1] and outHist <= 0
histB_IsUp = outHist > outHist[1] and outHist <= 0

// MACD line color conditions
macd_IsAbove = outMacD >= outSignal
macd_IsBelow = outMacD < outSignal

plot_color = hist_colorChange ? (histA_IsUp ? color.aqua : histA_IsDown ? color.blue : histB_IsDown ? color.red : histB_IsUp ? color.maroon : color.yellow) : color.gray
macd_color = macd_colorChange ? (macd_IsAbove ? color.lime : color.red) : color.red
signal_color = macd_colorChange ? color.yellow : color.lime

// Plotting
plot(smd and na(outMacD) == false ? outMacD : na, title="MACD", color=macd_color, linewidth=4)
plot(smd and na(outSignal) == false ? outSignal : na, title="Signal Line", color=signal_color, style=plot.style_line, linewidth=2)
plot(sh and na(outHist) == false ? outHist : na, title="Histogram", color=plot_color, style=plot.style_histogram, linewidth=4)
plot(sd and ta.cross(outMacD, outSignal) ? outSignal : na, title="Cross", style=plot.style_circles, linewidth=4, color=macd_color)
hline(0, '0 Line', linestyle=hline.style_solid, linewidth=2, color=color.white)

// Trading strategy
if (ta.crossover(outMacD, outSignal)) 
    strategy.entry("MACD Long", strategy.long)
    strategy.close("MACD Short")

if (ta.crossunder(outMacD, outSignal))
    strategy.entry("MACD Short", strategy.short)
    strategy.close("MACD Long")
