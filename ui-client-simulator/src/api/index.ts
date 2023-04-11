export default {
  async getToken(getTokenRequest: GetTokenRequest): Promise<GetTokenResponse> {
    const basicAuth = btoa(`${getTokenRequest.clientID}:${getTokenRequest.clientSecret}`)
    const response = await fetch('https://oauth.local.tashima.space/api/token', {
      method: 'POST',
      headers: new Headers({ 'Authorization': `Basic ${basicAuth}` }),
      body: new URLSearchParams({
        grant_type: getTokenRequest.grantType,
        code: getTokenRequest.code,
        redirect_uri: getTokenRequest.redirectURI,
        client_id: getTokenRequest.clientID,
      }).toString(),
    })
    if (response.status !== 200) {
      throw new Error('Failed to get token')
    }
    return response.json()
  }
}

interface GetTokenRequest {
  grantType: string
  code: string
  redirectURI: string
  clientID: string
  clientSecret: string
}

interface GetTokenResponse {
  accessToken: string
  expiresIn: number
  refreshToken: string
  refreshTokenExpiresIn: string
  tokenType: string
}