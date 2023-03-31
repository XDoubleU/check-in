import { type Options } from "@mikro-orm/core"
import sharedConfig from "./shared-config"

export * from "./entities"
export * from "./seeders"

// eslint-disable-next-line @typescript-eslint/naming-convention
const config: Options = {
  ...sharedConfig,
  migrations: {
    disableForeignKeys: false,
    path: "./dist/src/migrations",
    pathTs: "./src/migrations"
  }
}

export default config
