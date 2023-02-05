import { ReactNode } from "react"
import Head from "next/head"

type BaseLayoutProps = {
  children: ReactNode,
  title?: string
}

export default function BaseLayout({children, title}: BaseLayoutProps){
  const fullTitle = title ? `${title} - Check-In` : "Check-In"
  
  return (
  <>
    <Head>
      <title>{fullTitle}</title>
    </Head>

    {children}
  </>
  )   
}