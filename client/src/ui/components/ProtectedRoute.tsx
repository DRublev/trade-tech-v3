import { Navigate, useLocation } from "react-router-dom"
import { useAuth } from "../features/auth/useAuth"
import React from "react"

type Props = { sandboxOnly?: boolean, children: any };

export const ProtectedRoute = ({ children, sandboxOnly }: Props) => {
    const { isAuthorized, isSandbox, isLoaded, account } = useAuth();
    const location = useLocation();

    // TODO: Loader
    if (!isLoaded) return <>Loading...</>;

    if (!isAuthorized || (sandboxOnly && !isSandbox)) return <Navigate to="/register" state={{ from: location, shouldBeSandbox: sandboxOnly }} />;
    if (!account) return <Navigate to="/register/select-account" state={{ from: location, shouldBeSandbox: sandboxOnly }} />;

    return children;
}

export const ProtectedHoc = (Component: React.FC, options?: Omit<Props, 'children'>) => (props: any) => <ProtectedRoute {...options}><Component {...props} /></ProtectedRoute>