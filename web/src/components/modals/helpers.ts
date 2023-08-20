import { type APIResponse } from "api-wrapper"
import { type FieldValues, type UseFormSetError } from "react-hook-form"

export function setErrors<T, U extends FieldValues>(response: APIResponse<T>, setError: UseFormSetError<U>) {
  if (typeof response.message === "string") {
    setError("root", {
      message: response.message ?? "Something went wrong"
    })
  } else {
    for (const field in response.message as object) {
      setError(field as never, {
        message: (response.message as never)[field] as string
      })
    }
  }
}