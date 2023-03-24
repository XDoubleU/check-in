// Source: https://github.com/jwa-lab/nest-mikro-orm/blob/trx-complex/test/test.middleware.ts

import { RequestContext, type Transaction } from "@mikro-orm/core"
import { type Knex } from "@mikro-orm/postgresql"
import { Injectable, type NestMiddleware } from "@nestjs/common"
import { type Request, type Response, type NextFunction } from "express"

@Injectable()
export class ContextManager {
  private context?: Transaction<Knex.Transaction>

  public setContext(context: Transaction<Knex.Transaction>): void {
    this.context = context
  }

  public resetContext(): Transaction<Knex.Transaction> | undefined {
    const context = this.context
    delete this.context
    return context
  }

  public getContext(): Transaction<Knex.Transaction> | undefined {
    return this.context
  }
}

@Injectable()
export class TransactionContextMiddleware implements NestMiddleware {
  private readonly contextManager: ContextManager

  public constructor(contextManager: ContextManager) {
    this.contextManager = contextManager
  }

  public use(_req: Request, _res: Response, next: NextFunction): void {
    const em = RequestContext.getEntityManager()
    const ctx = this.contextManager.getContext()

    if (em && ctx) {
      em.setTransactionContext(ctx)
    }

    next()
  }
}
