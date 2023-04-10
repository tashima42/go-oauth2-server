<script setup lang="ts">
import { RouterLink } from 'vue-router';
import LoginForm from '@/components/LoginForm.vue';
import IconGolang from '@/components/icons/IconGolang.vue';
</script>

<template>
  <main>
    <a href="https://github.com/tashima42/go-oauth2-server" target="_blank" rel="noopener noreferrer">
      <IconGolang class="go-icon"/>
    </a>
    <LoginForm class="login-form" />
    <p>Don't have an account? <RouterLink to="signup">Sign up!</RouterLink></p>
  </main>
</template>

<style scoped>
  main {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
  }
  .go-icon {
    height: 40px;
    width: 100px;
  }
  .login-form {
    margin-top: 1rem;
    margin-bottom: 1rem;
  }
</style>
