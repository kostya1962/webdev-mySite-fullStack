<script setup lang="ts">
useSeoMeta({
    title: 'Контакты',
    description: 'Контактная информация интернет магазина Shopper',
    ogDescription: 'Контактная информация интернет магазина Shopper',
});
</script>
<template>
     <div class="block">
            <div class="contacts__title">
                <h1>Контакты</h1>
            </div>

            <div class="contacts">
                <div class="contacts-info">
                    <p class="text">
                        <strong>Адрес:</strong><br />
                        г. Чебоксары, Московский проспект, 15.
                    </p>

                    <p class="text">
                        <strong>Телефоны:</strong><br />
                        +7 (495) 123-45-67<br />
                        +7 (999) 765-43-21
                    </p>

                    <p class="text">
                        <strong>Режим работы:</strong><br />
                        Пн–Пт: 10:00–19:00<br />
                        Сб–Вс: выходной
                    </p>
                </div>

                <div class="map">

                        <iframe  
                            src="https://yandex.ru/map-widget/v1/?from=mapframe&ll=47.224290%2C56.145715&mode=search&oid=1725040712&ol=biz&source=mapframe&um=constructor%3A9d9e6a5f3f6d2e3f3b4e5f6a7b8c9d0e&utm_source=mapframe&z=16.98" 
                            width="560" 
                            height="400" 
                            frameborder="1" 
                            allowfullscreen="true" 
                            style="position:relative;">
                        </iframe>

                </div>
            </div>
        </div>
</template>

<style scoped>
    .block {
        margin-bottom: 200px;

    }
    .contacts__title {
        margin-top: 10px;
        margin-bottom: 200px;
        font-weight: 900;
    }

    .contacts {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 40px;
    }

    .contacts-info .text {
        margin-bottom: 30px;
    }

    .text {
        max-width: 100%;
        font-size: 20px;
        line-height: 1.7;
        color: #555;
        margin-bottom: 24px;
    }

    .map {
        width: 100%;
        border-radius: 12px;
        overflow: hidden;
    }

    /* Адаптив */
    @media (max-width: 768px) {
        .contacts {
            grid-template-columns: 1fr;
        }
    }

</style>