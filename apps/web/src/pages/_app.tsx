import "bootstrap/dist/css/bootstrap.min.css"
import "bootstrap-icons/font/bootstrap-icons.css"
import "@/styles/globals.css"
import "@/styles/scss/global.scss"

import { type AppProps } from "next/app"
import Head from "next/head"
import { useEffect } from "react"
import { AuthProvider } from "@/contexts/authContext"
import { RedirectsProvider } from "@/contexts/redirectsContext"
import { LoadingProvider } from "@/contexts/loadingContext"

export default function App({
  // eslint-disable-next-line @typescript-eslint/naming-convention
  Component,
  pageProps: { session, ...pageProps }
}: AppProps) {
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
      <LoadingProvider>
        <AuthProvider>
          <RedirectsProvider>
            <Component {...pageProps} />
          </RedirectsProvider>
        </AuthProvider>
      </LoadingProvider>
    </main>
  )
}
