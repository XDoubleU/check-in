import { type CheckIn } from "./checkin"


export interface School {
  id: number
  name: string
  checkIns: CheckIn[]
}