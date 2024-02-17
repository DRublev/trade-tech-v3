import React from 'react';

const useCurrentTciker = () => {
    // TODO: СОздать контекст, оттуда брать
    return 'SBER';
}

const useOCHL = (): any[] => {
    const ticker = useCurrentTciker();
    // TODO: Создавать стрим с бэком (может ws?)
    // Поднимаем на стороне ноды сервак, слушаем стрим grpc
    // со стороны grps открываем стрим и отправляем все по ws на фронт

    return []
}

export const withChartData = (Component: React.FC) => (props: React.ComponentProps<typeof Component>) => {
    const ochl = useOCHL();

    return <Component {...props} />;
};
