export default {
  async login(LoginRequest: LoginRequest): Promise<LoginResponse> {
    return fetch('/api/login', {
      method: 'POST',
      body: JSON.stringify(LoginRequest),
    }).then(response => response.json())
  },
  async userInfo(): Promise<UserInfoResponse> {
    const response = await fetch('/api/userinfo', { method: 'GET' })
    if (response.status !== 200) {
      throw new Error('Failed to get user info')
    }
    return response.json()
  },
  async clientInfo(clientID: string): Promise<ClientInfoResponse> {
    const response = await fetch(`/api/clients/${clientID}`, { method: 'GET' })
    if (response.status !== 200) {
      throw new Error('Failed to get client info')
    }
    return response.json()
  },
}

interface LoginRequest {
  username: string;
  password: string;
}

interface LoginResponse {
  token: string;
}

interface UserInfoResponse {
  username: string
  type: string
  scopes: Array<string>
}

interface ClientInfoResponse {
  clientID: string
  name: string
}