<script setup lang="ts">
    const email = ref<string>('')

    const toastMessage = ref<string>('')
    const toastState = ref<'success' | 'error'>('success')
    const showToast = ref<boolean>(false)

    function openToast(message: string, state: 'success' | 'error'): void {
        toastMessage.value = message
        toastState.value = state

        showToast.value = true

        setTimeout(() => {
            showToast.value = false
        }, 3000)
    }

    function restore(): void {
        if (!email.value.trim()) {
            openToast('Введите email для восстановления пароля', 'error')
            return
        }

        console.log(`Отправка письма на почту по адресу ${email.value}`)

        openToast('Ссылка для восстановления отправлена на почту', 'success')
        email.value = ''
    }
</script>
<template>
    <div class="restore-root">
        <h1 class="page-title">Забыли пароль?</h1>

        <div class="restore-text">
            <p>
                Если вы забыли пароль, то введите свой email и мы отправим вам ссылку
                на восстановление
            </p>
        </div>

        <form class="restore-form" @submit.prevent="restore">
        <div class="restore-form__fileds">
            <InputFiled
            v-model="email"
            variant="gray"
            placeholder="Email"
            />
        </div>

        <div class="restore-form__actions">
            <div class="full-wrap">
            <ActionButton color="primary" @click.stop.prevent="restore">
                Сбросить пароль
            </ActionButton>
            </div>
        </div>
        </form>

    
        <InfoWindow
            v-if="showToast"
            :message="toastMessage"
            :variant="toastState"
            :duration="10000"
        />
    </div>
</template>

<style scoped>
    .restore-root{
        max-width: 520px;
        margin: 64px auto 270px auto;
        padding: 0 16px;
    }

    .page-title{
        text-align: center;
        font-size: 28px;
        margin: 0;
    }

    .restore-text{
        font-size: 20px;
        color: var(--color-black);
        text-align: center;
        margin-top: 60px;
        margin-bottom: 80px;
    }


    .restore-form{
        display: flex;
        gap: 24px;
        flex-direction: column;
        margin-top: 32px;
        cursor: pointer;
    }

    .restore-form__fileds{
        display:flex;
        flex-direction: column;
        gap: 28px;
    }

    .restore-form__actions{
        display:flex;
        justify-content: center;
        margin-top: 80px;
        font-weight: 75%;
    }

    .full-wrap{
        width: 100%;
        max-width: 520px;
    }

    .full-wrap :deep(button) {
        width: 100%;
        display: block;
    }



</style>