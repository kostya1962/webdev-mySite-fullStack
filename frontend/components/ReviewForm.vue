<script setup lang="ts">
import CheckboxFiled from './CheckboxFiled.vue';
import InputArea from './InputArea.vue';
import InputRating from './InputRating.vue';
import InputFiled from './InputFiled.vue';

const emit = defineEmits<{
    (e: 'submit-review', payload: { name: string; text: string; rating: number; save_to_account?: boolean }): void
}>()

const name = ref<string>('')
const text = ref<string>('')
const rating = ref<number>(0)
const saveData = ref<boolean>(false)
const isSubmitted = ref(false)

function onSubmit(e: Event) {
    e.preventDefault()
    if (!name.value || !text.value || rating.value <= 0) {
        return
    }
    isSubmitted.value = true
    emit('submit-review', {
        name: name.value.trim(),
        text: text.value.trim(),
        rating: Math.floor(rating.value),
        save_to_account: saveData.value,
    })
}

</script>

<template>
    <div class="rev">
        <div>
            Добавить отзыв
        </div>
        <div>
            Обязательные поля помечены *
        </div>
    <form @submit="onSubmit">
        <div class="input__area">
            <InputFiled v-model="name" placeholder="Ваше имя*" variant="gray" required/>
            <InputArea v-model="text" placeholder="Отзыв*" variant="gray" required/>
            <CheckboxFiled v-model="saveData">
                Сохранить данные для следующих отзывов
            </CheckboxFiled>
            <div class="input__rating">
                <p>Рейтинг*</p>
                <InputRating v-model="rating"/>
            </div>
        </div>
        <ActionButton class="input__button" type="submit" :disabled="isSubmitted">
            Отправить отзыв
        </ActionButton>
    </form>
    </div>
</template>

<style scoped>
.rev{
    display: flex;
    flex-direction: column;
    gap: 30px;
    margin-top: 16px;
}

.input__area{
    display: flex;
    flex-direction: column;
    gap: 50px;
    margin-bottom: 55px;
}

.input__rating{
    display: flex;
    flex-direction: column;
    gap: 30px;
    font-size: 14px;
    line-height: 20px;
    color: #333333;
}

.input__button{
    width: 200px;
    height: 50px;
    
}
</style>