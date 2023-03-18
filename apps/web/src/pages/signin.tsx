import { type FormEvent, useState } from "react"
import styles from "./signin.module.css"
import { Col, Form } from "react-bootstrap"
import BaseLayout from "@/layouts/BaseLayout"
import CustomButton from "@/components/CustomButton"
import { signin } from "my-api-wrapper"
import Router from "next/router"

// TODO: implement remember me

// eslint-disable-next-line max-lines-per-function
export default function SignIn() {
  const [userInfo, setUserInfo] = useState({
    username: "",
    password: "",
    rememberMe: false
  })
  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    const response = await signin(userInfo.username, userInfo.password)
    if (response === null) {
      await Router.push("/")
    } else {
      console.log(response)
    }
  }

  return (
    <BaseLayout title="Sign In" showLinks={true}>
      <Col md={4} style={{ margin: "auto" }}>
        <h1 className="text-center">Sign In</h1>
        <br />

        <Form className={styles.customForm} onSubmit={() => handleSubmit}>
          <Form.Group className="mb-3">
            <Form.Label>Username</Form.Label>
            <Form.Control
              type="text"
              placeholder="Username"
              value={userInfo.username}
              onChange={({ target }) =>
                setUserInfo({ ...userInfo, username: target.value })
              }
              required
            ></Form.Control>
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>Password</Form.Label>
            <Form.Control
              type="password"
              placeholder="Password"
              value={userInfo.password}
              onChange={({ target }) =>
                setUserInfo({ ...userInfo, password: target.value })
              }
              required
            ></Form.Control>
          </Form.Group>

          <CustomButton type="submit">Sign In</CustomButton>
        </Form>
      </Col>
    </BaseLayout>
  )
}
