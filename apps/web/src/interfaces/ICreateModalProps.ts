import { type FieldValues, type UseFormReturn } from "react-hook-form"

export interface ICreateModalProps<T extends FieldValues> {
  form: UseFormReturn<T>
  fetchData: () => Promise<void>
}
