const storageKey = "api-token"

function saveAPIToken(token) {
  console.log(token)
  localStorage.setItem(storageKey, token)
}

function setAuthHeader(fetchObject) {
  var token = localStorage.getItem(storageKey)

  fetchObject.headers['X-Api-Token'] = token

  console.log(fetchObject)

  return fetchObject
}

function clearAPIToken() {
  localStorage.removeItem(storageKey)
}

export {
  saveAPIToken,
  setAuthHeader,
  clearAPIToken
}