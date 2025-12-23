<script setup lang="ts">
import { useFavoriteStore } from '~/state/favorite.state';

const { id } = defineProps<{
    id: number;
}>(); 

const favoriteState = useFavoriteStore();
</script>

<template>
    <button class="fav" @click.stop.prevent="() => favoriteState.toggleFavorite(id)">
        <Icon name="icons:empty-fav" size="18px" />
    </button>
</template>

<style scoped>
    .fav{
        background: none;
        border: none;
        cursor: pointer;
    }
</style>