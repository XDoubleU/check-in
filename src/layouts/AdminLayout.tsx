import { ReactNode } from "react"
import BaseLayout from "./BaseLayout"
import Navigation from "@/components/Navigation"

type AdminLayoutProps = {
  children: ReactNode,
  title: string,
  isSuperUser: boolean
}

export default function AdminLayout({children, title, isSuperUser}: AdminLayoutProps){
  return (
    <BaseLayout title={title}>
      <Navigation isSuperUser={isSuperUser}/>

      <div className="container content">
        {children}
      </div>

      <br/>
      <br/>

      <footer className="text-center">
        <br/>

        <p>Made with <i className="bi bi-heart-fill" style={{"color": "red"}}></i> by <a href="https://xdoubleu.com">XDoubleU</a> for <a href="https://bruggestudentenstad.be/">Brugge Studentenstad</a></p>
      </footer>
    </BaseLayout>
  )   
}