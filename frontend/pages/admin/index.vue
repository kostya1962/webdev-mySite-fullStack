<script setup lang="ts">
import { ref } from 'vue'
import ProductsTable from '~/components/admin/ProductsTable.vue'
import CategoriesTable from '~/components/admin/CategoriesTable.vue'
import OrdersTable from '~/components/admin/OrdersTable.vue'
import NewsTable from '~/components/admin/NewsTable.vue'
import BannersTable from '~/components/admin/BannersTable.vue'
import UsersTable from '~/components/admin/UsersTable.vue'

    useSeoMeta({
        title: 'Админ-панель',
        description: 'Админ-панель интернет магазина Shopper',
        ogDescription: 'Админ-панель интернет магазина Shopper',
    });

const tabs = ['products','categories','orders','news','banners','users']
const tabNames: Record<string,string> = {
  products: 'Товары',
  categories: 'Категории',
  orders: 'Заказы',
  news: 'Новости',
  banners: 'Баннеры',
  users: 'Пользователи'
}
const active = ref('products')

// Interface for table components
interface TableComponent {
  fetchList(): void
}

// Ref для дочернего компонента, чтобы вызывать его методы
const resourceTableRef = ref<TableComponent | null>(null)

// Функция reload вызывает fetchList в дочернем компоненте
// для обновления списка после создания/редактирования/удаления записи
function _reload() {
  resourceTableRef.value?.fetchList()
}
</script>

<template>
  <div class="admin-page">
    <h1>Админ-панель</h1>

    <div class="tabs">
      <button
        v-for="tab in tabs"
        :key="tab"
        :class="{ active: tab === active }"
        @click="active = tab"
      >
        {{ tabNames[tab] }}
      </button>
    </div>

    <div class="panel">
      <ProductsTable v-if="active === 'products'" ref="resourceTableRef" />
      <CategoriesTable v-if="active === 'categories'" ref="resourceTableRef" />
      <OrdersTable v-if="active === 'orders'" ref="resourceTableRef" />
      <NewsTable v-if="active === 'news'" ref="resourceTableRef" />
      <BannersTable v-if="active === 'banners'" ref="resourceTableRef" />
      <UsersTable v-if="active === 'users'" ref="resourceTableRef" />
    </div>
  </div>
</template>

<style scoped>
.admin-page{
  max-width:1200px;
  margin:32px auto;
  padding:0 16px}
.tabs{display:flex;gap:8px;margin-bottom:16px
}
.tabs button{padding:8px 12px;border:1px solid #ddd;background:#fff}
.tabs button.active{background:#111;color:#fff}
.panel{background:#fff;padding:16px;border:1px solid #eee}
</style>
