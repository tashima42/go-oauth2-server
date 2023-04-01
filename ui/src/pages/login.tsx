import Head from 'next/head'
import styles from '@ui/styles/Login.module.css'

export default function Login() {
  return (
    <>
      <Head>
        <title>Login Go Oauth2 Server</title>
      </Head>
      <main className={styles.main}>
        <div className={styles.loginForm}>
          <input className={styles.credentialsInput} type="text" placeholder='Username'/>
          <input className={styles.credentialsInput} type="password" placeholder='Password'/>
          <button className={styles.loginButton}>Login</button>
        </div>
      </main>
    </>
  )
}