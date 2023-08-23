import { randomUUID } from "crypto"
import { type ChangeEventHandler } from "react"
import { Form } from "react-bootstrap"
import { type FieldError, type UseFormRegisterReturn } from "react-hook-form"

interface FormInputProps<T extends string> {
  label: string
  type: string
  placeholder?: string | number
  required?: boolean
  value?: string
  onChange?: ChangeEventHandler<HTMLInputElement | HTMLTextAreaElement>
  register?: UseFormRegisterReturn<T>
  max?: string | number
  min?: string | number
  autocomplete?: string
  // eslint-disable-next-line redundant-undefined/redundant-undefined
  errors?: FieldError | undefined
}

export default function FormInput<T extends string>({
  label,
  type,
  placeholder,
  required,
  value,
  onChange,
  register,
  max,
  min,
  autocomplete,
  errors
}: FormInputProps<T>) {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  let registerOrOnChange: any

  if (register && !onChange) {
    registerOrOnChange = register
  } else if (onChange && !register) {
    registerOrOnChange = {
      onChange: onChange
    }
  }

  return (
    <Form.Group className="mb-3" controlId={randomUUID()}>
      <Form.Label>{label}</Form.Label>
      <Form.Control
        type={type}
        placeholder={placeholder}
        required={required}
        value={value}
        max={max}
        min={min}
        autoComplete={autocomplete}
        isInvalid={!!errors}
        {...registerOrOnChange}
      ></Form.Control>
      {errors && (
        <Form.Control.Feedback type="invalid">
          {errors.message}
        </Form.Control.Feedback>
      )}
    </Form.Group>
  )
}
