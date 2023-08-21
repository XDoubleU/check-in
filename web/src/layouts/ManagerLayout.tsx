import { type ReactNode } from "react"
import BaseLayout from "./BaseLayout"
import { Col, Container, Row } from "react-bootstrap"

interface ManagerLayoutProps {
  children: ReactNode
  title: string
  titleButton?: ReactNode
}

export default function ManagerLayout({
  children,
  title,
  titleButton
}: ManagerLayoutProps) {
  return (
    <BaseLayout title={title} showLinks={true} showNav={true}>
      <Row>
        <Col>
          <h1>{title}</h1>
        </Col>
        <Col className="text-end">{titleButton}</Col>
      </Row>
      <br />

      <Container style={{ minHeight: "65vh" }}>{children}</Container>
    </BaseLayout>
  )
}
