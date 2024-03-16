import { AuthClient } from "../../../grpcGW/auth";
import type { AuthClient as IAuthClient } from "../../../grpcGW/auth";
import { createService } from "./utils";

export const authService = createService<IAuthClient>(AuthClient,);
