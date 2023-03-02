import { NestFactory } from "@nestjs/core"
import { AppModule } from "./app.module"
import cookieParser from "cookie-parser"

const corsOptions = {
  credentials: true,
  origin: process.env.WEB_URL
}

async function bootstrap(): Promise<void> {
  const app = await NestFactory.create(AppModule)
  app.enableCors(corsOptions)
  app.use(cookieParser())
  await app.listen(process.env.API_PORT ?? 8000)
}
void bootstrap()
