<script setup lang="ts">
import CheckboxFiled from './CheckboxFiled.vue';
import InputArea from './InputArea.vue';
import InputRating from './InputRating.vue';

const rating = ref<number>(0)
const isSubmitted = ref(false)

</script>

<template>
    <div class="rev">
        <div>
            Добавить отзыв
        </div>
        <div>
            Обязательные поля помечены *
        </div>
    <form >
        <InputFiled placeholder="Ваше имя*" required/>
        <InputArea placeholder="Отзыв*" required/>
        <CheckboxFiled placeholder="Сохранить данные для следующих отзывов" :data="false" />
        <div>
            Рейтинг*
            <InputRating v-model="rating"/>
        </div>
        <ActionButton type="submit" :disabled="isSubmitted">
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
</style>