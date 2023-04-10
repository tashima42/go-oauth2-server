<script setup lang="ts">
import { ref } from 'vue';
import api from '@/api';
import router from '@/router';

const username = ref('');
const password = ref('');

async function login() {
  const loginRequest = {
    username: username.value,
    password: password.value
  }
  await api.login(loginRequest);
  router.push('/authorize')
}
</script>

<template>
  <form autocomplete="on" @submit="e => e.preventDefault()" class="login-form">
    <h1>Welcome back!</h1>
    <div class="form-group">
      <label for="username">Username:</label>
      <input type="text" id="username" v-model="username" />
    </div>
    <div class="form-group">
      <label for="password">Password:</label>
      <input type="password" id="password" v-model="password" />
    </div>
    <button v-on:click="login">Login</button>
  </form>
</template>

<style scoped>
.login-form {
  display: flex;
  flex-direction: column;
  align-items: center;
  background-color: var(--color-background-soft);
  min-width: 400px;
  min-height: 200px;
  border-radius: 10px;
  color: var(--color-text);
  padding: 1rem;
}

.login-form h1 {
  font-weight: bold;
}

.form-group {
  display: flex;
  flex-direction: column;
  margin: 0.5rem;
  width: 100%;
}

.login-form button {
  margin-top: 1rem;
  margin-bottom: 1rem;
  width: 100%;
  padding: 0.5rem;
  border-radius: 5px;
  border: 1px solid var(--color-border);
  background-color: var(--color-primary);
  color: var(--vt-c-white);
  font-weight: bold;
}

.login-form button:hover {
  border: 1px solid var(--vt-c-white);
}

.form-group input {
  margin-top: 0.5rem;
  padding: 0.5rem;
  border-radius: 5px;
  border: 1px solid var(--color-border);
  background-color: var(--color-background);
  color: var(--color-text);
}
</style>