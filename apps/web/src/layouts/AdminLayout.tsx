import { type ReactNode } from "react"
import BaseLayout from "@/layouts/BaseLayout"
import { Col, Row } from "react-bootstrap"

interface AdminLayoutProps {
  children: ReactNode
  title: string
  titleButton?: ReactNode
}

export default function AdminLayout({
  children,
  title,
  titleButton
}: AdminLayoutProps) {
  return (
    <BaseLayout title={title} showLinks={true} showNav={true}>
      <Row>
        <Col>
          <h1>{title}</h1>
        </Col>
        <Col className="text-end">{titleButton}</Col>
      </Row>
      <br />

      {children}
    </BaseLayout>
  )
}
