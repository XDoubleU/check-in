import { Collection, Entity, OneToMany, PrimaryKey, Property, Unique } from "@mikro-orm/core"
import { School } from "types-custom"
import { CheckInEntity } from "./checkin"

type MikroSchoolInterface = Omit<School, "checkIns"> & { checkIns: Collection<CheckInEntity> }

@Entity({ tableName: "School" })
export class SchoolEntity implements MikroSchoolInterface {
  @PrimaryKey()
  id: number

  @Property()
  @Unique()
  name: string

  @OneToMany(() => CheckInEntity, checkIn => checkIn.school, { eager: true })
  checkIns = new Collection<CheckInEntity>(this)

  constructor(name: string) {
    this.name = name
  }
}