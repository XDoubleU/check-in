import { type ReactElement } from "react"
import { Alert, Form } from "react-bootstrap"
import { type FieldErrors } from "react-hook-form"
import CustomButton from "components/CustomButton"

interface BaseFormProps {
  errors: FieldErrors
  onSubmit: () => void
  submitBtnText: string
  children?: ReactElement | ReactElement[]
  className?: string
  onCancelCallback?: () => void
  submitBtnDisabled?: boolean
}

export default function BaseForm({
  className,
  children,
  errors,
  onSubmit,
  submitBtnText,
  onCancelCallback,
  submitBtnDisabled
}: BaseFormProps) {
  const floatDir = onCancelCallback ? "right" : "left"

  return (
    <Form className={className ?? ""} onSubmit={onSubmit}>
      {children}
      <br />

      {errors.root && <Alert key="danger">{errors.root.message}</Alert>}

      {onCancelCallback ? (
        <>
          <CustomButton
            type="button"
            style={{ float: "left" }}
            onClick={onCancelCallback}
          >
            Cancel
          </CustomButton>
        </>
      ) : (
        <></>
      )}

      <CustomButton
        type="submit"
        style={{ float: floatDir }}
        disabled={submitBtnDisabled ?? false}
      >
        {submitBtnText}
      </CustomButton>
    </Form>
  )
}
