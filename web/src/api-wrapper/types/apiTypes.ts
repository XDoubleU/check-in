import { type PartialBy, type DeepRequired } from "typing-helpers"
import { type definitions } from "./schema"

export type CheckIn = DeepRequired<definitions["CheckInDto"]>
export type CreateCheckInDto = DeepRequired<definitions["CreateCheckInDto"]>
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

export type CheckInsLocationEntryRawMap = Record<
  string,
  CheckInsLocationEntryRaw
>

export type State = DeepRequired<definitions["State"]>
export type StateDto = definitions["StateDto"]

export const DATE_FORMAT = "YYYY-MM-DD"
export const TIME_FORMAT = "HH:mm"
export const FULL_FORMAT = "YYYY-MM-DD HH:mm"
