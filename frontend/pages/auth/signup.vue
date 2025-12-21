<script setup lang="ts">
import type { LoginResponse } from '~/interfaces/auth.interface'
import { useAuthStore } from '~/state/auth.state'
import { useFavoriteStore } from '~/state/favorite.state'

const API_URL = useAPI()
const authStore = useAuthStore()
const favoriteStore = useFavoriteStore()

const email = ref<string | undefined>()
const password = ref<string | undefined>()
const password2 = ref<string | undefined>()
const agree = ref(false)
const error = ref<string | null>(null)
const loading = ref(false)

const canSubmit = computed(() => {
    return !!email.value && !!password.value && password.value === password2.value && agree.value && !loading.value
})

async function signup() {
    error.value = null
    if (!canSubmit.value) return
    loading.value = true
    try {
        const data = await $fetch<LoginResponse>(API_URL + '/auth/register', {
            method: 'POST',
            body: {
                email: email.value,
                password: password.value,
            },
        })

        authStore.setToken(data.token)
        authStore.setEmail(data.user.email)
        await favoriteStore.restore(data.user.email)

        navigateTo('/account')
    } catch (e: any) {
        error.value = e?.data?.error || e?.message || 'Ошибка регистрации'
    } finally {
        loading.value = false
    }
}
</script>

<template>
    <div class="signup-root">
        <h1 class="page-title">Мой аккаунт</h1>

        <div class="tabs">
            <button class="passive_btn" @click.prevent="() => navigateTo('/auth/login')">Войти</button>
            <button class="active_btn">Зарегистрироваться</button>
        </div>

        <form class="signup-form" @submit.prevent="signup">
            <div class="signup-form__fileds">
                <InputFiled v-model="email" variant="gray" placeholder="Email" />
                <InputFiled v-model="password" variant="gray" type="password" placeholder="Пароль" />
                <InputFiled v-model="password2" variant="gray" type="password" placeholder="Повторите пароль" />
            </div>

            <label class="agree">
                <input type="checkbox" v-model="agree" />
                <span>Согласен на обработку персональных данных</span>
            </label>

            <div class="signup-form__actions">
                <div class="full-wrap">
                    <ActionButton color="primary" @click.stop.prevent="signup">
                        Зарегистрироваться
                    </ActionButton>
                </div>
            </div>

            <div v-if="error" class="form-error">{{ error }}</div>
        </form>
    </div>
</template>

<style scoped>
    .signup-root{
        max-width: 520px;
        margin: 64px auto 0 auto;
        padding: 0 16px;
    }

    .page-title{
        text-align: center;
        font-size: 28px;
        margin: 0;
    }

    .tabs{
        display: flex;
        gap: 12px;
        justify-content: space-between;
        margin: 18px auto 0 auto;
        background-color: var(--color-gray);
        border-radius: 4px;
        width: 100%;
        height: 50px;
    }

    .active_btn{
        background-color: var(--color-white);
        border: none;
        margin: 5px 8px;
        border-radius: 4px;
        width: 50%;
    }

    .passive_btn{
        background-color: transparent;
        border: none;
        margin: 5px 8px;
        border-radius: 4px;
        width: 50%;
        cursor: pointer;
    }

    .signup-form{
        display: flex;
        gap: 24px;
        flex-direction: column;
        margin-top: 32px;
        cursor: pointer;
    }

    .signup-form__fileds{
        display:flex;
        flex-direction: column;
        gap: 28px;
    }

    .agree{
        display: flex;
        align-items: center;
        gap: 12px;
        font-size: 14px;
        color: var(--color-dark-gray);
    }

    .signup-form__actions{
        display:flex;
        justify-content: center;
        margin-top: 10px;
    }

    .full-wrap{
        width: 100%;
        max-width: 520px;
    }

    .full-wrap ::v-deep button{
        width: 100%;
        display: block;
    }

    .form-error{
        color: red;
        text-align: center;
        margin-top: 12px;
    }

</style>
