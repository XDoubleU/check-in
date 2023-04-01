import { type Options } from "@mikro-orm/core"
import sharedConfig from "./shared-config"

export * from "./entities"
export * from "./seeders"

// eslint-disable-next-line @typescript-eslint/naming-convention
const config: Options = {
  ...sharedConfig,
  driverOptions: {
    ...(process.env.NODE_ENV === "production" && {
      connection: { ssl: { ca: process.env.CA_CERT } }
    })
  },
  migrations: {
    disableForeignKeys: false,
    path: "../dist/src/migrations",
    pathTs: "../src/migrations"
  }
}

export default config
