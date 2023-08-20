import { Form } from "react-bootstrap"
import { type UseFormRegisterReturn } from "react-hook-form"

interface TimeZoneInputProps {
  register: UseFormRegisterReturn<"timeZone">
}

export default function TimeZoneInput({ register }: TimeZoneInputProps) {
  return (
    <Form.Group
      className="mb-3"
      hidden={process.env.NEXT_PUBLIC_EDIT_TIME_ZONE !== "true"}
    >
      <Form.Label>Time zone</Form.Label>
      <Form.Select {...register}>
        {Intl.supportedValuesOf("timeZone").map((timeZone) => {
          return (
            <option key={timeZone} value={timeZone}>
              {timeZone}
            </option>
          )
        })}
      </Form.Select>
    </Form.Group>
  )
}
