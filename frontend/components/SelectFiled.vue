<script setup lang="ts">
const model = defineModel<string>();
const { options } = defineProps<{
    options: {label: string; value: string}[];
}>();
</script>

<template>
    
    <div class="select-wrapper">
        <select v-model="model" class="select">
            <option v-for="item in options" :key="item.value" :value="item.value">
                {{ item.label }}
            </option>
        </select>
        <Icon name="icons:select-dawn" class="select-dawn" size="10px"/>
    </div>
</template>

<style scoped>
    .select-wrapper{
        position: relative;
    }

    .select-dawn{
        position: absolute;
        right: 12px;
        top: 50%;
        transform: translateY(-50%);
        pointer-events: none;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        color: inherit;
    }


    .select {
        appearance: none;
        background-color: transparent;
        border: 1px solid var(--color-gray);
        padding: 10px 40px 10px 12px;
        box-sizing: border-box;
        margin: 0;
        width: 100%;
        font-family: inherit;
        font-size: 14px;
        line-height: 22px;

    }
</style>