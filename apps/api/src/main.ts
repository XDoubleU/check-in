import { NestFactory } from "@nestjs/core"
import { AppModule } from "./app.module"
import cookieParser from "cookie-parser"
import helmet from "helmet"
import * as Sentry from "@sentry/node"
import "@sentry/tracing"
import { WsAdapter } from "@nestjs/platform-ws"
import { ValidationPipe } from "@nestjs/common"

const corsOptions = {
  credentials: true,
  origin: process.env.WEB_URL ?? ""
}

async function bootstrap(): Promise<void> {
  const app = await NestFactory.create(AppModule)

  if (process.env.API_SENTRY_DSN) {
    Sentry.init({
      dsn: process.env.API_SENTRY_DSN,
      debug: process.env.NODE_ENV === "development",
      environment: process.env.NODE_ENV ?? "unknown",
      release: process.env.RELEASE ?? "unknown",
      tracesSampleRate: 0.7,
      integrations: [
        // enable HTTP calls tracing
        new Sentry.Integrations.Http({ tracing: true }),
        // Automatically instrument Node.js libraries and frameworks
        ...Sentry.autoDiscoverNodePerformanceMonitoringIntegrations()
      ]
    })

    app.use(Sentry.Handlers.requestHandler())
    app.use(Sentry.Handlers.tracingHandler())
  }

  app.enableCors(corsOptions)
  app.use(helmet())
  app.use(cookieParser())
  app.useGlobalPipes(new ValidationPipe({ transform: true }))

  app.useWebSocketAdapter(new WsAdapter(app))

  app.enableShutdownHooks()

  await app.listen(process.env.API_PORT ?? 8000)
}
void bootstrap()
