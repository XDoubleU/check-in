import NavItem from "./NavItem"

type NavigationProps = {
  isSuperUser: boolean
}

function DynamicNavItems({isSuperUser}: NavigationProps){
  if (!isSuperUser){
    return <NavItem name="My location" url="#" />
  }

  return (
    <>
      <NavItem name="Locations" url="#" />
      <NavItem name="Schools" url="#" />
    </>
  )
}

export default function Navigation({isSuperUser}: NavigationProps){
  return (
    <>
      <nav className="navbar navbar-expand-lg bg-primary-custom navbar-dark">
        <div className="container">
          <a className="navbar-brand bold" href="#">CHECK-IN</a>
          <button className="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
            <span className="navbar-toggler-icon"></span>
          </button>
          <div className="collapse navbar-collapse" id="navbarSupportedContent">
            <ul className="navbar-nav me-auto mb-2 mb-lg-0">
              <DynamicNavItems isSuperUser={isSuperUser} />
            </ul>
            <ul className="navbar-nav ms-auto mb-2 mb-lg-0">
              <NavItem name="Log out" url="#" />
            </ul>
          </div>
        </div>
      </nav>
    </>
  )
}