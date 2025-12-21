<script setup lang="ts">
import "~/assets/styles/main.css";

</script>

<template>
  <NuxtLayout>
    <main class="main-content">
      <NuxtPage 
        :transition="{
          name: 'anim',
          mode: 'out-in'
        }" 
      />
    </main>
  </NuxtLayout>
</template>


<style scoped>
  .main-content{
    margin: 0 auto;
    padding: 40px 16px;
    max-width: 1248px;
  }

  .anim-enter-active,
  .anim-leave-active {
    transition: all 0.3s ease;
  }

  .anim-enter-from,
  .anim-leave-to {
    opacity: 0;
    transform: translateY(20px);
  }
</style>