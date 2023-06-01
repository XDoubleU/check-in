import { type APIResponse } from "api-wrapper"
import { type ReactElement } from "react"
import { type FieldValues, type UseFormReturn } from "react-hook-form"

export interface IModalProps<T extends FieldValues, Y> {
  children: ReactElement | ReactElement[]
  form: UseFormReturn<T>
  handler: (data: T) => Promise<APIResponse<Y>>
  fetchData: () => Promise<void>
  typeName: string
}
