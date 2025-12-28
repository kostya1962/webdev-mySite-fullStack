<!-- eslint-disable @typescript-eslint/no-unused-vars -->
<script setup lang="ts">
    import { useAuthStore } from '~/state/auth.state'


    interface ResourceItem {
        id?: number | string
        [key: string]: unknown
    }

    interface Category {
        id: number | string
        name: string
        [key: string]: unknown
    }

    interface CategoriesResponse {
        categories: Category[]
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

    // Переменные для работы с бэкапами
    interface BackupInfo {
        name: string
        path: string
        modTime: string
        size: number
    }

    const backups = ref<BackupInfo[]>([])
    const selectedBackup = ref('')
    const showBackupModal = ref(false)
    const backupLoading = ref(false)
    const backupMessage = ref('')

    // Загружаем категории (используется для таблицы товаров)
    const { data: categoriesData } = await useFetch<CategoriesResponse>(API_URL + '/categories');
    const categoriesList = computed(() => categoriesData.value?.categories ?? []);
    const resource = computed(() => props.resource);

    const hiddenColumns: Record<string, string[]> = {
        products: ['created_at', 'images', 'long_description', 'short_description', 'updated_at'],
        news: ['image'],
        users: ['created_at', 'updated_at', 'password', 'secret'],
    }

    // Колонки, которые отображаются в таблице (фильтруются на основе `cols` и `hiddenColumns`)
    const displayCols = ref<string[]>([])


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

                const raw = Object.keys(first).filter(k => k !== 'password' && k !== 'secret' )

                function orderCols(list: string[]) {
                    const preferred = ['id', 'name', 'title', 'email']
                    const result: string[] = []

                    // add id first if exists
                    if (list.includes('id')) {
                        result.push('id')
                    }

                    // then add first preferred after id (name/title/email)
                    for (const p of preferred) {
                        if (p === 'id') continue
                        if (list.includes(p) && !result.includes(p)) {
                            result.push(p)
                            break
                        }
                    }

                    // then add the rest in original order
                    for (const col of list) {
                        if (!result.includes(col)) result.push(col)
                    }

                    return result
                }

                cols.value = orderCols(raw)

                // Отфильтровать колонки для визуального отображения, но оставить `cols` целыми
                const hidden = hiddenColumns[resource.value] || []
                displayCols.value = cols.value.filter(c => !hidden.includes(c))
                headers.value = [...displayCols.value]
            } else {
                cols.value = ['id']
                displayCols.value = ['id']
                headers.value = ['id']
            }
        } catch (e) {
            console.error('fetchList error:', e)
        }
    }

    function formatCell(item: ResourceItem, col: string) {
        const v = item[col]
        if (v === null || v === undefined) return ''

        // Для товаров заменяем id_categories на имя категории
        if (col === 'id_categories') {
            const cat = categoriesList.value.find(c => String(c.id) === String(v))
            if (cat) return String(cat.name)
        }

        return typeof v === 'object' ? JSON.stringify(v) : String(v)
    }

    function openCreate() {
    editing.value = cols.value.reduce<ResourceItem>(
        (acc, key) => {
        // init image fields as empty array/string appropriately
        if (key === 'images') acc[key] = []
        else if (key === 'status' && resource.value === 'orders') acc[key] = 'оплачен'
        else if (key === 'created_at' || key === 'updated_at') acc[key] = new Date().toISOString()
        else acc[key] = ''
        return acc
        },
        {}
    )
    // ensure date fields are present even if cols exclude them
    editing.value['created_at'] = new Date().toISOString()
    editing.value['updated_at'] = new Date().toISOString()
    showForm.value = true
    }

    function openEdit(item: ResourceItem) {
    // keep all fields internally (id is needed for PUT), but we'll hide auto-fields in the form
    editing.value = JSON.parse(JSON.stringify(item))
    showForm.value = true
    }

    function closeForm() {
    showForm.value = false
    }

    async function save() {
    try {
        // prepare image fields for backend
        const payload = JSON.parse(JSON.stringify(editing.value))
        // Для продуктов отправляем массив изображений (AdminCreateProduct ожидает массив)

        if (resource.value === 'news' || resource.value === 'banners') {
            // ensure single image string
            if (Array.isArray(payload.image)) payload.image = payload.image[0] || ''
        }

        if (editing.value.id) {
            await $fetch(`${API_URL}${props.endpoint}/${editing.value.id}`, {
                method: 'PUT',
                headers: {
                    Authorization: `Bearer ${auth.token}`
                },
                // remove id only; allow created_at/updated_at if provided by admin
                body: (() => {
                    delete payload.id
                    return payload
                })()
            })
        } else {
            // remove id only; allow created_at/updated_at when creating
            delete payload.id

            await $fetch(`${API_URL}${props.endpoint}`, {
                method: 'POST',
                headers: {
                    Authorization: `Bearer ${auth.token}`
                },
                body: payload
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

    // Функции для работы с бэкапами БД
    async function loadBackups() {
        try {
            backupLoading.value = true
            const response = await $fetch<{ backups: BackupInfo[] }>(`${API_URL}/admin/backups`, {
                method: 'GET',
                headers: {
                    Authorization: `Bearer ${auth.token}`
                }
            })
            backups.value = response.backups || []
        } catch (e) {
            console.error('loadBackups error:', e)
            backupMessage.value = 'Ошибка при загрузке списка бэкапов'
        } finally {
            backupLoading.value = false
        }
    }

    async function createBackup() {
        try {
            backupLoading.value = true
            backupMessage.value = ''
            await $fetch(`${API_URL}/admin/backup`, {
                method: 'POST',
                headers: {
                    Authorization: `Bearer ${auth.token}`
                }
            })
            backupMessage.value = 'Бэкап успешно создан'
            await loadBackups()
            setTimeout(() => {
                backupMessage.value = ''
            }, 3000)
        } catch (e) {
            console.error('createBackup error:', e)
            backupMessage.value = 'Ошибка при создании бэкапа'
        } finally {
            backupLoading.value = false
        }
    }

    async function restoreBackup() {
        if (!selectedBackup.value) {
            backupMessage.value = 'Пожалуйста, выберите бэкап для восстановления'
            return
        }

        if (!confirm('Вы уверены? Текущая база данных будет заменена на выбранный бэкап. Автоматический бэкап текущего состояния будет сохранён.')) {
            return
        }

        try {
            backupLoading.value = true
            backupMessage.value = ''
            await $fetch(`${API_URL}/admin/restore`, {
                method: 'POST',
                headers: {
                    Authorization: `Bearer ${auth.token}`
                },
                body: {
                    backupName: selectedBackup.value
                }
            })
            backupMessage.value = 'База данных успешно восстановлена'
            selectedBackup.value = ''
            await loadBackups()
            setTimeout(() => {
                backupMessage.value = ''
                showBackupModal.value = false
            }, 2000)
        } catch (e) {
            console.error('restoreBackup error:', e)
            backupMessage.value = 'Ошибка при восстановлении бэкапа'
        } finally {
            backupLoading.value = false
        }
    }

    function openBackupModal() {
        selectedBackup.value = ''
        backupMessage.value = ''
        loadBackups()
        showBackupModal.value = true
    }

    function closeBackupModal() {
        showBackupModal.value = false
        selectedBackup.value = ''
        backupMessage.value = ''
    }

    // Computed property для фильтрации полей формы (исключаем 'id')
    const editingFields = computed(() => {
        return Object.entries(editing.value).filter(([key]) => key !== 'id')
    })

    watch(
        () => props.resource,
        () => fetchList(),
        { immediate: true }
    )

    // Делаем fetchList доступным из родительского компонента через ref
    // eslint-disable-next-line vue/no-expose-after-await
    defineExpose({
        fetchList
    })

    async function uploadFileToServer(file: File) {
        try {
            const fd = new FormData()
            fd.append('file', file)

            const res = await fetch(`${API_URL}/admin/upload/${resource.value}`, {
                method: 'POST',
                headers: {
                    Authorization: `Bearer ${auth.token}`
                },
                body: fd
            })

            if (!res.ok) throw new Error('upload failed')
            const data = await res.json()
            return data.path || data.url || ''
        } catch (e) {
            console.error('upload error', e)
            return ''
        }
    }

    async function onFilesChange(e: Event, key: string) {
        const input = e.target as HTMLInputElement
        if (!input.files || input.files.length === 0) return

        const files = Array.from(input.files)
        const uploaded: string[] = []

        for (const f of files) {
            const p = await uploadFileToServer(f)
            if (p) uploaded.push(p)
        }

        if (key === 'images') {
            // ensure array; if stored as JSON string, parse it
            if (!Array.isArray(editing.value[key])) {
                try {
                    const parsed = JSON.parse(String(editing.value[key] || '[]'))
                    editing.value[key] = Array.isArray(parsed) ? parsed : []
                } catch (e) {
                    editing.value[key] = []
                }
            }
            editing.value[key] = [...(editing.value[key] as ResourceItem[]), ...uploaded]
        } else {
            // single image
            editing.value[key] = uploaded[0] || ''
        }
    }

    function toLocalInput(val: unknown) {
        if (!val) return ''
        try {
            const d = new Date(String(val))
            const pad = (n: number) => String(n).padStart(2, '0')
            return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
        } catch (e) {
            return ''
        }
    }

    function onDateInput(e: Event, key: string) {
        const input = e.target as HTMLInputElement
        if (!input.value) {
            editing.value[key] = ''
            return
        }
        // input.value is local datetime like '2025-12-27T10:20' -> convert to ISO
        const d = new Date(input.value)
        editing.value[key] = d.toISOString()
    }
</script>

<template>
<div>
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:12px">
        <div>
            <button @click="openCreate">Создать</button>
            <button @click="fetchList">Обновить</button>
        </div>
        <div style="display:flex;gap:8px">
            <button :disabled="backupLoading" @click="createBackup">Создать бэкап БД</button>
            <button :disabled="backupLoading" @click="openBackupModal">Загрузить бэкап</button>
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
        <td v-for="col in displayCols" :key="col">
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

        <div v-for="[key] in editingFields" :key="key" class="field">
        <label>{{ key }}</label>

        <template v-if="key === 'category_id' && resource === 'products'">
            <select v-model="editing[key]">
                <option value="">Выбрать категорию</option>
                <option v-for="c in categoriesList" :key="c.id" :value="c.id">{{ c.name }}</option>
            </select>
        </template>

        <template v-else-if="(key === 'images' && resource === 'products') || (key === 'image' && (resource === 'news' || resource === 'banners'))">
            <div>
                <input :multiple="key === 'images'" type="file" @change="onFilesChange($event, key)" />
                <div style="margin-top:8px">
                    <template v-if="key === 'images'">
                        <div v-if="Array.isArray(editing[key])">
                            <div v-for="p in editing[key]" :key="p">{{ p }}</div>
                        </div>
                        <div v-else>{{ editing[key] }}</div>
                    </template>
                    <template v-else>
                        <div>{{ editing[key] }}</div>
                    </template>
                </div>
            </div>
        </template>

        <template v-else-if="key === 'status' && resource === 'orders'">
            <select v-model="editing[key]">
                <option value="оплачен">оплачен</option>
                <option value="отправлен">отправлен</option>
                <option value="доставлен">доставлен</option>
                <option value="закрыт">закрыт</option>
            </select>
        </template>

        <template v-else-if="key === 'created_at' || key === 'updated_at'">
            <input :value="toLocalInput(editing[key])" type="datetime-local" @input="onDateInput($event, key)" />
        </template>

        <template v-else>
            <input v-model="editing[key]" />
        </template>
        </div>

        <div style="margin-top:12px">
        <button @click="save">Сохранить</button>
        <button @click="closeForm">Отмена</button>
        </div>
    </div>
    </div>

    <!-- Модальное окно для восстановления бэкапа -->
    <div v-if="showBackupModal" class="modal">
    <div class="modal-body">
        <h3>Восстановление бэкапа БД</h3>

        <div v-if="backupMessage" :style="{ color: backupMessage.includes('Ошибка') ? 'red' : 'green', marginBottom: '12px' }">
            {{ backupMessage }}
        </div>

        <div class="field">
            <label>Выбрать бэкап</label>
            <select v-model="selectedBackup" :disabled="backupLoading">
                <option value="">-- Выберите бэкап --</option>
                <option v-for="backup in backups" :key="backup.name" :value="backup.name">
                    {{ backup.name }} ({{ backup.modTime }}, {{ (backup.size / 1024 / 1024).toFixed(2) }} MB)
                </option>
            </select>
        </div>

        <div style="margin-top:12px">
            <button :disabled="backupLoading" @click="restoreBackup">
                {{ backupLoading ? 'Обработка...' : 'Восстановить' }}
            </button>
            <button :disabled="backupLoading" @click="closeBackupModal">Отмена</button>
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

.field select {
    padding: 6px;
    border: 1px solid #ddd;
}
</style>
