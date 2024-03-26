import React from "react";
import { HashRouter, Route, Routes } from "react-router-dom";
import { RegisterRoute } from "./features/register/RegisterRoute";
import { ProtectedHoc } from "./components/ProtectedRoute";
import { SpaceRoute } from "./features/space/SpaceRoute";
import { RegisterForm } from "./features/register/RegisterForm";
import { SelectAccountForm } from "./features/register/SelectAccountForm";
import OfflineBanner from './components/OfflineBanner';


export const Router = () => {
    return (
        <HashRouter>
            <React.Fragment>
                <OfflineBanner />
                <Routes>
                    <Route path="/" Component={ProtectedHoc(SpaceRoute)} />
                    <Route path="/register" Component={RegisterRoute}>
                        <Route path="select-account" Component={SelectAccountForm} />
                        <Route index Component={RegisterForm} />
                    </Route>
                </Routes>
            </React.Fragment>
        </HashRouter >
    )
};
