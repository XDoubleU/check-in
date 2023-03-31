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
  const app = await NestFactory.create(AppModule, {
    forceCloseConnections: true
  })

  Sentry.init({
    dsn: process.env.API_SENTRY_DSN ?? "",
    debug: process.env.NODE_ENV === "development",
    environment: process.env.NODE_ENV ?? "unknown",
    release: process.env.RELEASE ?? "unknown",
    tracesSampleRate: 1.0,
    integrations: [
      // enable HTTP calls tracing
      new Sentry.Integrations.Http({ tracing: true }),
      // Automatically instrument Node.js libraries and frameworks
      ...Sentry.autoDiscoverNodePerformanceMonitoringIntegrations(),
    ],
  })

  app.use(Sentry.Handlers.requestHandler())

  app.enableCors(corsOptions)
  app.use(helmet())
  app.use(cookieParser())
  app.enableShutdownHooks()

  await app.listen(process.env.API_PORT ?? 8000)
}
void bootstrap()
