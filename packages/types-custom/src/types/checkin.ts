import { School } from "./school"
import { Location } from "./location"


export interface CheckIn {
  id: number
  location: Location
  school: School
  capacity: number
  createdAt: Date
}