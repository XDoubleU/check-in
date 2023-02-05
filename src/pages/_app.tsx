import "bootstrap/dist/css/bootstrap.min.css"
import "bootstrap-icons/font/bootstrap-icons.css"
import "@/styles/globals.css"

import type { AppProps } from "next/app"
import Head from "next/head"
import { useEffect } from "react"
import localFont from "@next/font/local"
import { SessionProvider } from "next-auth/react"

const brandon = localFont({
  src: [
    {
      path: "../../public/fonts/brandon_bld.otf",
      weight: "bold",
      style: "bold"
    },
    {
      path: "../../public/fonts/brandon_reg.otf",
      weight: "400",
      style: "normal"
    }
  ]
})


export default function App({ Component, pageProps: {session, ...pageProps} }: AppProps) {
  useEffect(() => {
    require("bootstrap/dist/js/bootstrap.bundle.min.js")
  }, [])

  return (
    <main className={brandon.className}>
      <Head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
      </Head>
      <SessionProvider session={session}>
        <Component {...pageProps} />
      </SessionProvider>
    </main>
  )
}
