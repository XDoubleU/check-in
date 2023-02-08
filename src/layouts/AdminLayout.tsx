import { ReactNode } from "react"
import BaseLayout from "./BaseLayout"
import { User } from "next-auth/core/types"

type AdminLayoutProps = {
  children: ReactNode,
  title: string,
  user: User
}

export default function AdminLayout({children, title, user}: AdminLayoutProps){
  return (
    <BaseLayout title={title} showLinks={true} showNav={true} user={user} >
      <h1>{title}</h1>
      <br/>

      {children}
    </BaseLayout>
  )   
}