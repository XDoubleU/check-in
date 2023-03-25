import { Entity, ManyToOne, PrimaryKey, Property } from "@mikro-orm/core"
import { SchoolEntity } from "./school"
import { LocationEntity } from "./location"
import { type CheckIn } from "types-custom"

type MikroCheckInInterface = Omit<CheckIn, "location" | "school"> & {
  location: LocationEntity
  school: SchoolEntity
}

@Entity({ tableName: "CheckIn" })
export class CheckInEntity implements MikroCheckInInterface {
  @PrimaryKey()
  public id!: number

  @ManyToOne({ onDelete: "cascade" })
  public location: LocationEntity

  @ManyToOne({ default: 1, onDelete: "set default", eager: true })
  public school: SchoolEntity

  @Property()
  public capacity: number

  @Property()
  public createdAt = new Date()

  public constructor(location: LocationEntity, school: SchoolEntity) {
    this.location = location
    this.school = school
    this.capacity = location.capacity
  }
}
