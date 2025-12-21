<script setup lang="ts">

</script>

<template>

    <div class="footer">

        <hr />

        <div class="footer__top">
            <ul>
                <li><NuxtLink to="#">КОНТАКТЫ</NuxtLink></li>
                <li><NuxtLink to="#">УСЛВИЯ ПОКУПКИ </NuxtLink></li>
                <li><NuxtLink to="#">ДОАСТВКА И ВОЗВРАТ</NuxtLink></li>
            </ul>
            <div class="footer__form">
                <InputFiled placeholder="Ваш email для акций и предложений" variat="black"/>
                <button class="footer__sub">
                    <Icon name="icons:sub-mark" size="9px" />
                </button>
            </div>
        </div>

        <div class="footer__bottom">
            <div>© {{ new Date().getFullYear() }} Shopper</div>
            <div class="footer__social">
                <NuxtLink to="#">
                    <Icon name="icons:rutube" size="25px"/>
                </NuxtLink>
                <NuxtLink to="#">
                    <Icon name="icons:dzen" size="25px"/>
                </NuxtLink>
                <NuxtLink to="#">
                    <Icon name="icons:vk" size="25px"/>
                </NuxtLink>
                <NuxtLink to="#">
                    <Icon name="icons:tgc" size="25px"/>
                </NuxtLink>
            </div>
        </div>

    </div>

</template>

<style scoped>
.footer{
    margin: 0 auto;
    max-width: 1248px;
    display: flex;
    flex-direction: column;
    gap: 50px;
    padding: 0 16px;
    margin-bottom: 50px;
}

.footer hr {
    border: none;
    border-bottom: 1px solid var(--color-gray);
    height: 1px;
}

.footer__top{
    display: flex;
    justify-content: space-between;
}

.footer__top ul{
    display: flex;
    list-style: none;
    padding: 0;
    gap: 30px;
}

.footer__top ul a{
    font-size: 16px;
    font-weight: 400;
    text-transform: uppercase;
    color: var(--color-dark-gray);
    text-decoration: none;
}

.footer__top ul a:hover{
    font-size: 16px;
    font-weight: 400;
    text-transform: uppercase;
    color: var(--color-dark);
    cursor: pointer;
}

.footer__bottom{
    display: flex;
    justify-content: space-between;
}

.footer__social{
    display: flex;
    gap: 20px;

}

.footer__social a{
    color: var(--color-dark-gray);
}

.footer__social a:hover{
    color: var(--color-black);
}

.footer__form{
    position: relative;
}

.footer__form input{
    min-width: 396px;
}

.footer__sub{
    position:absolute;
    right: 0;
    top: 10px;
    background: none;
    border: none;
    color: var(--color-dark-gray);

}

.footer__sub:hover{
    cursor: pointer;
    color: var(--color-black);
}
</style>