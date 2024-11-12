import styles from "./error.module.css"
import Head from "next/head"
import { type MouseEventHandler, useRef } from "react"

export default function CustomError() {
  const leftEye = useRef<HTMLDivElement>(null)
  const rightEye = useRef<HTMLDivElement>(null)

  const moveEyes: MouseEventHandler<HTMLDivElement> = (event) => {
    if (!leftEye.current || !rightEye.current) return

    const lx = leftEye.current.offsetLeft + leftEye.current.offsetWidth / 2
    const ly = leftEye.current.offsetTop + leftEye.current.offsetHeight / 2
    const lrad = Math.atan2(event.pageX - lx, event.pageY - ly)
    const lrot = lrad * (180 / Math.PI) * -1 + 180

    const rx = rightEye.current.offsetLeft + rightEye.current.offsetWidth / 2
    const ry = rightEye.current.offsetTop + rightEye.current.offsetHeight / 2
    const rrad = Math.atan2(event.pageX - rx, event.pageY - ry)
    const rrot = rrad * (180 / Math.PI) * -1 + 180

    leftEye.current.style.transform = `rotate(${lrot.toString()}deg)`
    rightEye.current.style.transform = `rotate(${rrot.toString()}deg)`
  }

  return (
    <>
      <Head>
        <title>Something went wrong</title>
      </Head>

      <div className={styles.body} onMouseMove={moveEyes}>
        <div>
          <span className={styles.errorNum}>5</span>
          <div ref={leftEye} className={styles.eye}></div>
          <div ref={rightEye} className={styles.eye}></div>
          <p className={styles.subText}>
            Something went wrong. We&apos;re <i>looking</i> into what happened.
          </p>
        </div>
      </div>
    </>
  )
}
