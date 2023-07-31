import { useAuth } from "contexts/authContext"
import { useRouter } from "next/router"
import { type MouseEventHandler } from "react"
import { Container, Nav, Navbar } from "react-bootstrap"

interface NavItemProps {
  children: string
  href?: string
  onClick?: MouseEventHandler<HTMLAnchorElement>
  active?: boolean
}

function NavItem({ children, href, onClick, active }: NavItemProps) {
  const router = useRouter()

  if (onClick !== undefined) {
    return <Nav.Link onClick={onClick}>{children}</Nav.Link>
  }

  if (href === undefined) {
    throw new Error("href or onClick has to be defined.")
  }

  const selected = router.pathname.includes(href.toLowerCase())
  return (
    <Nav.Link
      className={`nav-link ${selected || active ? "active" : ""}`}
      href={href}
    >
      {children}
    </Nav.Link>
  )
}

export default function Navigation() {
  const { user } = useAuth()

  return (
    <Navbar expand="lg" bg="primary" variant="dark">
      <Container>
        <Navbar.Brand className="bold">CHECK-IN</Navbar.Brand>
        <Navbar.Toggle aria-controls="navbar-nav" />
        <Navbar.Collapse id="navbar-nav">
          <Nav className="me-auto mb-2 mb-lg-0">
            {user?.role === "default" && user.location ? (
              <NavItem
                active={true}
                href={`/settings/locations/${user.location.id}`}
              >
                My location
              </NavItem>
            ) : (
              <></>
            )}
            {user?.role === "manager" || user?.role === "admin" ? (
              <>
                <NavItem href="/settings/locations">Locations</NavItem>
                <NavItem href="/settings/schools">Schools</NavItem>
              </>
            ) : (
              <></>
            )}
            {user?.role === "admin" ? (
              <>
                <NavItem href="/settings/users">Users</NavItem>
              </>
            ) : (
              <></>
            )}
          </Nav>
          <Nav className="ms-auto mb-2 mb-lg-0">
            <NavItem href="/signout">Sign out</NavItem>
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  )
}
