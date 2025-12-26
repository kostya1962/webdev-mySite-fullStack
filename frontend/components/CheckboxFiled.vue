<script setup lang="ts">
const data = defineModel<boolean>();
</script>

<template>
    <div>
        <div class="check">
            <input v-model="data"  type="checkbox" placeholder="Введите данные"/>
            <span><slot /></span>
        </div>
    </div>
</template>

<style scoped>

    .check{
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 14px;
        line-height: 20px;
        color: #333333;
    }
    input[type="checkbox"]{
        width: 16px;
        height: 16px;
        border: 1px solid #CCCCCC;
        border-radius: 4px;
        background-color: #FFFFFF;
        cursor: pointer;
    }
</style>
