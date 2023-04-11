<script setup lang="ts">
import IconGolang from '@/components/icons/IconGolang.vue';
</script>

<template>
  <main>
    <a href="https://github.com/tashima42/go-oauth2-server" target="_blank" rel="noopener noreferrer">
      <IconGolang class="go-icon"/>
    </a>
    <div>
      <h1>Congrats! You're logged in</h1>
    </div>
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
</style>
