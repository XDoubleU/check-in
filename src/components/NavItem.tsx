import { useRouter } from "next/router"

type NavItemProps = {
  name: string,
  url: string
}


export default function NavItem({name, url}: NavItemProps){
  const router = useRouter()
  const selected = router.pathname.includes(name)

  return (
    <li className="nav-item">
      <a className={`nav-link ${selected ? "active" : ""}`} href={`${url}`}>{name}</a>
    </li>
  )
}