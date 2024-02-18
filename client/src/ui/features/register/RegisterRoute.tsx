import React from "react";
import { Layout } from "../../components/Layout";
import { Outlet } from "react-router-dom";

export const RegisterRoute = () => {
    // TODO: Использовать layout route вместо этого
    // https://reactrouter.com/en/main/start/concepts#layout-routes
    return <Layout>
        <Outlet />
    </Layout>;
};