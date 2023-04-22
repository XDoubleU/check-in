import { type Options } from "@mikro-orm/core"
import { TsMorphMetadataProvider } from "@mikro-orm/reflection"

export * from "./entities"
export * from "./seeders"

// eslint-disable-next-line @typescript-eslint/naming-convention
const sharedConfig: Options = {
  type: "postgresql",
  clientUrl: process.env.DATABASE_URL ?? "",
  entities: ["./dist/entities/*.js"],
  entitiesTs: ["./src/entities/*.ts"],
  baseDir: __dirname + "/..",
  metadataProvider: TsMorphMetadataProvider,
  cache: {
    enabled: false
  },
  schemaGenerator: {
    managementDbName: process.env.DATABASE_NAME ?? "postgres"
  }
}

export default sharedConfig
