import {
  Collection,
  Entity,
  OneToMany,
  PrimaryKey,
  Property,
  Unique
} from "@mikro-orm/core"
import { type School } from "types-custom"
import { CheckInEntity } from "./checkin"

type MikroSchoolInterface = Omit<School, "checkIns"> & {
  checkIns: Collection<CheckInEntity>
}

@Entity({ tableName: "School" })
export class SchoolEntity implements MikroSchoolInterface {
  @PrimaryKey()
  public id!: number

  @Property()
  @Unique()
  public name: string

  @OneToMany(() => CheckInEntity, (checkIn) => checkIn.school)
  public checkIns = new Collection<CheckInEntity>(this)

  public constructor(name: string) {
    this.name = name
  }
}
