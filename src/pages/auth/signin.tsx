import { FormEventHandler, useState } from "react"
import { signIn } from "next-auth/react"
import BaseLayout from "@/layouts/BaseLayout"
import styles from "./signin.module.css"
import { Col, Form } from "react-bootstrap"
import CustomButton from "@/components/CustomButton"

export default function SignIn(){
  const [userInfo, setUserInfo] = useState({ username: "", password: ""})
  const handleSubmit: FormEventHandler<HTMLFormElement> = async (event) => {
    event.preventDefault()

    await signIn("credentials", {
      email: userInfo.username,
      password: userInfo.password,
      callbackUrl: `${window.location.origin}`
    })
  }

  return (
    <BaseLayout title="Sign In" showLinks={true} >
      <Col md={4} style={{"margin": "auto"}}>
        <h1 className="text-center">Sign In</h1>
        <br/>

        <Form className={styles.customForm} onSubmit={handleSubmit}>
          <Form.Group className="mb-3">
            <Form.Label>Username</Form.Label>
            <Form.Control type="text" placeholder="Username" value={userInfo.username} onChange={({ target}) => setUserInfo({ ...userInfo, username: target.value })}></Form.Control>
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>Password</Form.Label>
            <Form.Control type="password" placeholder="Password" value={userInfo.password} onChange={({ target}) => setUserInfo({ ...userInfo, password: target.value })}></Form.Control>
          </Form.Group>
          
          <CustomButton type="submit">Sign In</CustomButton>
        </Form>
      </Col>
    </BaseLayout>
  )
}