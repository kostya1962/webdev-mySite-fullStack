<script setup lang="ts">
import { useDebounceFn } from '@vueuse/core';
import type { GetCategoryResponse } from '~/interfaces/category.interface';
import type { GetProductsResponse } from '~/interfaces/product.interface';

useSeoMeta({
    title: 'Каталог товаров',
    description: 'Каталог товаров интернет магазина Shopper с  ювелирными изделиями',
    ogDescription: 'Каталог товаров интернет магазина Shopper с  ювелирными изделиями',
});

const API_URL = useAPI();
const route = useRoute(); // даёт текущий роут
const router = useRouter(); // управляет состоянием роутера (запроса)


const select = ref(route.query.select?.toString() ?? "");
const search = ref(route.query.search?.toString() ?? "");

const changeRoute = useDebounceFn((select, search) => {
    router.replace(
            { 
                query: 
                    { 
                        select: select.value,
                        search: search.value,
                    } 
            }
        )
}, 100); // будет ждать 100 мс с последнего вызова 


watch([select, search], () => {
    changeRoute(select, search)
});

const query = computed(() => (
    {
        limit: route.query.limit ?? 20,
        offset: route.query.offset ?? 0,
        category_id: route.query.select || undefined,
        search: route.query.search || undefined,
    }
));

const {data} = await useFetch<GetCategoryResponse>(API_URL + '/categories'); 

const defaultCategories = {
    value: '',
    label: 'Категории'
};
const categories = computed(() => {
    return data.value ?
            data.value?.categories.map((c) => ({
                value: c.id.toString(),
                label: c.name,
            })).concat([defaultCategories])
        : [defaultCategories];
});

const { data: productsData } = await useFetch<GetProductsResponse>(
    API_URL + '/products', 
    {
        query,
        key: 'get-products',
    }
); 

</script>

<template>
    <div>
        <h1 class="left">Каталог товаров</h1>
        <div class="catalog">
            <div class="catalog__filter">
                <div class="catalog__search">
                    <InputFiled v-model="search" variant="gray" placeholder="Поиск..." />
                    <Icon name="icons:search" size="19" />
                </div>
                
                <SelectFiled 
                    v-model="select"
                    :options="categories"
                />
            </div>
            <div class="catalog__list">
                <CatalogCard 
                    v-for="product in productsData?.products" 
                    :key="product.id" 
                    v-bind="product"
                />
            </div>
        </div>

    </div>
</template>

<style scoped>
    .catalog{
        display: flex;
    }

    .catalog__filter{
        width: 260px;
        margin-right: 45px;
        display: flex;
        flex-direction: column;
        gap: 24px;
    }

    .catalog__list{
        display: grid;
        width: 100%;
        gap: 64px 12px;
        grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    }

    .catalog__search{
        position: relative;
    }

    .catalog__search .iconify {
        position: absolute;
        top: 11px;
        right: 9px;
    }
</style>

