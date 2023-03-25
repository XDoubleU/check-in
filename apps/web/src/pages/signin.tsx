import styles from "./signin.module.css"
import { Col, Form } from "react-bootstrap"
import BaseLayout from "@/layouts/BaseLayout"
import { signin } from "my-api-wrapper"
import { useRouter } from "next/router"
import { useForm, type SubmitHandler } from "react-hook-form"
import { type SignInDto } from "types-custom"
import { useAuth } from "@/contexts"
import BaseForm from "@/components/forms/BaseForm"

// eslint-disable-next-line max-lines-per-function
export default function SignIn() {
  const router = useRouter()
  const { setUser } = useAuth()

  const {
    register,
    handleSubmit,
    setError,
    formState: { errors }
  } = useForm<SignInDto>({
    defaultValues: {
      rememberMe: true
    }
  })

  const onSubmit: SubmitHandler<SignInDto> = (data) => {
    void signin(data).then((response) => {
      if (response.ok) {
        setUser(response.data)
        return router.push("/")
      }
      setError("root", {
        message: response.message ?? "Something went wrong"
      })
      return new Promise((resolve) => resolve(true))
    })
  }

  return (
    <BaseLayout title="Sign In" showLinks={true}>
      <Col md={4} style={{ margin: "auto" }}>
        <h1 className="text-center">Sign In</h1>
        <br />

        <BaseForm
          className={styles.customForm}
          onSubmit={handleSubmit(onSubmit)}
          errors={errors}
          submitBtnText="Sign In"
        >
          <Form.Group className="mb-3">
            <Form.Label>Username</Form.Label>
            <Form.Control
              type="text"
              placeholder="Username"
              required
              {...register("username")}
            ></Form.Control>
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>Password</Form.Label>
            <Form.Control
              type="password"
              placeholder="Password"
              required
              {...register("password")}
            ></Form.Control>
          </Form.Group>
          <Form.Check
            label="Remember me"
            type="checkbox"
            {...register("rememberMe")}
          ></Form.Check>
        </BaseForm>
        <br />
        <br />
      </Col>
    </BaseLayout>
  )
}
