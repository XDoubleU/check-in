import {
  Entity,
  Enum,
  OneToOne,
  PrimaryKey,
  Property,
  Unique
} from "@mikro-orm/core"
import { hashSync } from "bcrypt"
import { Role, type User } from "types-custom"
import { v4 } from "uuid"
import { LocationEntity } from "./location"

type MikroUserInterface = Omit<User, "location"> & { location?: LocationEntity }

@Entity({ tableName: "User" })
export class UserEntity implements MikroUserInterface {
  @PrimaryKey({ type: "uuid" })
  public id = v4()

  @Property()
  @Unique()
  public username: string

  @Property({ hidden: true })
  public passwordHash: string

  @Enum({ default: [Role.User] })
  public roles = [Role.User]

  @OneToOne({ mappedBy: "user", eager: true })
  public location?: LocationEntity

  public constructor(username: string, password: string, role?: Role) {
    this.username = username
    this.passwordHash = hashSync(password, 12)

    if (role) {
      this.roles = [role]
    }
  }

  public update(username?: string, password?: string): void {
    this.username = username ?? this.username
    this.passwordHash = password ? hashSync(password, 12) : this.passwordHash
  }
}
