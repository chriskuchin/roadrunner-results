import { auth } from '../../firebase'

const storageKey = "api-token"

function saveAPIToken(token) {
  localStorage.setItem(storageKey, token)
}

async function setAuthHeader(fetchObject) {
  var token
  if (auth.currentUser) {
    token = await auth.currentUser.getIdToken(true)
  } else {
    token = localStorage.getItem(storageKey)
  }

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