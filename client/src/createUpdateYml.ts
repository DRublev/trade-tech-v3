import { app } from 'electron';

import path from "path";
import fs from "fs";

export const createUpdateYml = () => {
    try {
        let yaml = '';

        yaml += "provider: generic\n"
        yaml += "url: https://github.com/DRublev/trade-tech-v3/releases\n"
        yaml += "useMultipleRangeRequest: false\n"
        yaml += "channel: latest\n"
        yaml += "updaterCacheDirName: " + app.getName()

        const update_file = [path.join(process.resourcesPath, 'app-update.yml'), yaml]
        const dev_update_file = [path.join(process.resourcesPath, 'dev-app-update.yml'), yaml]
        const chechFiles = [update_file, dev_update_file]

        for (const f of chechFiles) {
            if (!fs.existsSync(f[0])) {
                fs.writeFileSync(f[0], f[1])
            }
        }
    } catch (e) {
        console.error("Error creating update yml file", e);
    }
}