import { type ReactNode } from "react"
import BaseLayout from "@/layouts/BaseLayout"

interface AdminLayoutProps {
  children: ReactNode
  title: string
}

export default function AdminLayout({ children, title }: AdminLayoutProps) {
  return (
    <BaseLayout title={title} showLinks={true} showNav={true}>
      <h1>{title}</h1>
      <br />

      {children}
    </BaseLayout>
  )
}
