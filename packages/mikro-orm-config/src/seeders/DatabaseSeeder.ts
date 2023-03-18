/* eslint-disable @typescript-eslint/require-await */
import { type EntityManager } from "@mikro-orm/core"
import { Seeder } from "@mikro-orm/seeder"
import { SchoolEntity } from "../entities"

export class DatabaseSeeder extends Seeder {
  public async run(em: EntityManager): Promise<void> {
    const andere = new SchoolEntity("Andere")
    em.persist(andere)
  }
}
