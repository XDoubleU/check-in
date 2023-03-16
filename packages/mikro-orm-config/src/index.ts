import { type Options } from "@mikro-orm/core"
import { TsMorphMetadataProvider } from "@mikro-orm/reflection"

export * from "./entities"

export const Config: Options = {
  type: "postgresql",
  clientUrl: process.env.DATABASE_URL ?? "",
  entities: ["./dist/src/entities/*.js"],
  entitiesTs: ["./src/entities/*.ts"],
	baseDir: __dirname + "/..",
  metadataProvider: TsMorphMetadataProvider,
  cache: {
    enabled: false
  }
}
