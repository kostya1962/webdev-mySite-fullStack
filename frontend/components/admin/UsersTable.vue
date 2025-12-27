<script setup lang="ts">
import { ref } from 'vue'
import ResourceTable from '~/components/admin/ResourceTable.vue'

const tableRef = ref<InstanceType<typeof ResourceTable> | null>(null)
defineExpose({
  fetchList: () => tableRef.value?.fetchList()
})
</script>

<template>
  <ResourceTable ref="tableRef" resource="users" endpoint="/admin/users" />
</template>
