import { useEffect, useState } from 'react'
import Head from 'next/head'
import { useRouter } from 'next/router'
import styles from '@ui/styles/Authorize.module.css'
import buildOauth2API from '@ui/api'

export default function Authorize() {
  const oauth2API = buildOauth2API()
  const { query, isReady } = useRouter()
  const [clientName, setClientName] = useState<string>('')
  const [scopes, setScopes] = useState<ScopeInfo[]>([])
  useEffect(() => {
    if (!isReady) return

    getClientInfo(query.client_id as string).then((clientInfo) => {
      setClientName(clientInfo.name)
    })
    setScopes(parseScopes(query.scope as string))
  }, [isReady])
  return (
    <>
      <Head>
        <title>Login Go Oauth2 Server</title>
      </Head>
      <main className={styles.main}>
        <div className={styles.acceptForm}>
          <h1>{clientName.toUpperCase()}</h1>
          <p>Want's the following permissions</p>
          <div className={styles.scopes}>
            {scopes.map((scope) => {
              return (
                <div key={scope.scope} className={styles.scope}>
                  <p>{scope.scope}: {scope.description}</p>
                </div>
              )
            })
            }
          </div>
          <button className={styles.acceptButton} onClick={redirectToAuthorizeServer} >GIVE PERMISSIONS</button>
        </div>
      </main>
    </>
  )

  async function getClientInfo(clientID: string): Promise<{ clientID: string, name: string }> {
    return await oauth2API.getClientInfo(clientID)
  }

  function parseScopes(scopes: string): ScopeInfo[] {
    return scopes.split('+').map((scope) => {
      return {
        scope,
        description: 'TODO'
      }
    })
  }

  function redirectToAuthorizeServer() {
    // TODO: rework this, this method allows an attack where the user is not presented
    // to the authorize page and it's just redirect direclty to the server
    window.location.href = 'http://localhost:8096/authorize?'+ encodeURIComponent(queryAsString())
  }
  function queryAsString(): string {
    return Object.keys(query).map((key) => {
      return `${key}=${query[key]}`
    }).join('&')
  }
}

type ScopeInfo = {
  scope: string
  description: string
}