import { getMyUser, signOut } from "api-wrapper"
import Router, { useRouter } from "next/router"
import { MouseEventHandler, useEffect, useState } from "react"
import { Container, Nav, Navbar } from "react-bootstrap"
import { Role, User } from "types"

type NavItemProps = {
  children: string,
  href?: string,
  onClick?: MouseEventHandler<HTMLAnchorElement>,
  active?: boolean
}

function NavItem({children, href, onClick, active}: NavItemProps) {
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
    <Nav.Link className={`nav-link ${selected || active ? "active" : ""}`} href={href}>
      {children}
    </Nav.Link>
  )
}

export default function Navigation(){
  const [user, setUser] = useState<User | undefined>(undefined)

  useEffect(() => {
    getMyUser()
      .then(data => {
        if (data === null) {
          Router.push("/signin")
        } else {
          setUser(data)
        }
      })
  }, [])


  const signOutHandler = () => {
    signOut()
    Router.push("/signin")
  }

  return (
    <Navbar expand="lg" bg="primary" variant="dark">
      <Container>
        <Navbar.Brand className="bold">CHECK-IN</Navbar.Brand>
        <Navbar.Toggle aria-controls="navbar-nav" />
        <Navbar.Collapse id="navbar-nav">
          <Nav className="me-auto mb-2 mb-lg-0">
            {
              user && !user.roles.includes(Role.Admin) ? (
                <NavItem active={true} href={`/settings/locations/${user.locationId}`} >My location</NavItem>
              ) : (
                <>
                  <NavItem href="/settings/locations">Locations</NavItem>
                  <NavItem href="/settings/schools">Schools</NavItem>
                </>
              )
            }
          </Nav>
          <Nav className="ms-auto mb-2 mb-lg-0">
            <NavItem onClick={signOutHandler}>Sign out</NavItem>
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  )
}