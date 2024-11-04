import {
  type FieldError,
  type FieldErrors,
  type UseFormRegister,
  type UseFormWatch
} from "react-hook-form"
import FormInput from "./FormInput"

interface Inputs {
  username?: string
  password?: string
  repeatPassword?: string
}

interface UserInputsProps<T extends Inputs> {
  required: boolean
  register: UseFormRegister<T>
  watch: UseFormWatch<T>
  errors: FieldErrors<T>
}

export default function UserInputs<T extends Inputs>({
  required,
  register,
  watch,
  errors
}: Readonly<UserInputsProps<T>>) {
  return (
    <>
      <FormInput
        label="Username"
        type="text"
        placeholder="Username"
        required={required}
        register={register("username" as never)}
        errors={errors.username as FieldError}
      />
      <FormInput
        label="Password"
        type="password"
        placeholder="Password"
        required={required}
        autocomplete="new-password"
        register={register("password" as never)}
      />

      <FormInput
        label="Repeat password"
        type="password"
        placeholder="Repeat password"
        autocomplete="new-password"
        register={register("repeatPassword" as never, {
          validate: (val: string | undefined) => {
            if ((watch("password" as never) as unknown as string) !== val) {
              return "Your passwords do no match"
            }
            return undefined
          }
        })}
        errors={errors.repeatPassword as FieldError}
      />
    </>
  )
}
