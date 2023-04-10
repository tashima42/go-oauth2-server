<script setup lang="ts">
import IconGolang from '@/components/icons/IconGolang.vue';
import AuthorizeForm from '@/components/AuthorizeForm.vue';
</script>

<template>
  <main>
    <a href="https://github.com/tashima42/go-oauth2-server" target="_blank" rel="noopener noreferrer">
      <IconGolang class="go-icon"/>
    </a>
    <AuthorizeForm class="authorize" />
  </main>
</template>

<style scoped>
  main {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
  }
  .go-icon {
    height: 40px;
    width: 100px;
  }
  .authorize-form {
    margin-top: 1rem;
    margin-bottom: 1rem;
  }
</style>
