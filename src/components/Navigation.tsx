import { signOut } from "next-auth/react"
import { useRouter } from "next/router"
import { MouseEventHandler } from "react"
import { Container, Nav, Navbar } from "react-bootstrap"

type NavigationProps = {
  isAdmin: boolean
}

type NavItemProps = {
  children: string,
  href?: string,
  onClick?: MouseEventHandler<HTMLAnchorElement>
}

function NavItem({children, href, onClick}: NavItemProps) {
  const router = useRouter()
  
  if (onClick !== undefined) {
    return (
      <Nav.Link onClick={onClick}>
        {children}
      </Nav.Link>
    )
  }

  if (href === undefined) {
    throw new Error("href or onClick has to be defined.")
  }

  const selected = router.pathname.includes(href.toLowerCase())
  return (
    <Nav.Link className={`nav-link ${selected ? "active" : ""}`} href={href}>
      {children}
    </Nav.Link>
  )
}

export default function Navigation({isAdmin}: NavigationProps){
  return (
    <Navbar expand="lg" bg="primary" variant="dark">
      <Container>
        <Navbar.Brand className="bold">CHECK-IN</Navbar.Brand>
        <Navbar.Toggle aria-controls="navbar-nav" />
        <Navbar.Collapse id="navbar-nav">
          <Nav className="me-auto mb-2 mb-lg-0">
            {
              !isAdmin ? (
                <NavItem href="#" >My location</NavItem>
              ) : (
                <>
                  <NavItem href="/settings/locations">Locations</NavItem>
                  <NavItem href="/settings/schools">Schools</NavItem>
                </>
              )
            }
          </Nav>
          <Nav className="ms-auto mb-2 mb-lg-0">
            <NavItem onClick={() => signOut()}>Log out</NavItem>
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  )
}