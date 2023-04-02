import { useState } from 'react'
import Head from 'next/head'
import { useRouter } from 'next/router'
import styles from '@ui/styles/Login.module.css'
import buildOauth2API from '@ui/api'

export default function Login() {
  const router = useRouter()
  const oauth2API = buildOauth2API()
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  return (
    <>
      <Head>
        <title>Login Go Oauth2 Server</title>
      </Head>
      <main className={styles.main}>
        <div className={styles.loginForm}>
          <input className={styles.credentialsInput} onChange={(e) => setUsername(e.target.value)} type="text" placeholder='Username' />
          <input className={styles.credentialsInput} onChange={(e) => setPassword(e.target.value)} type="password" placeholder='Password' />
          <button className={styles.loginButton} onClick={login}>Login</button>
        </div>
      </main>
    </>
  )

  async function login() {
    const response = await oauth2API.login(username, password)
    if (response.success === true) {
      router.push({ pathname: "/authorize", query: router.query })
    }
    console.log(response)
  }
}