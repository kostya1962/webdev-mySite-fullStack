<script setup lang="ts">
import { useFavoriteStore } from '~/state/favorite.state';

const { id, isShown } = defineProps<{
    id: number;
    isShown: boolean;
}>(); 

const favoriteState = useFavoriteStore();
</script>

<template>
    <button v-show="isShown || favoriteState.isFavorite(id)" class="fav" @click.stop.prevent="() => favoriteState.toggleFavorite(id)">
        <Icon name="icons:favorite-add" size="18px"/>
    </button>
</template>

<style scoped>
    .fav{
        background: none;
        border: none;
        cursor: pointer;
    }
</style>