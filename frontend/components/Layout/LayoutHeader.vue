<script setup lang="ts">

</script>

<template>
    <div class="header">
        <NuxtLink to="/">
            <IconLogo />
        </NuxtLink>
        <div class="header__right">
            <div class="header__menu">
                <NuxtLink to="/catalog">МАГАЗИН</NuxtLink>
                <NuxtLink to="/about">О НАС</NuxtLink>
            </div>
            <div class="header__hr"></div>
            <div class="header__icons">
                <NuxtLink to="#">
                    <Icon name="icons:search" size="21px" style="color: #707070"/>
                </NuxtLink>
                <NuxtLink to="/cart">
                    <Icon name="icons:cart" size="21px" />
                </NuxtLink>
                <NuxtLink to="/favorites">
                    <Icon name="icons:fav" size="21px" />
                </NuxtLink>
                <NuxtLink to="/account">
                    <Icon name="icons:user" size="21px" />
                </NuxtLink>
            </div>
        </div>
    </div>
</template>

<style scoped>
.header{
    margin: 0 auto;
    max-width: 1248px;
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    margin-top: 48px;
    border-bottom: 1px solid var(--color-gray);
    padding: 0 16px;
}

.header__right{
    display: flex;
    gap: 60px;
    align-items: center;
}

.header__menu{
    display: flex;
    gap: 60px;
    align-items: center;
    
}

.header__menu a{
    font-size: 16px;
    line-height: 27px;
    color: var(--color-black);
    text-decoration: none;
    padding-bottom: 48px;
}

.header__menu a:hover{
    color: var(--color-dark-gray);
}

.header__menu a.router-link-active{
    border-bottom: 1px solid var(--color-black);
}

.header__icons{
    display: flex;
    gap: 40px;
    align-items: center;
}

.header__icons a{
    padding-bottom: 48px;
}

.header__icons a.router-link-active{
    border-bottom: 1px solid var(--color-black);
}

.header__hr{
    border-left: 1px solid var(--color-dark-gray);
    height: 17px;
    align-self: flex-start;
    margin-top: 4px;
}
</style>