'use strict'

const heroku = "https://afternoon-wave-39227.herokuapp.com/purplecat/api/"

const availableUrl = (urlString) => {
    const url = new URL(urlString)
    return url.hostname != "" && url.pathname.endsWith(".pom")
}

const checkUrl = (e) => {
    const text = document.getElementById("pomurl")

    const runButton = document.getElementById('runButton')
    runButton.disabled = !availableUrl(text.value)
}

const reset = (e) => {
    const text = document.getElementById("pomurl")
    text.value = ""

    const area = document.getElementById("resultArea")
    area.innerText = ""

    const runButton = document.getElementById('runButton')
    runButton.disbled = true

    const message = document.getElementById('message')
    message.innerText = ""
}

const showMessage = (message) => {
    const messageArea = document.getElementById('message')
    messageArea.innerText = message
    messageArea.classList.remove("warning")
}

const showError = (message) => {
    const messageArea = document.getElementById('message')
    messageArea.innerText = message
    messageArea.classList.add("warning")
}

const showResult = (jsonString) => {
    const json = JSON.parse(jsonString)
    const str = JSON.stringify(json, null, "  ")
    const area = document.getElementById('resultArea')
    console.log(str)
    area.innerText = str
}

const executePurplecat = (e) => {
    const url = document.getElementById("pomurl").value
    const request = new XMLHttpRequest()
    request.onreadystatechange = () => {
        console.log(`readyState: ${request.readyState}`)
        if (request.readyState == 4) {
            console.log(`done: http status: ${request.status}`)
            if (request.status == 200) {
                showMessage(`GET license data from ${url}`)
                showResult(request.responseText)
            } else {
                showError(request.responseText)
            }
        } else {
            showMessage("running purplecat...")
        }
    }
    console.log(`${heroku}licenses?target=${url}`)
    request.open("GET", `${heroku}licenses?target=${url}`, true)
    request.send(null)
}

const init = () => {
    const text = document.getElementById("pomurl")
    text.addEventListener('change', checkUrl)

    const runButton = document.getElementById('runButton')
    runButton.disbled = true
    runButton.addEventListener('click', executePurplecat)

    const resetButton = document.getElementById('resetButton')
    resetButton.addEventListener('click', reset)

    console.log('init done')
}

window.onload = () => {
    init()
}
