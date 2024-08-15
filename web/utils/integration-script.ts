import fs from "fs"
import UglifyJs from "uglify-js"

export function generateIntegrationScripts() {
    const dir = "./public/scripts"

  if (!fs.existsSync(dir)){
    fs.mkdirSync(dir, { recursive: true });
  }

  ["eng", "nl"].forEach(language => {
    fs.writeFileSync(`${dir}/${language}.js`, getScript(process.env.NEXT_PUBLIC_API_URL ?? "", language))
  });  
}

function getNow(language: string) {
    switch (language) {
        case "eng":
            return "`<b>${location.available}</b> of the <b>${location.capacity}</b> spots remaining`"
        case "nl":
            return "`Nog <b>${location.available}</b> van de <b>${location.capacity}</b> plaatsen vrij`"
        default:
            return ""
    }
}

function getYesterday(language: string) {
    switch (language) {
        case "eng":
            return "`Yesterday <b>${location.availableYesterday}</b> of the <b>${location.capacityYesterday}</b> spots remained`"
        case "nl":
            return "`Gisteren nog <b>${location.availableYesterday}</b> van de <b>${location.capacityYesterday}</b> plaatsen over`"
        default:
            return ""
    }
}

function getYesterdayFull(language: string) {
    switch (language) {
        case "eng":
            return "`Yesterday full at ${time}`"
        case "nl":
            return "`Gisteren volzet om ${time}`"
        default:
            return ""
    }
}

// eslint-disable-next-line max-lines-per-function
function getScript(apiUrl: string, language: string){
    const wsUrl = apiUrl.replace("http", "ws")

    const now = getNow(language)
    const yesterday = getYesterday(language)
    const yesterdayFull = getYesterdayFull(language)

    const script = `
        const webSocket = new WebSocket("${wsUrl}")

        webSocket.onopen = async () => {
            webSocket.send(
                JSON.stringify({
                subject: "all-locations"
                })
            )
        }

        webSocket.onmessage = (event) => {
            let data = JSON.parse(event.data)
            
            if(!Array.isArray(data)) {
                data = [data]
            }

            data.forEach((location) => {
                fill(location)
            })
        }

        window.onbeforeunload = () => {
            webSocket.close()
        }

        function fill(location){
            let element = document.getElementById(location.normalizedName);

            if(element){
                const capacity = document.createElement("p")
                const yesterdayFullAt = document.createElement("p")
                yesterdayFullAt.style.marginTop = "-15px"

                capacity.innerHTML = ${now}

                let output = ${yesterday}
                if (location.yesterdayFullAt) {
                    const time = new Date(location.yesterdayFullAt).toUTCString([], {
                        timeStyle: "short",
                        hourCycle: "h23"
                    })

                    output = ${yesterdayFull}
                }

                yesterdayFullAt.innerHTML = output

                element.innerHTML = ""
                element.appendChild(capacity)
                element.appendChild(yesterdayFullAt)
            }
        }
    `

    return UglifyJs.minify(script).code
}
