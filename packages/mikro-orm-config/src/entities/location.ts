import { Check, Collection, Entity, OneToMany, OneToOne, PrimaryKey, Property } from "@mikro-orm/core"
import { v4 } from "uuid"
import { UserEntity } from "./user"
import { CheckInEntity } from "./checkin"
import { Location } from "types-custom"

type MikroLocationInterface = Omit<Location, "checkIns"|"userId"> & { user: UserEntity, checkIns: Collection<CheckInEntity> }

@Entity({ tableName: "Location" })
export class LocationEntity implements MikroLocationInterface {
  @PrimaryKey({ type: 'uuid' })
  id = v4()

  @Property()
  name: string

  @Property()
  @Check({ expression: "capacity >= 0" })
  capacity: number

  @OneToOne({ inversedBy: "location", serializedName: "userId", serializer: user => user.id })
  user: UserEntity

  @OneToMany(() => CheckInEntity, checkIn => checkIn.location)
  checkIns = new Collection<CheckInEntity>(this)

  constructor(name: string, capacity: number, user: UserEntity) {
    this.name = name
    this.capacity = capacity
    this.user = user
  }

  @Property({ persist: false })
  get normalizedName(): string {
    return this.name.toLowerCase()
                    .replace(" ", "-")
                    .replace("[^A-Za-z0-9\-]+", "")
  }

  @Property({ persist: false })
  get available(): number {
    const [today, tomorrow] = this.getDates()
    let checkInsToday: number = 0

    void (this.checkIns as unknown as Collection<CheckInEntity>).matching({
      where: {
        createdAt: {
          $gte: today,
          $lt: tomorrow
        }
      }
    })
    .then(result => {
      checkInsToday = result.length
    })

    return this.capacity - checkInsToday
  }

  private getDates(): Date[] {
    const today = new Date()
    today.setHours(0)
    today.setMinutes(0)
    today.setSeconds(0)

    const tomorrow = new Date(today)
    tomorrow.setDate(tomorrow.getDate() + 1)

    return [today, tomorrow]
  }
}