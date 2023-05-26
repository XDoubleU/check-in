import {
  Collection,
  Entity,
  Formula,
  OneToMany,
  OneToOne,
  PrimaryKey,
  Property
} from "@mikro-orm/core"
import { v4 } from "uuid"
import { UserEntity } from "./user"
import { CheckInEntity } from "./checkin"
import { type Location } from "types-custom"
import { normalizeName } from "../helpers/normalization"

type MikroLocationInterface = Omit<Location, "checkIns" | "userId"> & {
  user: UserEntity
  checkIns: Collection<CheckInEntity>
}

@Entity({ tableName: "Location" })
export class LocationEntity implements MikroLocationInterface {
  @PrimaryKey({ type: "uuid" })
  public id = v4()

  @Property()
  public name: string

  @Property()
  public capacity: number

  @OneToOne({
    inversedBy: "location",
    serializedName: "userId",
    serializer: (user: UserEntity) => user.id
  })
  public user: UserEntity

  @OneToMany(() => CheckInEntity, (checkIn) => checkIn.location)
  public checkIns = new Collection<CheckInEntity>(this)

  @Formula(
    (alias) =>
      `(
        SELECT CAST(COUNT(*) as int) 
        FROM "CheckIn"
        WHERE location_id = ${alias}.id 
          AND DATE(created_at) = DATE(NOW())
      )`,
    { persist: false, hidden: true }
  )
  public readonly checkInsToday!: number

  @Formula(
    (alias) =>
      `(
        SELECT (EXTRACT(EPOCH FROM MAX(created_at)) * 1000)::numeric::bigint
        FROM "CheckIn"
        INNER JOIN (
          SELECT location_id, COUNT(*) AS total_checkins, MAX(capacity) AS max_capacity
          FROM "CheckIn"
          WHERE DATE(created_at) = (DATE(NOW()) - INTERVAL '1' DAY)
          GROUP BY location_id
        ) daily_stats ON "CheckIn".location_id = daily_stats.location_id
        WHERE "CheckIn".location_id = ${alias}.id
          AND DATE("CheckIn".created_at) = (DATE(NOW()) - INTERVAL '1' DAY)
          AND daily_stats.total_checkins >= daily_stats.max_capacity
      )`,
    { persist: false }
  )
  public readonly yesterdayFullAt!: string | null

  public constructor(name: string, capacity: number, user: UserEntity) {
    this.name = name
    this.capacity = capacity
    this.user = user
  }

  @Property({ persist: false })
  public get normalizedName(): string {
    return normalizeName(this.name)
  }

  @Property({ persist: false })
  public get available(): number {
    return this.capacity - this.checkInsToday
  }
}
