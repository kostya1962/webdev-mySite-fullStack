<script setup lang="ts">
import LayoutFooter from '~/components/Layout/LayoutFooter.vue';

</script>

<template>
    <div>
        <LayoutHeader />
        <slot />
        <LayoutFooter />
    </div>
</template>