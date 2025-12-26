<script setup lang="ts">
    import { ref, watch, computed } from 'vue'
    import { useAuthStore } from '~/state/auth.state'


    interface ResourceItem {
        id?: number | string
        [key: string]: unknown
    }

    const props = defineProps<{
        resource: string
        endpoint: string
    }>()


    const auth = useAuthStore()

    const items = ref<ResourceItem[]>([])
    const editing = ref<ResourceItem>({})
    const showForm = ref(false)

    const cols = ref<string[]>([])
    const headers = ref<string[]>([])
    const API_URL = useAPI();

    const title = computed(() => props.resource)


    async function fetchList() {
        try {
            const data = await $fetch<ResourceItem[] | { data: ResourceItem[] }>(
            `${API_URL}${props.endpoint}`,
            {
                method: 'GET',
                headers: {
                Authorization: `Bearer ${auth.token}`
                }
            }
            )

            const arr = Array.isArray(data) ? data : data.data || []
            items.value = arr

            if (arr.length > 0) {
            const first = arr[0]
                if (!first) return

                cols.value = Object.keys(first).filter(
                    k => k !== 'password' && k !== 'secret'
                )
                headers.value = [...cols.value]
            } else {
                cols.value = ['id']
                headers.value = ['id']
            }
        } catch (e) {
            console.error('fetchList error:', e)
        }
    }

    function formatCell(item: ResourceItem, col: string) {
    const v = item[col]
    if (v === null || v === undefined) return ''
    return typeof v === 'object' ? JSON.stringify(v) : String(v)
    }

    function openCreate() {
    editing.value = cols.value.reduce<ResourceItem>(
        (acc, key) => {
        acc[key] = ''
        return acc
        },
        {}
    )
    showForm.value = true
    }

    function openEdit(item: ResourceItem) {
    editing.value = JSON.parse(JSON.stringify(item))
    showForm.value = true
    }

    function closeForm() {
    showForm.value = false
    }

    async function save() {
    try {
        if (editing.value.id) {
        await $fetch(`${API_URL}${props.endpoint}/${editing.value.id}`, {
            method: 'PUT',
            headers: {
            Authorization: `Bearer ${auth.token}`
            },
            body: editing.value
        })
        } else {
        await $fetch(`${API_URL}${props.endpoint}`, {
            method: 'POST',
            headers: {
            Authorization: `Bearer ${auth.token}`
            },
            body: editing.value
        })
        }

        showForm.value = false
        await fetchList()
    } catch (e) {
        console.error('save error:', e)
    }
    }

    async function remove(item: ResourceItem) {
    if (!item.id) return
    if (!confirm('Удалить запись?')) return

    try {
        await $fetch(`${API_URL}${props.endpoint}/${item.id}`, {
        method: 'DELETE',
        headers: {
            Authorization: `Bearer ${auth.token}`
        }
        })
        await fetchList()
    } catch (e) {
        console.error('remove error:', e)
    }
    }


    watch(
        () => props.resource,
        () => fetchList(),
        { immediate: true }
    )

    // Делаем fetchList доступным из родительского компонента через ref
    defineExpose({
        fetchList
    })
</script>

<template>
<div>
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:12px">
    <h2>{{ title }}</h2>
    <div>
        <button @click="openCreate">Создать</button>
        <button @click="fetchList">Обновить</button>
    </div>
    </div>

    <table v-if="items.length" class="resource-table">
    <thead>
        <tr>
        <th v-for="h in headers" :key="h">{{ h }}</th>
        <th>Действия</th>
        </tr>
    </thead>
    <tbody>
        <tr v-for="item in items" :key="item.id">
        <td v-for="col in cols" :key="col">
            {{ formatCell(item, col) }}
        </td>
        <td>
            <button @click="openEdit(item)">Ред.</button>
            <button @click="remove(item)">Удал.</button>
        </td>
        </tr>
    </tbody>
    </table>

    <p v-else>Записей нет</p>

    <div v-if="showForm" class="modal">
    <div class="modal-body">
        <h3>{{ editing.id ? 'Редактировать' : 'Создать' }}</h3>

        <div v-for="(val, key) in editing" :key="key" class="field">
        <label>{{ key }}</label>
        <input v-model="editing[key]" />
        </div>

        <div style="margin-top:12px">
        <button @click="save">Сохранить</button>
        <button @click="closeForm">Отмена</button>
        </div>
    </div>
    </div>
</div>
</template>



<style scoped>
.resource-table {
    width: 100%;
    border-collapse: collapse;
}
.resource-table th,
.resource-table td {
    border: 1px solid #eee;
    padding: 8px;
    text-align: left;
    }

.modal {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.4);
    display: flex;
    align-items: center;
    justify-content: center;
}

.modal-body {
    background: #fff;
    padding: 16px;
    min-width: 400px;
    border-radius: 6px;
}

.field {
    margin-bottom: 8px;
    display: flex;
    flex-direction: column;
}

.field input {
    padding: 6px;
    border: 1px solid #ddd;
}
</style>
