<script setup lang="ts">
import { useFavoriteStore } from '~/state/favorite.state';

const { id, isShown, curcolor } = defineProps<{
    id: number;
    isShown: boolean;
    curcolor: string;
}>(); 
console.log(`${curcolor} in AddFavorite.vue`);
const favoriteState = useFavoriteStore();
</script>

<template>
    <button v-show="isShown || favoriteState.isFavorite(id)" class="fav" @click.stop.prevent="() => favoriteState.toggleFavorite(id)">
        <Icon name="icons:favorite-add" size="18px" color="curcolor"/>
    </button>
</template>

<style scoped>
    .fav{
        background: none;
        border: none;
        cursor: pointer;
    }
</style>