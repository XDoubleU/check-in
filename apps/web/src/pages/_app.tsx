import "bootstrap/dist/css/bootstrap.min.css"
import "bootstrap-icons/font/bootstrap-icons.css"
import "@/styles/globals.css"
import "@/styles/scss/global.scss"

import  { type AppProps } from "next/app"
import Head from "next/head"
import { useEffect } from "react"
import NextNProgress from "nextjs-progressbar"


// eslint-disable-next-line @typescript-eslint/naming-convention
export default function App({ Component, pageProps: {session, ...pageProps} }: AppProps) {
  useEffect(() => {
    // eslint-disable-next-line @typescript-eslint/no-require-imports
    require("bootstrap/dist/js/bootstrap.bundle.min.js")
  }, [])

  return (
    <main>
      <Head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
      </Head>
      <NextNProgress color="red" />
      <Component {...pageProps} />
    </main>
  )
}
