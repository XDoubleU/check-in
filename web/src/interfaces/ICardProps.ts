import { type User } from "api-wrapper/types/apiTypes"

export interface ICardProps<T> {
  data: T
  user?: User
  fetchData: () => Promise<void>
}
