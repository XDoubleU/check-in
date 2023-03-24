/* eslint-disable @typescript-eslint/no-extraneous-class */
// Source: https://github.com/jwa-lab/nest-mikro-orm/blob/trx-complex/test/test.module.ts

import {
  type MiddlewareConsumer,
  Module,
  type NestModule
} from "@nestjs/common"
import { ContextManager, TransactionContextMiddleware } from "./test.middleware"

@Module({
  providers: [ContextManager]
})
class TestMiddlewareSubModule implements NestModule {
  public configure(consumer: MiddlewareConsumer): void {
    consumer.apply(TransactionContextMiddleware).forRoutes("*")
  }
}

@Module({
  imports: [TestMiddlewareSubModule]
})
class TestMiddlewareModule {}

@Module({
  imports: [TestMiddlewareModule]
})
export class TestModule {}
