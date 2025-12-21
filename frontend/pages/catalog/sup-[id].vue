<script setup lang="ts">
import GallerayProd from '~/components/GallerayProd.vue';
import type { Product } from '~/interfaces/product.interface';




const route = useRoute(); //извлекаю текущий маршрут (путь и параметры)
const API_URL = useAPI();


const {data} = await useFetch<{product: Product}>(
    API_URL + '/products/' + route.params.id
);

useSeoMeta({
    title: `Купить ${data.value?.product.name}`,
    description: data.value?.product.short_description,
});

</script>

<template>
    <div>
        <div>
            <GallerayProd :images="data?.product.images ?? []" />
        </div>
    </div>
</template>