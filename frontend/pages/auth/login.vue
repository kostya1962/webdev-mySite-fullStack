<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts">
import type { LoginResponse } from '~/interfaces/auth.interface';
import { useAuthStore } from '~/state/auth.state';
import { useFavoriteStore } from '~/state/favorite.state';

    const email = ref<string | undefined>();
    const password = ref<string | undefined>();
    const API_URL = useAPI(); 
    const authStore = useAuthStore();
    const favoriteStore = useFavoriteStore();

    const showToast = ref<boolean>(false)
    const toastMessage = ref<string>('')
    const toastState = ref<string>('')

    function openToast(message: string, state: string): void {
        toastMessage.value = message
        toastState.value = state
        showToast.value = true

        setTimeout(() => {
            showToast.value = false
        }, 3000)
    }

    async function login(){
        if (!email.value || !password.value) {
            openToast('Пожалуйста, заполните все поля', 'info')
            return
        }

        try{
            const data = await $fetch<LoginResponse>(API_URL + '/auth/login', {
                method: 'POST',
                body: {
                    email: email.value,
                    password: password.value,
                },
            });
            authStore.setToken(data.token);
            authStore.setEmail(data.user.email);
            authStore.setRole?.(data.user.role || 'user');
            await favoriteStore.restore(data.user.email);

            navigateTo('/account');
        
        } catch (error: any) {

            console.log('error', error);
        
            if (error?.status === 401) {
                openToast('Неверный email или пароль', 'error')
                return
            }
            openToast('Ошибка авторизации. Попробуйте позже', 'error')
        }
    }
</script>

<template>
    <div class="auth-root">
        <h1 class="page-title">Войти в аккаунт</h1>

        <div class="tabs">
            <button class="active_btn">Войти</button>
            <button class="passive_btn" @click.prevent="() => navigateTo('/auth/signup')">Зарегистрироваться</button>
        </div>

        <form action="" class="login-form">
            <div class="login-form__fileds">
                <InputFiled v-model="email" variant="gray" placeholder="Email" required/>
                <InputFiled v-model="password" variant="gray" type="password" placeholder="Пароль" required/>
            </div>
            <div class="login-form__actions">
                <ActionButton @click.stop.prevent="login" >
                    Вход
                </ActionButton>
                <NuxtLink to="/auth/restore">Забыли пароль?</NuxtLink>
            </div>
            
        </form>
        
        <InfoWindow
            v-if="showToast"
            :message="toastMessage"
            :duration="100000"
            :variant="toastState"
        />
    </div>
</template>

<style scoped>
    .auth-root{
        max-width: 520px;
        margin: 64px auto 270px auto;
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

    .login-form{
        display: flex;
        gap: 70px;
        flex-direction: column;
        margin: 0 auto;
        max-width: 500px;
        margin-top: 32px;
    }

    .login-form__fileds{
        display:flex;
        flex-direction: column;
        gap: 46px;
    }

    .login-form__actions{
        display: flex;
        flex-direction: column;
        gap: 22px;
    }

    .login-form__actions a{
        text-decoration: none;
        margin: 0 auto;
        color: var(--color-black);
    }

    .login-form__actions a:hover{
        color: var(--color-black-hover);
    }
</style>