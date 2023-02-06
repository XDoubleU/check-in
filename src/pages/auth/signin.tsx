import AdminLayout from "@/layouts/AdminLayout"
import { FormEventHandler, useState } from "react"
import { signIn } from "next-auth/react"
import FormItem from "@/components/FormInput"

export default function SignIn(){
  const [userInfo, setUserInfo] = useState({ username: "", password: ""})
  const handleSubmit: FormEventHandler<HTMLFormElement> = async (event) => {
    event.preventDefault()

    const res = await signIn("credentials", {
      email: userInfo.username,
      password: userInfo.password,
      callbackUrl: `${window.location.origin}`
    })
    console.log(res)
  }

  return (
    <AdminLayout title="Sign In" isSuperUser={false} showNav={false}>
      <div className="col-md-4" style={{"margin": "auto"}}>
        <h1 className="text-center">Sign In</h1>
        <br/>
        
        <form className="custom-form" onSubmit={handleSubmit}>
          <FormItem name="Username" type="username" value={userInfo.username} onChange={({ target}) => setUserInfo({ ...userInfo, username: target.value })}/>
          <FormItem name="Password" type="password" value={userInfo.password} onChange={({ target}) => setUserInfo({ ...userInfo, password: target.value })}/>
          <button className="btn btn-custom" type="submit">Sign In</button>
        </form>
    </div>
    </AdminLayout>
  )
}