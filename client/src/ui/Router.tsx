import React from "react";
import { RouterProvider, createHashRouter } from "react-router-dom";
import { RegisterRoute } from "./features/register/RegisterRoute";

const router = createHashRouter([
    {
        path: '/',
        Component: RegisterRoute,
    }
])

export const Router = () => {
    return <RouterProvider router={router} />
};
