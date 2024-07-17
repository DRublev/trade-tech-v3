import { Cross2Icon, MagnifyingGlassIcon } from "@radix-ui/react-icons";
import {
    Box,
    Card,
    Flex,
    Text,
    ScrollArea,
    TextField,
    type BoxProps,
    IconButton,
} from "@radix-ui/themes";
import React, { useMemo, useState, type FC } from "react";

import { Share } from "../../../../node/grpc/contracts/shares";
import { useTodaysSchedules, useSharesFromStore } from "../hooks";
import s from "./styles.css";
import { useCurrentInstrument } from "../../../utils/useCurrentInstrumentId";
import { quantToNumber } from "../../../../utils";


const isContainsWithIgnoreCase = (value: string, term: string): boolean => {
    return value.toLocaleLowerCase().includes(term.toLocaleLowerCase());
};

type Props = BoxProps & { share: Share; isAvailable: boolean };
const ShareLine: FC<Props> = ({ share, isAvailable, ...props }) => (
    <Box {...props} className={s.shareItem} mb="2" p="3">
        <Flex key={share.ticker} justify="between">
            <span>{share.name}</span>
            {isAvailable ? (
                <Text color="gray">{quantToNumber(share.minPriceIncrement)}</Text>
            ) : (
                <Text color="red">Недоступен</Text>
            )}
        </Flex>
        <Text color="gray" size="1">
            {" "}
            {share.ticker}
        </Text>
    </Box>
);

const useExchangesTradingStatus = () => {
    const schedules = useTodaysSchedules();
    const schedulesTradingNow = useMemo<Record<string, boolean>>(() => {
        const enableMap: Record<string, boolean> = schedules.reduce(
            (acc, s) => ({
                ...acc,
                [s.exchange]:
                    s.days[0].isTradingDay &&
                    s.days[0].endTime > new Date() &&
                    s.days[0].startTime < new Date(),
            }),
            {}
        );
        enableMap["MOEX_EVENING_WEEKEND"] = enableMap["MOEX_CLOSE"];
        return enableMap;
    }, [schedules]);

    return schedulesTradingNow;
};

const useShares = (): [Share[], string, (s: string) => void] => {
    const { shares } = useSharesFromStore();
    const [term, setTerm] = useState("");
    const filteredShares = useMemo(() => {
        if (!shares || !shares.length) return [];

        const fitsSearch = (share: Share) =>
            isContainsWithIgnoreCase(share.name, term) ||
            isContainsWithIgnoreCase(share.ticker, term) ||
            isContainsWithIgnoreCase(share.figi, term) ||
            isContainsWithIgnoreCase(share.uid, term);
        return shares.filter(fitsSearch);
    }, [term, shares]);

    return [filteredShares, term, setTerm];
};

export const InstrumentSelect = () => {
    const [searchInputFocused, setSearchInputFocused] = useState(false);
    const [, setCurrentInstrument, currentInstrument] = useCurrentInstrument();

    const schedulesTradingNow = useExchangesTradingStatus();
    const [instruments, search, setInstrumentSearch] = useShares();

    const onSearchChange: React.ChangeEventHandler<HTMLInputElement> = ({
        target,
    }) => {
        setInstrumentSearch(target.value);
    };

    const open = () => setSearchInputFocused(true);
    const close = () => setSearchInputFocused(false);

    const handleInstrumentSelect = (instrument: Share) => {
        if (!instrument.uid) return;

        setCurrentInstrument(instrument.uid);
        close();
        setInstrumentSearch("");
    };


    return (
        <>
            <TextField.Root
                placeholder="Выберите инструмент"
                radius="large"
                size="3"
                onClick={open}
                onChange={onSearchChange}
                autoFocus={true}
                value={searchInputFocused ? search : currentInstrument?.ticker}
            >
                <TextField.Slot>
                    <MagnifyingGlassIcon height="16" width="16" />
                </TextField.Slot>
                <TextField.Slot>
                    <IconButton size="1" variant="ghost" onClick={close}>
                        <Cross2Icon height="14" width="14" />
                    </IconButton>
                </TextField.Slot>
            </TextField.Root>
            <Card data-focused={searchInputFocused && "on"} className={s.container}>
                <ScrollArea
                    type="auto"
                    scrollbars="vertical"
                    className={s.scrollContainer}
                >
                    <Box p="2">
                        {!instruments.length && <Text>Не нашлось инструментов по вашему запросу</Text>}
                        {instruments.map((share) => (
                            <ShareLine
                                key={share.uid}
                                share={share}
                                onClick={() => handleInstrumentSelect(share)}
                                isAvailable={schedulesTradingNow[share.exchange.toUpperCase()]}
                            />
                        ))}
                    </Box>
                </ScrollArea>
            </Card>
        </>
    );
};
