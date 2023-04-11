<script setup lang="ts">
import router from '@/router';
import { ref, onMounted } from 'vue';
import api from "@/api";
const query = router.currentRoute.value.query;
const code = query.code as string;
const clientID = ref("")
const clientSecret = ref("")
const redirectURI = ref("")

async function getToken() {
  try {
    const response = await api.getToken({
        grantType: "authorization_code",
        code,
        redirectURI: redirectURI.value,
        clientID: clientID.value,
        clientSecret: clientSecret.value,
    })
    // TODO: show response in some place
    console.log(response)
  } catch (error) {
    console.error(error)
  }
}
</script>

<template>
  <form autocomplete="on" @submit="e => e.preventDefault()" class="client-info-form">
    <div class="form-group">
      <label for="clientID">Client ID:</label>
      <input type="text" id="clientID" v-model="clientID" />
    </div>
    <div class="form-group">
      <label for="clientSecret">Client Secret:</label>
      <input type="text" id="clientSecret" v-model="clientSecret" />
    </div>
    <div class="form-group">
      <label for="redirectURI">Redirect URI:</label>
      <input type="text" id="redirectURI" v-model="redirectURI" />
    </div>
    <button v-on:click="getToken">Get Token</button>
  </form>
</template>

<style scoped>
.client-info-form {
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

.client-info-form h1 span {
  font-weight: bold;
}

.form-group {
  display: flex;
  flex-direction: column;
  margin: 0.5rem;
  width: 100%;
}

.client-info-form button {
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

.authorize-form button:hover {
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