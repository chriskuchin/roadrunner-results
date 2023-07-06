const storageKey = "api-token"

function saveAPIToken(token) {
  localStorage.setItem(storageKey, token)
}

function setAuthHeader(fetchObject) {
  var token = localStorage.getItem(storageKey)

  if (!fetchObject.headers)
    fetchObject.headers = {}

  fetchObject.headers['X-Api-Token'] = token

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