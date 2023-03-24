import {
  Check,
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
  @Check({ expression: "capacity >= 0" })
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
      `(select CAST(count(*) as int) from "CheckIn" ci where ci.location_id = ${alias}.id and CAST(ci.created_at AS DATE) = CAST(NOW() AS DATE))`,
    { persist: false }
  )
  public readonly checkInsToday!: number

  public constructor(name: string, capacity: number, user: UserEntity) {
    this.name = name
    this.capacity = capacity
    this.user = user
  }

  @Property({ persist: false })
  public get normalizedName(): string {
    return this.name
      .toLowerCase()
      .replace(" ", "-")
      .replace("[^A-Za-z0-9-]+", "")
  }

  @Property({ persist: false })
  public get available(): number {
    return this.capacity - this.checkInsToday
  }
}
