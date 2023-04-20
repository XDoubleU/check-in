import { type ReactNode } from "react"
import Head from "next/head"
import { Container } from "react-bootstrap"
import Navigation from "components/Navigation"

interface BaseLayoutProps {
  children: ReactNode
  title?: string
  showLinks?: boolean
  showNav?: boolean
}

export default function BaseLayout({
  children,
  title,
  showLinks,
  showNav
}: BaseLayoutProps) {
  const fullTitle = title ? `${title} - Check-In` : "Check-In"

  return (
    <>
      <Head>
        <title>{fullTitle}</title>
      </Head>

      {showNav ? <Navigation /> : <></>}

      <Container className="content">{children}</Container>

      <br />

      <footer className="text-center">
        <br />

        {showLinks ? (
          <p>
            Made with{" "}
            <i className="bi bi-heart-fill" style={{ color: "red" }}></i> by{" "}
            <a href="https://xdoubleu.com">XDoubleU</a> for{" "}
            <a href="https://bruggestudentenstad.be/">Brugge Studentenstad</a>
          </p>
        ) : (
          <p>
            Made with{" "}
            <i className="bi bi-heart-fill" style={{ color: "red" }}></i> by
            XDoubleU for Brugge Studentenstad
          </p>
        )}
      </footer>
    </>
  )
}
