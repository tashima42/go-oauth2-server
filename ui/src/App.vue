<script setup lang="ts">
import { RouterView } from 'vue-router'
</script>

<template>
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Go Oauth2 Server</title>
    </head>
  <RouterView class="view" />
</template>

<style scoped>
.view {
  width: 100vw;
  height: 100vh;
}
</style>
