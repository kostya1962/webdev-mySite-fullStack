<script setup lang="ts">
import { ref } from 'vue'
import AdminResource from '~/components/admin/ResourceTable.vue'

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

// Ref для дочернего компонента, чтобы вызывать его методы
const resourceTableRef = ref<InstanceType<typeof AdminResource>>()

// Функция reload вызывает fetchList в дочернем компоненте
// для обновления списка после создания/редактирования/удаления записи
function reload() {
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
      <AdminResource
        ref="resourceTableRef"
        :resource="active"
        :endpoint="`/admin/${active}`"
        @saved="reload"
      />
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
