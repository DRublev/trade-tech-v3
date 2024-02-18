import { Navigate, useLocation } from "react-router-dom"
import { useAuth } from "../features/auth/useAuth"
import React from "react"

type Props = { sandboxOnly?: boolean, children: any }

export const ProtectedRoute = ({ children, sandboxOnly }: Props) => {
    const { isAuthorized, isSandbox, isLoaded } = useAuth()
    const location = useLocation()
    console.log("10 ProtectedRoute", isAuthorized, sandboxOnly, isSandbox);


    // TODO: Loader
    if (!isLoaded) return <>Loading...</>

    if (!isAuthorized || (sandboxOnly && !isSandbox)) return <Navigate to="/register" state={{ from: location, shouldBeSandbox: sandboxOnly }} />

    return children
}

export const ProtectedHoc = (Component: React.FC, options?: Omit<Props, 'children'>) => (props: any) => <ProtectedRoute {...options}><Component {...props} /></ProtectedRoute>