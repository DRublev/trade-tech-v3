import { AuthClient } from "../../../contracts/auth";
import type { AuthClient as IAuthClient } from "../../../contracts/auth";
import { createService } from "./utils";

export const authService = createService<IAuthClient>(AuthClient,);
