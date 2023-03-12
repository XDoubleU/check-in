import { Options } from "@mikro-orm/core"
import { TsMorphMetadataProvider } from "@mikro-orm/reflection"

const config: Options = {
  type: "postgresql",
  clientUrl: process.env.DATABASE_URL,
  entities: ["./dist/src/entities/*.js"],
  entitiesTs: ["./src/entities/*.ts"],
	baseDir: __dirname + "/..",
  metadataProvider: TsMorphMetadataProvider,
  cache: {
    enabled: false
  },
  debug: true
}

export default config
export * from "./entities"