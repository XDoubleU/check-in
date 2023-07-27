import { type definitions } from "./schema"

type DeepRequired<T> = Required<{
  [P in keyof T]-?: T[P] extends object | undefined
    ? DeepRequired<Required<T[P]>>
    : T[P]
}>

type PartialBy<T, K extends keyof T> = Omit<T, K> & Partial<Pick<T, K>>

export type CheckIn = DeepRequired<definitions["CheckIn"]>
export type CheckInDto = DeepRequired<definitions["CheckInDto"]>
export type CheckInsLocationEntryRaw = DeepRequired<
  definitions["CheckInsLocationEntryRaw"]
>
export type CreateLocationDto = DeepRequired<definitions["CreateLocationDto"]>
export type CreateUserDto = DeepRequired<definitions["CreateUserDto"]>
export type ErrorDto = DeepRequired<definitions["ErrorDto"]>
export type Location = DeepRequired<definitions["Location"]>
export type LocationUpdateEvent = DeepRequired<
  definitions["LocationUpdateEvent"]
>
export type PaginatedLocationsDto = DeepRequired<
  definitions["PaginatedLocationsDto"]
>
export type PaginatedSchoolsDto = DeepRequired<
  definitions["PaginatedSchoolsDto"]
>
export type PaginatedUsersDto = DeepRequired<definitions["PaginatedUsersDto"]>
export type Role = DeepRequired<definitions["Role"]>
export type School = DeepRequired<definitions["School"]>
export type SchoolDto = DeepRequired<definitions["SchoolDto"]>
export type SignInDto = DeepRequired<definitions["SignInDto"]>
export type SubscribeMessageDto = definitions["SubscribeMessageDto"]
export type UpdateLocationDto = definitions["UpdateLocationDto"]
export type UpdateUserDto = definitions["UpdateUserDto"]
export type User = PartialBy<DeepRequired<definitions["User"]>, "location">

export type WebSocketSubject = DeepRequired<definitions["WebSocketSubject"]>

export type CheckInsLocationEntryRawMap = {
  [x: number]: CheckInsLocationEntryRaw
}

export const API_DATE_FORMAT = "yyyy-MM-ddx"
