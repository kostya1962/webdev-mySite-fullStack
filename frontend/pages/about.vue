<script setup lang="ts">
useSeoMeta({
    title: 'О нас',
    description: 'О компании интернет магазина Shopper',
    ogDescription: 'О компании интернет магазина Shopper',
});


</script>
<template>
    <section class="about">

        <div class="about-header">
            <h1>О нас</h1>
            <p class="subtitle">Мы делаем шикарные украшения для вас</p>
            <p class="description">
                Shopper — это ювелирный интернет-магазин, где современный дизайн сочетается
                с вниманием к деталям и качеству исполнения. Мы подбираем украшения,
                которые становятся частью образа и сохраняют свою актуальность вне времени.
            </p>
        </div>

        <div class="about-block">
        <h2>Тренды в украшениях</h2>

        <div class="image-card">
            <img
                src="../assets/img/about1.jpg"
                alt="Тренды в украшениях"
            />
        </div>

        <p class="text">
            Мы внимательно следим за мировыми ювелирными трендами, чтобы предлагать
            украшения, которые отражают дух времени. В коллекциях Shopper вы найдёте
            баланс между минимализмом, элегантностью и выразительными акцентами.
        </p>

        <ul class="list">
            <li>лаконичные формы и чистые линии</li>
            <li>актуальные материалы и универсальные дизайны</li>
        </ul>
        </div>

        <div class="about-block">
        <h2>Сделано с любовью</h2>

        <div class="image-card">
            <img
                src="../assets/img/about2.jpg"
                alt="Сделано с любовью"
            />
        </div>

        <p class="text">
            Каждое украшение в нашем магазине проходит тщательный отбор.
            Мы уделяем внимание не только внешнему виду, но и комфорту в носке,
            качеству материалов и долговечности изделий. Для нас важно,
            чтобы украшения приносили радость каждый день.
        </p>
        </div>
    </section>
</template>

<style scoped>
    .about {
        max-width: 900px;
        margin: 0 auto;
        padding: 80px 20px;
        font-family: "Inter", sans-serif;
        color: #1a1a1a;
    }

    .about-header {
        text-align: center;
        margin-bottom: 80px;
    }

    .about-header h1 {
        font-size: 36px;
        font-weight: 600;
        margin-bottom: 16px;
    }

    .subtitle {
        font-size: 18px;
        font-weight: 500;
        margin-bottom: 50px;
    }

    .description {
        max-width: 100%;
        margin: 5% 0;
        font-size: 14px;
        line-height: 1.7;
        color: #666;
        text-align: left;
    }

    .about-block {
        margin-bottom: 100px;
    }

    .about-block h2 {
        font-size: 24px;
        font-weight: 600;
        margin-bottom: 32px;
    }

    .image-card {
        width: 100%;
        height: 50%;
        border-radius: 12px;
        overflow: hidden;
        margin-bottom: 32px;
        background: var(--color-white);
    }

    .image-card img {
        width: 100%;
        height: auto;
        display: block;
    }

    .text {
        max-width: 600px;
        font-size: 14px;
        line-height: 1.7;
        color: #555;
        margin-bottom: 24px;
    }

    .list {
        padding-left: 18px;
        color: #555;
        font-size: 14px;
    }

    .list li {
        margin-bottom: 10px;
    }
</style>