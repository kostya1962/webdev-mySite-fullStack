<script setup lang="ts">
import { ref } from 'vue'
import ResourceTable from '~/components/admin/ResourceTable.vue'

const tableRef = ref<InstanceType<typeof ResourceTable> | null>(null)

// expose fetchList to parent
// eslint-disable-next-line vue/no-expose-after-await
defineExpose({
  fetchList: () => tableRef.value?.fetchList()
})
</script>

<template>
  <ResourceTable ref="tableRef" resource="products" endpoint="/admin/products" />
</template>

