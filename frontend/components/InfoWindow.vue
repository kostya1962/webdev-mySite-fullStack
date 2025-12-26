<script setup lang="ts">
    interface Props {
        message?: string
        duration?: number
        variant?: string
    }

    const props = withDefaults(defineProps<Props>(), {
        message: 'Операция выполнена успешно',
        duration: 3000,
        variant: 'success'
    })

    const visible = ref<boolean>(false)
    let timer: number | null = null
    const state = ref<string>(props.variant ?? 'success')

    onMounted(() => {
        visible.value = true

        timer = window.setTimeout(() => {
            visible.value = false
        }, props.duration)
    });

    onBeforeUnmount(() => {
        if (timer !== null) {
            clearTimeout(timer)
        }
    })
</script>

<template>
    <Transition name="toast">
        <div v-if="visible" class="toast">
            <Icon v-if="state === 'success'" name="icons:access" size="21px"/>
            <Icon v-else-if="state === 'error'" name="icons:error" size="21px"/>
            <Icon v-else name="icons:info" size="21px"/>
            <span class="text">{{ message }}</span>
        </div>
    </Transition>
</template>

<style scoped>
    .toast {
        position: fixed;
        left: 50%;
        bottom: 24px;
        transform: translateX(-50%) translateY(0);
        background: var(--color-gray);
        color: var(--color-black);
        padding: 14px 20px;
        border-radius: 8px;
        display: flex;
        align-items: center;
        gap: 10px;
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
        font-size: 14px;
        z-index: 9999;
    }

    /* Общая анимация */
    .toast-enter-active,
    .toast-leave-active {
    transition: opacity 0.35s ease, transform 0.35s ease;
    }

    /* Появление */
    .toast-enter-from {
    opacity: 0;
    transform: translateX(-50%) translateY(30px);
    }

    .toast-enter-to {
    opacity: 1;
    transform: translateX(-50%) translateY(0);
    }

    /* Исчезновение */
    .toast-leave-from {
    opacity: 1;
    transform: translateX(-50%) translateY(0);
    }

    .toast-leave-to {
    opacity: 0;
    transform: translateX(-50%) translateY(30px);
    }

</style>