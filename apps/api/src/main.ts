import { NestFactory } from "@nestjs/core"
import { AppModule } from "./app.module"
import cookieParser from "cookie-parser"
import helmet from "helmet"
import { SentryService } from "@ntegral/nestjs-sentry"

const corsOptions = {
  credentials: true,
  origin: process.env.WEB_URL ?? ""
}

async function bootstrap(): Promise<void> {
  const app = await NestFactory.create(AppModule)
  app.useLogger(SentryService.SentryServiceInstance())
  app.enableCors(corsOptions)
  app.use(helmet())
  app.use(cookieParser())
  app.enableShutdownHooks()
  await app.listen(process.env.API_PORT ?? 8000)
}
void bootstrap()
