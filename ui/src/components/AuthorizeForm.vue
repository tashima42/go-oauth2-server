<script setup lang="ts">
// import router from '@/router';
import router from '@/router';
import { ref, onMounted } from 'vue';
import api from "@/api";
const clientName = ref('TODO REPLACE');
const query = router.currentRoute.value.query;
const scopeRaw = query.scope ? query.scope as string : '';
const scopes = scopeRaw.split(' ');

onMounted(async () => {
  try {
    await api.userInfo()
  } catch (error) {
    console.error(error)
    router.push('/')
  }
  const { name } = await api.clientInfo(query.client_id as string)  
  clientName.value = name
})

function authorize() {
  window.location.href = window.location.href.replace("/authorize", "/api/authorize")
}
</script>

<template>
  <form autocomplete="on" @submit="e => e.preventDefault()" class="authorize-form">
    <h1><span>{{ clientName }}</span> wants access to:</h1>
    <ul class="scopes-list">
      <li v-for="scope, i in scopes" :key="i">{{ scope }}</li>
    </ul>
    <button v-on:click="authorize">Authorize</button>
  </form>
</template>

<style scoped>
.authorize-form {
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

.scopes-list {
  text-align: left;
}

.authorize-form h1 span {
  font-weight: bold;
}

.form-group {
  display: flex;
  flex-direction: column;
  margin: 0.5rem;
  width: 100%;
}

.authorize-form button {
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
</style>