/* eslint-disable sonarjs/no-duplicate-string */
import { type Role } from "api-wrapper/types/apiTypes"
import Loader from "components/Loader"
import { AuthRedirecter, useAuth } from "contexts/authContext"
import ManagerLayout from "layouts/ManagerLayout"
import Link from "next/link"
import { useRouter } from "next/router"
import { useEffect, useState } from "react"

export const ImageWidth = "60%"

const supportedLanguages = ["en", "nl"]

// eslint-disable-next-line max-lines-per-function
export function ManualNavigation() {
  const router = useRouter()
  const { user } = useAuth()
  const [language, setLanguage] = useState("en")
  const [currentEndpoint, setCurrentEndpoint] = useState("")

  useEffect(() => {
    const regex = "/([a-z]{2})/"
    const match = (router.pathname.match(regex) as RegExpMatchArray)[1]
    setLanguage(match)

    const endpoint = router.pathname.split(`/manual/${match}/`)[1]
    setCurrentEndpoint(endpoint)
  }, [router])

  const linkTitles = new Map<string, Map<string, string>>([
    [
      "manager",
      new Map<string, string>([
        ["en", "Manual Manager"],
        ["nl", "Handleiding Beheerder"]
      ])
    ],
    [
      "location",
      new Map<string, string>([
        ["en", "Manual Location"],
        ["nl", "Handleiding Locatie"]
      ])
    ]
  ])

  return (
    <>
      <ul>
        {language === "en" ? (
          <li>
            <Link href={`/manual/nl/${currentEndpoint}`}>
              Verander naar Nederlands
            </Link>
          </li>
        ) : (
          <li>
            <Link href={`/manual/en/${currentEndpoint}`}>
              Switch to English
            </Link>
          </li>
        )}

        {user?.role === "admin" || user?.role === "manager" ? (
          <>
            <li>
              <Link href={`/manual/${language}/manager`}>
                {linkTitles.get("manager")?.get(language)}
              </Link>
            </li>
            <li>
              <Link href={`/manual/${language}/location`}>
                {linkTitles.get("location")?.get(language)}
              </Link>
            </li>
          </>
        ) : (
          <></>
        )}
      </ul>
    </>
  )
}

export const ManagerRedirects = new Map<Role, string>([
  ["default", "/manual/location"]
])

export default function ManualHome() {
  const [redirects, setRedirects] = useState(
    new Map<Role, string>([
      ["admin", "/manual/en/manager"],
      ["manager", "/manual/en/manager"],
      ["default", "/manual/en/location"]
    ])
  )

  useEffect(() => {
    let detectedLanguage = window.navigator.language
    if (!supportedLanguages.includes(detectedLanguage)) {
      detectedLanguage = "en"
    }

    setRedirects(
      new Map<Role, string>([
        ["admin", `/manual/${detectedLanguage}/manager`],
        ["manager", `/manual/${detectedLanguage}/manager`],
        ["default", `/manual/${detectedLanguage}/location`]
      ])
    )
  }, [])

  return (
    <AuthRedirecter redirects={redirects}>
      <ManagerLayout title="">
        <Loader message="Loading manual." />
      </ManagerLayout>
    </AuthRedirecter>
  )
}
