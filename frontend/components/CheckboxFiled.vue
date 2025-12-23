<script setup lang="ts">
const data = defineModel<boolean>();
</script>

<template>
    <div>
        <div>
            <input v-model="data"  type="checkbox" />
            <span><slot /></span>
        </div>
    </div>
</template>