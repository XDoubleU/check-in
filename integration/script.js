const webSocket = new WebSocket("ws://localhost:8000")

webSocket.onopen = async () => {
  const response = await fetch("http://localhost:8000/locations/ws/")
  const data = await response.json()

  data.forEach((location) => {
    fill(location)
  })
}

webSocket.onmessage = (event) => {
  const data = JSON.parse(event.data)
  fill(data)
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

    capacity.innerHTML = `<b>${location.available}</b> of the <b>${location.capacity}</b> spots remaining`

    let output = "Wasn't full yesterday"
    if (location.yesterdayFullAt) {
      const time = new Date(parseInt(location.yesterdayFullAt)).toLocaleTimeString([], {
        timeStyle: "short",
        hourCycle: "h23"
      })

      output = `Yesterday full at: ${time}`
    }

    yesterdayFullAt.innerHTML = output

    element.innerHTML = ""
    element.appendChild(capacity)
    element.appendChild(yesterdayFullAt)
  }
}
