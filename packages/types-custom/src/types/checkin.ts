import { type School } from "./school"
import { type Location } from "./location"


export interface CheckIn {
  id: number
  location: Location
  school: School
  capacity: number
  createdAt: Date
}