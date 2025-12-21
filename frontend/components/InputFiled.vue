<script setup lang="ts">
const { variant="black" } = defineProps<{
    variant?: "gray" | "black";
}>();
const data = defineModel<string>();
</script>

<template>
    <input v-model="data" class="input" :class="variant"/>
</template>


<style scoped>
.input{
    padding: 12px 0;
    background: none;
    border: none;
    color: var(--color-black);
    width: 100%;
}

.input::placeholder{
    color: var(--color-dark-gray);
}

.gray{
    border-bottom: 1px solid var(--color-gray);
}

.black{
    border-bottom: 1px solid var(--color-black);
}
</style>