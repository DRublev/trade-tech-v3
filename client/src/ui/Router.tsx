import React from "react";
import { HashRouter, Route, Routes } from "react-router-dom";
import { RegisterRoute } from "./features/register/RegisterRoute";
import { ProtectedHoc } from "./components/ProtectedRoute";
import { SpaceRoute } from "./features/space/SpaceRoute";


export const Router = () => {
    return (
        <HashRouter>
            <Routes>
                <Route path="/" Component={ProtectedHoc(SpaceRoute)} />
                <Route path="/register" Component={RegisterRoute} />
            </Routes>
        </HashRouter>
    )
};
