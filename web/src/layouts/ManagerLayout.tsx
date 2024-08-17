import { useCallback, useEffect, useState, type ReactNode } from "react"
import BaseLayout from "./BaseLayout"
import { Col, Container, Row } from "react-bootstrap"
import { type State } from "api-wrapper/types/apiTypes"
import { getState } from "api-wrapper"
import StateAlert from "components/StateAlert"

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
  const [apiState, setApiState] = useState<State>();

  const fetchState = useCallback(async () => {
    setApiState((await getState()).data)
  }, [])

  useEffect(() => {
    void fetchState()
  }, [fetchState])

  return <>
    <BaseLayout title={title} showLinks={true} showNav={true}>
      <Row>
        <Col>
          <StateAlert state={apiState} />
        </Col>
      </Row>
      <Row>
        <Col>
          <h1>{title}</h1>
        </Col>
        <Col className="text-end">{titleButton}</Col>
      </Row>
      <br />

      <Container style={{ minHeight: "65vh" }}>{children}</Container>
    </BaseLayout>
  </>
}
