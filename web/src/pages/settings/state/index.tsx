import { useEffect } from "react"
import {
  type Role,
  type StateDto
} from "api-wrapper/types/apiTypes"
import { AuthRedirecter } from "contexts/authContext"
import ManagerLayout from "layouts/ManagerLayout"
import BaseForm from "components/forms/BaseForm"
import { Form } from "react-bootstrap"
import { useForm, type SubmitHandler } from "react-hook-form"
import { getState, updateState } from "api-wrapper"

// eslint-disable-next-line max-lines-per-function
export default function StateView() {
  const redirects = new Map<Role, string>([
    ["manager", "/settings"],
    ["default", "/settings"]
  ])

  const {
    register,
    handleSubmit,
    setError,
    setValue,
    formState: { errors }
  } = useForm<StateDto>()

  const onSubmit: SubmitHandler<StateDto> = (data) => {
    void updateState(data).then((response) => {
      if (response.ok) {
        setValue("isMaintenance", response.data?.isMaintenance)
      }
      else {
        setError("root", {
          message: (response.message as string) ?? "Something went wrong"
        })
      }
      
      return new Promise((resolve) => resolve(true))
    })
  }

  useEffect(() => {
    void getState().then((response) => {
      if (response.ok) {
        setValue("isMaintenance", response.data?.isMaintenance)
      }
    })
  }, [setValue])

  return (
    <AuthRedirecter redirects={redirects}>
      <ManagerLayout title="State">
      <BaseForm
            onSubmit={handleSubmit(onSubmit)}
            errors={errors}
            submitBtnText="Update"
          >
            <Form.Check
              label="Is maintenance enabled"
              type="checkbox"
              {...register("isMaintenance")}
            ></Form.Check>
          </BaseForm>
      </ManagerLayout>
    </AuthRedirecter>
  )
}
