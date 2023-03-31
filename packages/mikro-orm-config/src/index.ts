import { type Options } from "@mikro-orm/core"
import { TsMorphMetadataProvider } from "@mikro-orm/reflection"

export * from "./entities"
export * from "./seeders"

// eslint-disable-next-line @typescript-eslint/naming-convention
const config: Options = {
  type: "postgresql",
  clientUrl: process.env.DATABASE_URL ?? "",
  entities: ["./dist/src/entities/*.js"],
  entitiesTs: ["./src/entities/*.ts"],
  baseDir: __dirname + "/..",
  metadataProvider: TsMorphMetadataProvider,
  cache: {
    enabled: false
  },
  driverOptions: {
    ...(process.env.NODE_ENV === 'production' && {
      connection: { ssl: { rejectUnauthorized: false } },
    }),
  },
  schemaGenerator: {
    managementDbName: process.env.DATABASE_NAME ?? "postgres"
  },
  migrations: {
    disableForeignKeys: false,
    path: '../dist/src/migrations',
    pathTs: '../src/migrations'
  }
}

export default config