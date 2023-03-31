import { NestFactory } from "@nestjs/core"
import { AppModule } from "./app.module"
import cookieParser from "cookie-parser"
import helmet from "helmet"
import * as Sentry from "@sentry/node"

const corsOptions = {
  credentials: true,
  origin: process.env.WEB_URL ?? ""
}

async function bootstrap(): Promise<void> {
  const app = await NestFactory.create(AppModule)

  Sentry.init({
    dsn: process.env.API_SENTRY_DSN ?? "",
    debug: process.env.NODE_ENV === "development",
    environment: process.env.NODE_ENV ?? "unknown",
    tracesSampleRate: 1.0
  })

  app.enableCors(corsOptions)
  app.use(helmet())
  app.use(cookieParser())
  app.enableShutdownHooks()

  await app.listen(process.env.API_PORT ?? 8000)
}
void bootstrap()
