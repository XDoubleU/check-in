import { type Options } from "@mikro-orm/core"
import sharedConfig from "./shared-config"

// eslint-disable-next-line @typescript-eslint/naming-convention
const config: Options = {
  ...sharedConfig
}

export default config
