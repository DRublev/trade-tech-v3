import type { FC } from 'react';
import React from 'react';
import { ConfigScheme } from './types';

type Props = {
    scheme: ConfigScheme;
};

export const Form: FC<Props> = ({ scheme }: Props) => {
    return <>config form</>;
};
