import React, { useState, useEffect, useMemo, useCallback } from 'react';
import * as Toolbar from '@radix-ui/react-toolbar';
import { Flex, Spinner } from "@radix-ui/themes";
import { MixerHorizontalIcon, PersonIcon, PlayIcon, StopIcon } from '@radix-ui/react-icons';
import style from '../../../basicStyles.css';
import { useIpcInvoke, useLogger } from '../../../hooks';
import { useCurrentInstrument } from '../../../utils/useCurrentInstrumentId';
import { useNavigate } from 'react-router-dom';
import { ConfigChangeModal } from '../../config';
import { useDispatch } from 'react-redux';
import { setCurrentAccount } from '../../auth/authSlice';
import { Toast } from '../../../components/Toast/Toast';
import s from './styles.css';
import { StrategySelector, useStrategy } from '../../strategy';

const toolBarButtonProps = {
    className: style.button,
    style: { verticalAlign: 'middle', transform: 'scale(1.6)', marginRight: '20px' },
}

const useTradeToggle = (instrumentId: string, strategy: string, logger: ReturnType<typeof useLogger>) => {
    const startTrade = useIpcInvoke<unknown, { Ok: boolean, Error?: string }>('START_TRADE');
    const stopTrade = useIpcInvoke<unknown, { Ok: boolean, Error?: string }>('STOP_TRADE');
    const isStartedReq = useIpcInvoke<unknown, { Ok: boolean }>('IS_STARTED');
    const [isStarted, setIsStarted] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState(null);

    const toggleTrade = useCallback(async () => {
        try {
            setError(null);
            setIsLoading(true);
            console.log('35 ControlsPanel', strategy);

            logger.info("Switching trade", { isStarted, strategy });
            let res: any = {};

            if (isStarted) {
                res = await stopTrade({ instrumentId, strategy });
            } else {
                res = await startTrade({ instrumentId, strategy });
            }
            if (res.Ok) {
                setIsStarted(!isStarted);
            } else {
                throw new Error(res.Error || "Не удалось запустить стратегию. Попробуйте еще раз");
            }
        } catch (e) {
            setError(e);
            if (e.message.includes("no config found for")) {
                alert("Сначала установите настройки стратегии")
            }
            logger.error("Error switching trade state " + e);
        } finally {
            setIsLoading(false);
            logger.trace("Trade switched", { isStarted });
        }
    }, [strategy, instrumentId, isStarted]);

    useEffect(() => {
        console.log('63 ControlsPanel', strategy);

    }, [strategy])

    useEffect(() => {
        if (!instrumentId) return;
        isStartedReq({ instrumentId, strategy })
            .then(res => setIsStarted(res.Ok));
    }, [instrumentId, strategy]);

    return { isStarted, isLoading, error, toggle: toggleTrade };
};

export const ControlsPanel = () => {
    const navigate = useNavigate();
    const dispatch = useDispatch();
    const logger = useLogger({ component: 'ControlsPanel' });

    const [instrument] = useCurrentInstrument();
    const [strategy] = useStrategy();
    const [selectedStrategy, setSelectedStrategy] = useState(strategy);

    const { isStarted, isLoading, error, toggle } = useTradeToggle(instrument, selectedStrategy, logger);
    const [showErrorToast, setShowErrorToast] = useState(false);
    const [showSuccessToast, setShowSuccessToast] = useState(false);
    const StartIconComponent = useMemo(() => {
        if (isLoading) return Spinner;
        if (!isStarted) return PlayIcon;
        return StopIcon;
    }, [isStarted, isLoading]);

    const onAccountClick = () => {
        logger.info("Going to accounts select screen");
        dispatch(setCurrentAccount({ account: undefined }));
        navigate('/register/select-account');
    };

    const handleTradeToggle = () => toggle().then(() => {
        setShowSuccessToast(true);
    })

    useEffect(() => {
        setShowErrorToast(!!error);
    }, [error]);

    return (
        <Toolbar.Root>
            <Toolbar.ToggleGroup type="single">
                <Flex align="center" justify="center" gap="2" p="3">
                    <Flex mr="3">
                        <StrategySelector onChange={setSelectedStrategy} />
                    </Flex>
                    <Toolbar.Button value="start" asChild onClick={handleTradeToggle} {...toolBarButtonProps}>
                        <StartIconComponent className={isStarted ? s.stopIcon : s.startIcon} />
                    </Toolbar.Button>
                    <ConfigChangeModal
                        trigger={
                            <Toolbar.Button value="change-config" asChild {...toolBarButtonProps}>
                                <MixerHorizontalIcon color="white" />
                            </Toolbar.Button>
                        }
                    />
                    <Toolbar.Button value="logout" asChild onClick={onAccountClick} {...toolBarButtonProps}>
                        <PersonIcon />
                    </Toolbar.Button>
                </Flex>
            </Toolbar.ToggleGroup>

            <Toast
                open={showErrorToast}
                setOpen={setShowErrorToast}
                description={error?.message || error}
            />
            <Toast
                type="ok"
                open={showSuccessToast}
                setOpen={setShowSuccessToast}
                title={isStarted ? 'Запустили стратегию' : 'Остановили стратегию'}
            />
        </Toolbar.Root>
    )
}