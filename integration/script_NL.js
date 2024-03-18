const webSocket = new WebSocket("ws://localhost:8000")

webSocket.onopen = async () => {
  webSocket.send(
    JSON.stringify({
      subject: "all-locations"
    })
  )
}

webSocket.onmessage = (event) => {
  const data = JSON.parse(event.data)
  
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

    capacity.innerHTML = `Nog <b>${location.available}</b> van de <b>${location.capacity}</b> plaatsen vrij`

    let output = `Gisteren nog <b>${location.availableYesterday}</b> van de <b>${location.capacityYesterday}</b> plaatsen over`
    if (location.yesterdayFullAt) {
      const time = new Date(location.yesterdayFullAt).toUTCString([], {
        timeStyle: "short",
        hourCycle: "h23"
      })

      output = `Gisteren volzet om ${time}`
    }

    yesterdayFullAt.innerHTML = output

    element.innerHTML = ""
    element.appendChild(capacity)
    element.appendChild(yesterdayFullAt)
  }
}
