import { ReactNode } from "react"
import BaseLayout from "./BaseLayout"

type AdminLayoutProps = {
  children: ReactNode,
  title: string,
  isAdmin: boolean
}

export default function AdminLayout({children, title, isAdmin}: AdminLayoutProps){
  return (
    <BaseLayout title={title} showLinks={true} showNav={true} isAdmin={isAdmin} >
      {children}
    </BaseLayout>
  )   
}