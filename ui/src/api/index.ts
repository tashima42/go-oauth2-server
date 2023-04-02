// TODO: fix base url
const baseURL = 'http://localhost:8096'
// TODO: remove later, just for development
const authHeader = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRJRCI6IiIsImV4cCI6MTY4MDUyNTQ2OSwic2NvcGVzIjpbInVzZXJfYWNjb3VudDp1c2VyaW5mbzpyZWFkIiwiY2xpZW50OmNyZWF0ZSIsImNsaWVudDpsaXN0IiwiY2xpZW50OmluZm86cmVhZCJdLCJ1c2VyQWNjb3VudCI6eyJpZCI6IjI2MjA1MzE3LTk5MDYtNDgwMy1iZTU2LWFjNTBkZjM1OWQzYyIsInR5cGUiOiJ1c2VyIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.5hGTMAnfgxp6jH0G5prdN7QoQQtJXqGXHqQcqSQ5gdQ"
export default function buildOauth2API(): Oauth2API {
  return {
    login,
    getClientInfo,
  }
  async function login(username: string, password: string): Promise<LoginResponse> {
    return fetch(`${baseURL}/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    })
      .then((response) => {
        if (response.status !== 200) {
          response.json()
          throw new Error('Invalid credentials')
        }
      })
      .then(() => ({ success: true }))
      .catch((error) => {
        console.error(error)
        return { success: false, error: error.message, errorCode: 'INVALID_CREDENTIALS' }
      })
  }
  async function getClientInfo(clientID: string): Promise<ClientInfoResponse> {
    return fetch(`${baseURL}/clients/${clientID}`, {
      method: 'GET',
      headers: { 'Authorization': `Bearer ${authHeader}` }
    })
      .then((response) => {
        if (response.status !== 200) {
          throw new Error('Client not found')
        }
        return response.json()
      })
      .catch((error) => {
        console.error(error)
        return { clientID, name: 'Client not found' }
      })
  }
}
type LoginResponse = {
  success: boolean
  error?: string
  errorCode?: string
}
type ClientInfoResponse = {
  clientID: string
  name: string
}
interface Oauth2API {
  login: (username: string, password: string) => Promise<LoginResponse>
  getClientInfo(clientID: string): Promise<ClientInfoResponse>
}