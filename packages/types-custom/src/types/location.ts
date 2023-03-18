import { type CheckIn } from "./checkin"

export interface Location {
  id: string
  name: string
  capacity: number
  userId: string
  checkIns: CheckIn[]

  get normalizedName(): string
  get available(): number
}
